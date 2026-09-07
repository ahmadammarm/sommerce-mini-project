package toko

import (
	"github.com/ahmadammarm/sommerce-mini-project/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type TokoHandler struct {
	service TokoService
}

func NewTokoHandler(service TokoService) *TokoHandler {
	return &TokoHandler{service: service}
}

func (h *TokoHandler) GetMyToko(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	res, err := h.service.GetMyToko(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.WebResponse{
			Code:    fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Store profile retrieved successfully",
		Data:    res,
	})
}

func (h *TokoHandler) UpdateMyToko(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req UpdateTokoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	res, err := h.service.UpdateMyToko(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Store profile updated successfully",
		Data:    res,
	})
}

func (h *TokoHandler) GetTokoByID(c *fiber.Ctx) error {
	tokoID := c.Params("id")

	res, err := h.service.GetTokoByID(c.Context(), tokoID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.WebResponse{
			Code:    fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Store profile retrieved successfully",
		Data:    res,
	})
}
