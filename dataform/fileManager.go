package dataform

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/config"
	"github.com/nfnt/resize"
	"github.com/thedevsaddam/govalidator"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Default maximum file size in bytes (10MB)
var defaultMaxFileSize = 10 * 1024 * 1024

func sanitizeFileName(fileName string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	return reg.ReplaceAllString(fileName, "_")
}

func CheckFileExist(filepath string, fileName string, fileType string, ext string, i int) string {
	newFileName := ""
	if i > 0 {
		newFileName = fmt.Sprintf("%v", i) + "-" + fileName + ext
	} else {
		newFileName = fileName + ext
	}
	_, err := os.Stat(filepath + newFileName)

	if os.IsNotExist(err) {
		return newFileName
	} else {
		i = i + 1
		return CheckFileExist(filepath, fileName, fileType, ext, i)
	}
}

func sanitizeSvgContent(content []byte) ([]byte, error) {
	svgContent := string(content)

	// `<script>` болон бусад хортой элементүүдийг устгах
	if strings.Contains(svgContent, "<script") {
		return nil, errors.New("SVG file contains prohibited <script> tag")
	}

	// Doctype болон entity шалгалт хийх
	if strings.Contains(svgContent, "<!DOCTYPE") {
		return nil, errors.New("SVG file contains prohibited <!DOCTYPE>")
	}

	return content, nil
}

func makeUploadable(src io.Reader, fileType string, ext string, fileName string, mimeType string) (map[string]string, error) {
	var name = strings.TrimRight(fileName, ext)
	currentTime := time.Now()
	year := fmt.Sprintf("%v", currentTime.Year())
	month := fmt.Sprintf("%v", currentTime.Month())

	var publicPath string = "public"
	var uploadPath string = "/uploaded/" + fileType + "/" + year + "/" + month + "/"
	var fullPath string = publicPath + uploadPath

	if fileType == "sharedlib" {
		publicPath = ""
		uploadPath = ""
		fullPath = "sharedlib/"
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		os.MkdirAll(fullPath, 0700)
		// Create your file
	}

	var i int = 0
	newFileName := CheckFileExist(fullPath, name, fileType, ext, i)
	// Destination
	dst, err := os.Create(fullPath + newFileName)
	if err != nil {
		return map[string]string{
			"httpPath": "",
			"basePath": "",
			"fileName": "",
		}, errors.New("file create error")
	}
	defer dst.Close()

	// Copy and sanitize SVG content if necessary
	var content []byte
	if mimeType == "image/svg+xml" {
		content, err = io.ReadAll(src)
		if err != nil {
			return map[string]string{
				"httpPath": "",
				"basePath": "",
				"fileName": "",
			}, errors.New("file read error")
		}
		content, err = sanitizeSvgContent(content)

		if err != nil {
			return map[string]string{
				"httpPath": "",
				"basePath": "",
				"fileName": "",
			}, errors.New("file sanitization error")
		}
		_, err = dst.Write(content)
	} else {
		_, err = io.Copy(dst, src)
	}

	if err != nil {
		return map[string]string{
			"httpPath": "",
			"basePath": "",
			"fileName": "",
		}, errors.New("file write error")
	}

	// Create a thumbnail only for jpg, jpeg, and png files
	if fileType == "images" && (ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".JPG" || ext == ".JPEG" || ext == ".PNG") {
		thumbnailPath := fullPath + "thumb_" + newFileName
		err = createThumbnail(publicPath+uploadPath+newFileName, thumbnailPath, 500*1024) // 500KB
		if err != nil {
			return map[string]string{
				"httpPath": "",
				"basePath": "",
				"fileName": "",
			}, errors.New("thumbnail creation error")
		}
	}

	return map[string]string{
		"httpPath": uploadPath + newFileName,
		"basePath": fullPath,
		"fileName": newFileName,
	}, nil
}

func createThumbnail(inputPath, outputPath string, maxSize int64) error {
	// Open the image file
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Resize the image for the thumbnail
	img = resize.Thumbnail(500, 500, img, resize.Lanczos3)

	var buf bytes.Buffer

	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 75})
	case "png":
		err = png.Encode(&buf, img)
	default:
		return errors.New("unsupported image format for thumbnail")
	}

	if err != nil {
		return err
	}

	// Check the size of the thumbnail
	if int64(buf.Len()) > maxSize {
		return errors.New("thumbnail size exceeds the limit of 500KB")
	}

	// Save the thumbnail
	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

func Upload(c *fiber.Ctx) error {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"status":  "false",
			"message": "file not found",
		})
	}

	//
	src, err := file.Open()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"status":  "false",
			"message": "server error",
		})
	}
	defer src.Close()

	if config.Config.File.FileMaxSize > 0 {
		defaultMaxFileSize = config.Config.File.FileMaxSize * 1024 * 1024
	}

	var ext_ = filepath.Ext(file.Filename)
	ext := strings.ToLower(strings.TrimPrefix(ext_, "."))
	sanitizedFileName := sanitizeFileName(file.Filename)
	var fileType string = "images"
	rules := govalidator.MapData{
		"file:file": []string{"ext:jpg,png,jpeg,svg,gif,JPG,PNG,JPEG,SVG,GIF", fmt.Sprintf("size:%d", defaultMaxFileSize), "required"},
	}
	mimeTypes := []string{
		"image/svg+xml",
		"image/jpeg",
		"image/png",
		"image/gif",
		"application/octet-stream",
	}

	if ext == "dwg" || ext == "pdf" || ext == "zip" || ext == "swf" || ext == "doc" || ext == "docx" || ext == "csv" || ext == "xls" || ext == "xlsx" || ext == "ppt" || ext == "pptx" {
		rules = govalidator.MapData{
			"file:file": []string{"ext:xls,xlsx,doc,docx,pdf,ppt,pptx,csv,zip,dwg,XLS,XLSX,DOC,DOCX,PDF,PPT,PPTX,CSV,ZIP,DWG", fmt.Sprintf("size:%d", defaultMaxFileSize), "required"},
		}
		mimeTypes = []string{
			"application/acad",
			"application/pdf",
			"application/x-shockwave-flash",
			"application/x-shockwave-flash2-preview",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.ms-excel",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/vnd.ms-powerpoint",
			"application/vnd.openxmlformats-officedocument.presentationml.presentation",
			"text/csv",
		}
		fileType = "documents"
	}
	if ext == "mp4" || ext == "m4v" || ext == "avi" || ext == "webm" {
		rules = govalidator.MapData{
			"file:file": []string{"ext:mp4,m4v,avi,webm,MP4,M4V,AVI,WEBM", fmt.Sprintf("size:%d", defaultMaxFileSize), "required"},
		}
		mimeTypes = []string{
			"video/mp4",
			"video/x-m4v",
			"video/x-msvideo",
			"video/webm",
			"application/octet-stream",
		}
		fileType = "videos"
	}
	if ext == "mp3" || ext == "wav" || ext == "aac" || ext == "m4a" {
		rules = govalidator.MapData{
			"file:file": []string{"ext:mp3,wav,aac,m4a,MP3,WAV,AAC,M4A", fmt.Sprintf("size:%d", defaultMaxFileSize), "required"},
		}
		mimeTypes = []string{
			"audio/mpeg",
			"audio/wav",
			"audio/aac",
			"audio/m4a",
			"application/octet-stream",
		}
		fileType = "audios"
	}

	if ext == "so" {
		rules = govalidator.MapData{
			"file:file": []string{"ext:so", fmt.Sprintf("size:%d", defaultMaxFileSize), "required"},
		}
		mimeTypes = []string{
			"application/x-sharedlib",
		}
		fileType = "sharedlib"
	}

	mimeType := file.Header.Get("Content-Type")

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"status":  "false",
			"message": "can't parse file mime, server error",
		})
	}
	mimeAllowed := false
	for _, m := range mimeTypes {
		if m == mimeType {
			mimeAllowed = true
		}
	}

	if mimeAllowed == false {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"status":  "false",
			"message": "file mime not allowed, " + mimeType,
		})
	}

	messages := govalidator.MapData{
		"file:file": []string{"ext:file not allowed", "required:File required", "size:File size too big"},
	}

	r := http.Request{}

	fasthttpadaptor.ConvertRequest(c.Context(), &r, true)
	r.Host = string(c.Request().Host())
	opts := govalidator.Options{
		Request:  &r,    // request object
		Rules:    rules, // rules map,
		Messages: messages,
	}
	v := govalidator.New(opts)
	e := v.Validate()

	if len(e) >= 1 {
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status":  false,
			"message": e,
		})
	}
	upload, uerr := makeUploadable(src, fileType, ext_, sanitizedFileName, mimeType)

	if uerr != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"status":  "false",
			"message": uerr.Error(),
		})
	}
	return c.SendString(upload["httpPath"])
}
