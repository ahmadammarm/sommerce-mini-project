package upload

import (
	"fmt"
	"os"

	"github.com/ahmadammarm/sommerce-mini-project/pkg/response"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/upload"
	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadImage(c *fiber.Ctx) error {
	// 1. Get the file from form-data key "image"
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Image file is required (form-data key: 'image')",
		})
	}

	// 2. Validate file (Size & Type)
	if err := upload.ValidateFile(file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// 3. Generate unique filename
	filename := upload.GenerateFilename(file.Filename)

	// 4. Ensure directory exists
	savePath := fmt.Sprintf("./public/uploads/%s", filename)
	if err := os.MkdirAll("./public/uploads", os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.WebResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create upload directory",
		})
	}

	// 5. Save the file
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.WebResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to save file to server",
		})
	}

	// 6. Return the public URL
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	publicURL := fmt.Sprintf("%s/uploads/%s", baseURL, filename)

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "File uploaded successfully",
		Data: fiber.Map{
			"url": publicURL,
		},
	})
}
