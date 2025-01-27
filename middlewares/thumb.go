package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"os"
	"path/filepath"
)

func ThumbMiddleware(app *fiber.App) {
	app.Get("/uploaded/images/:year/:month/:file", func(c *fiber.Ctx) error {
		// Decode the file name from the URL
		rawFileName := c.Params("file")
		fileName, err := url.QueryUnescape(rawFileName)
		if err != nil {

			return c.Status(fiber.StatusBadRequest).SendString("Invalid file name")
		}

		// Handle thumb_ prefix logic
		finalFileName := handleThumbLogic(c.Params("year"), c.Params("month"), fileName)

		// File path on the local server
		filePath := filepath.Clean(fmt.Sprintf("public/uploaded/images/%s/%s/%s", c.Params("year"), c.Params("month"), finalFileName))

		// Check if file exists locally
		if fileExistsLocally(filePath) {
			return c.SendFile(filePath)
		}

		// File not found
		return c.Status(fiber.StatusNotFound).SendString("File not found")
	})
}

// Handles the logic for checking thumb_ prefix
func handleThumbLogic(year, month, fileName string) string {
	if len(fileName) >= 6 && fileName[:6] == "thumb_" {
		// Thumbnail file path
		thumbFilePath := filepath.Clean(fmt.Sprintf("public/uploaded/images/%s/%s/%s", year, month, fileName))
		if fileExistsLocally(thumbFilePath) {

			return fileName // Return the thumbnail if it exists
		}

		// Fallback to the original file if thumbnail doesn't exist
		originalFileName := fileName[6:] // Remove "thumb_" prefix

		return originalFileName
	}

	// No thumb_ prefix, return the original file name
	return fileName
}

func fileExistsLocally(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {

		return false
	}
	if info.IsDir() {

		return false
	}
	return true
}
