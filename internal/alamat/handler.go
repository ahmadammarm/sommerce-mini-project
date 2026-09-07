package alamat

import (
	"github.com/ahmadammarm/sommerce-mini-project/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type AlamatHandler struct {
	service AlamatService
}

func NewAlamatHandler(service AlamatService) *AlamatHandler {
	return &AlamatHandler{service: service}
}

func (h *AlamatHandler) CreateAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req CreateAlamatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	res, err := h.service.CreateAlamat(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.WebResponse{
		Code:    fiber.StatusCreated,
		Message: "Address created successfully",
		Data:    res,
	})
}

func (h *AlamatHandler) GetMyAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	res, err := h.service.GetMyAlamat(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.WebResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Addresses retrieved successfully",
		Data:    res,
	})
}

func (h *AlamatHandler) GetAlamatByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	alamatID := c.Params("id")

	res, err := h.service.GetAlamatByID(c.Context(), userID, alamatID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.WebResponse{
			Code:    fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Address retrieved successfully",
		Data:    res,
	})
}

func (h *AlamatHandler) UpdateAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	alamatID := c.Params("id")

	var req UpdateAlamatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	res, err := h.service.UpdateAlamat(c.Context(), userID, alamatID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Address updated successfully",
		Data:    res,
	})
}

func (h *AlamatHandler) DeleteAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	alamatID := c.Params("id")

	if err := h.service.DeleteAlamat(c.Context(), userID, alamatID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Address deleted successfully",
	})
}
