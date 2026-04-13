package krud

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/dataform"
	"github.com/lambda-platform/lambda/datagrid"
	"github.com/lambda-platform/lambda/krud/handlers"
	"github.com/lambda-platform/lambda/krud/krudMW"
	"github.com/lambda-platform/lambda/krud/utils"
)

func Set(e *fiber.App, GetGridMODEL func(schema_id string) datagrid.Datagrid, GetMODEL func(schema_id string) dataform.Dataform, krudMiddleWares []fiber.Handler, KrudWithPermission bool, publicForms []string, GetPermissionCustom func(c *fiber.Ctx, vbType string) krudMW.PermissionObj, ignoreList []string) {
	if config.Config.App.Migrate == "true" {
		utils.AutoMigrateSeed()
	}

	g := e.Group("/lambda/krud")
	if len(krudMiddleWares) >= 1 {
		for _, krudMiddleWare := range krudMiddleWares {
			g.Use(krudMiddleWare)
		}
	}

	p := e.Group("/lambda/puzzle")
	if len(krudMiddleWares) >= 1 {
		for _, krudMiddleWare := range krudMiddleWares {
			p.Use(krudMiddleWare)
		}
	}

	g.Post("/excel/:schemaId", agentMW.IsLoggedIn(), handlers.ExportExcel(GetGridMODEL))
	g.Post("/import-excel/:schemaId", agentMW.IsLoggedIn(), handlers.ImportExcel(GetGridMODEL))
	g.Post("/print/:schemaId", agentMW.IsLoggedIn(), handlers.Print(GetGridMODEL))
	if KrudWithPermission {
		g.Post("/:schemaId/filter-options", agentMW.IsLoggedIn(), handlers.FilterOptions(GetGridMODEL))
		g.Post("/update-row/:schemaId", agentMW.IsLoggedIn(), krudMW.PermissionDelete(GetPermissionCustom, ignoreList), handlers.UpdateRow(GetGridMODEL))
		g.Post("/:schemaId/:action", agentMW.IsLoggedIn(), krudMW.PermissionCreate(GetPermissionCustom, ignoreList), handlers.Crud(GetMODEL))
		g.Post("/:schemaId/:action/:id", agentMW.IsLoggedIn(), krudMW.PermissionEdit(GetPermissionCustom, ignoreList), handlers.Crud(GetMODEL))
		g.Delete("/delete/:schemaId/:id", agentMW.IsLoggedIn(), krudMW.PermissionDelete(GetPermissionCustom, ignoreList), handlers.Delete(GetGridMODEL))

	} else {
		g.Post("/:schemaId/filter-options", agentMW.IsLoggedIn(), handlers.FilterOptions(GetGridMODEL))
		g.Post("/update-row/:schemaId", agentMW.IsLoggedIn(), handlers.UpdateRow(GetGridMODEL))
		g.Post("/:schemaId/:action", agentMW.IsLoggedIn(), handlers.Crud(GetMODEL))
		g.Post("/:schemaId/:action/:id", agentMW.IsLoggedIn(), handlers.Crud(GetMODEL))
		g.Delete("/delete/:schemaId/:id", agentMW.IsLoggedIn(), handlers.Delete(GetGridMODEL))
	}

	/*
		OTHER
	*/
	g.Post("/upload", handlers.Upload)
	g.Options("/upload", handlers.Upload)
	//g.Post("/upload", handlers.Upload, agentMW.IsLoggedIn())
	//g.OPTIONS("/upload", handlers.Upload, agentMW.IsLoggedIn())
	g.Post("/unique", handlers.CheckUnique)
	g.Get("/today", handlers.Now)
	g.Get("/now", handlers.Now)
	g.Post("/check_current_password", agentMW.IsLoggedIn(), handlers.CheckCurrentPassword)

	// Regenerate all thumbnails with correct EXIF orientation
	g.Get("/regenerate-thumbs", agentMW.IsLoggedIn(), handleRegenerateThumbs)

	/*
		PUBLIC CURDS
	*/
	if len(publicForms) >= 1 {

		contains := func(slice []string, item string) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		}
		// Middleware to check if schemaId is in publicForms
		publicSchemaCheck := func(c *fiber.Ctx) error {
			schemaId := c.Params("schemaId")

			if schemaId == "" || !contains(publicForms, schemaId) {
				return c.Status(403).JSON(fiber.Map{"error": "Schema is not public"})
			}
			return c.Next()
		}

		public := e.Group("/lambda/krud-public")
		if len(krudMiddleWares) >= 1 {
			for _, krudMiddleWare := range krudMiddleWares {
				public.Use(krudMiddleWare)
			}
		}
		public.Post("/:schemaId/:action", publicSchemaCheck, handlers.Crud(GetMODEL))
		public.Post("/:schemaId/:action/:id", publicSchemaCheck, handlers.Crud(GetMODEL))
	}

}

// handleRegenerateThumbs regenerates all thumb_ prefixed images
// with correct EXIF orientation from the original source images.
func handleRegenerateThumbs(c *fiber.Ctx) error {
	rootDir := "public/uploaded/images"

	var totalFound, totalFixed, totalFailed, totalSkipped int
	var results []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		fileName := info.Name()
		if !strings.HasPrefix(fileName, "thumb_") {
			return nil
		}

		totalFound++

		dir := filepath.Dir(path)
		originalName := strings.TrimPrefix(fileName, "thumb_")
		originalPath := filepath.Join(dir, originalName)

		if _, err := os.Stat(originalPath); os.IsNotExist(err) {
			totalSkipped++
			results = append(results, fmt.Sprintf("⚠️ Skip (no original): %s", path))
			return nil
		}

		if err := regenerateThumb(originalPath, path); err != nil {
			totalFailed++
			results = append(results, fmt.Sprintf("❌ Failed: %s → %v", path, err))
			return nil
		}

		totalFixed++
		results = append(results, fmt.Sprintf("✅ Fixed: %s", path))
		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": fmt.Sprintf("Walk error: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Thumbnail regeneration complete",
		"total":   totalFound,
		"fixed":   totalFixed,
		"skipped": totalSkipped,
		"failed":  totalFailed,
		"details": results,
	})
}

func regenerateThumb(originalPath, thumbPath string) error {
	img, err := imaging.Open(originalPath, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	img = imaging.Fit(img, 500, 500, imaging.Lanczos)

	ext := strings.ToLower(filepath.Ext(originalPath))
	var buf bytes.Buffer

	switch ext {
	case ".jpeg", ".jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 75})
	case ".png":
		err = png.Encode(&buf, img)
	default:
		return fmt.Errorf("unsupported format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return os.WriteFile(thumbPath, buf.Bytes(), 0644)
}
