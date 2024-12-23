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
	"strings"
	"time"
)

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

func makeUploadable(src io.Reader, fileType string, ext string, fileName string) (map[string]string, error) {
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

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return map[string]string{
			"httpPath": "",
			"basePath": "",
			"fileName": "",
		}, errors.New("file create error")
	}

	if config.Config.Image.MaxSize > 0 {
		targetSizeBytes := int64(config.Config.Image.MaxSize * 1e6)
		errO := optimizeImage(publicPath+uploadPath+newFileName, targetSizeBytes)
		if errO != nil {
			fmt.Print(errO.Error())
			return map[string]string{
				"httpPath": "",
				"basePath": "",
				"fileName": "",
			}, errO

		}
	}

	return map[string]string{
		"httpPath": uploadPath + newFileName,
		"basePath": fullPath,
		"fileName": newFileName,
	}, nil

}
func optimizeImage(filePath string, targetSize int64) error {
	// Check the file size first
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// If the file is already smaller than the target size, return
	if fileInfo.Size() <= targetSize {
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Resize the image to reduce size
	img = resize.Thumbnail(1024, 1024, img, resize.Lanczos3)

	var buf bytes.Buffer

	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 75})
	case "png":
		err = png.Encode(&buf, img)
	default:
		return nil
	}

	if err != nil {
		return err
	}

	if int64(buf.Len()) > targetSize {
		return nil
	}

	// Write the optimized image back to disk
	return os.WriteFile(filePath, buf.Bytes(), 0644)
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
	//srcMime := src

	var ext_ = filepath.Ext(file.Filename)
	ext := strings.ToLower(strings.TrimPrefix(ext_, "."))
	var fileType string = "images"
	rules := govalidator.MapData{
		//"file:file": []string{"ext:jpg,png,jpeg,svg,JPG,PNG,JPEG,SVG", "size:100000", "mime:jpg,png,jpeg,svg,JPG,PNG,JPEG,SVG", "required"},
		"file:file": []string{"ext:jpg,png,jpeg,svg,gif,JPG,PNG,JPEG,SVG,GIF", "size:100000000", "required"},
	}
	mimeTypes := []string{
		"image/svg+xml",
		"image/jpeg",
		"image/png",
		"image/gif",
	}

	if ext == "dwg" || ext == "pdf" || ext == "zip" || ext == "swf" || ext == "doc" || ext == "docx" || ext == "csv" || ext == "xls" || ext == "xlsx" || ext == "ppt" || ext == "pptx" {
		rules = govalidator.MapData{
			"file:file": []string{"ext:xls,xlsx,doc,docx,pdf,ppt,pptx,csv,zip,dwg,XLS,XLSX,DOC,DOCX,PDF,PPT,PPTX,CSV,ZIP,DWG", "size:400000000", "required"},
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
			"file:file": []string{"ext:mp4,m4v,avi,webm,MP4,M4V,AVI,WEBM", "size:8000000000", "required"},
		}
		mimeTypes = []string{
			"video/mp4",
			"video/x-m4v",
			"video/x-msvideo",
			"video/webm",
		}
		fileType = "videos"
	}
	if ext == "mp3" || ext == "wav" {
		rules = govalidator.MapData{
			"file:file": []string{"ext:mp3,wav,MP3,WAV", "size:500000000", "required"},
		}
		mimeTypes = []string{
			"audio/mpeg",
			"audio/wav",
		}
		fileType = "audios"
	}

	if ext == "so" {
		rules = govalidator.MapData{
			"file:file": []string{"ext:so", "size:40000000", "required"},
		}
		mimeTypes = []string{
			"application/x-sharedlib",
		}
		fileType = "sharedlib"
	}

	//mimeType, _, err  := mimetype.DetectReader(srcMime)

	mimeType := "1"

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
		//return c.Status(http.StatusBadRequest).JSON(map[string]string{
		//	"status": "false",
		//	"message": "file mime not allowed",
		//})
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
	upload, uerr := makeUploadable(src, fileType, ext_, file.Filename)

	if uerr != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"status":  "false",
			"message": uerr.Error(),
		})
	}
	return c.SendString(upload["httpPath"])
}
