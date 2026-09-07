package produk

import (
	"github.com/ahmadammarm/sommerce-mini-project/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type ProdukHandler struct {
	service ProdukService
}

func NewProdukHandler(service ProdukService) *ProdukHandler {
	return &ProdukHandler{service: service}
}

func (h *ProdukHandler) CreateProduk(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req CreateProdukRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	res, err := h.service.CreateProduk(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.WebResponse{
		Code:    fiber.StatusCreated,
		Message: "Product created successfully",
		Data:    res,
	})
}

func (h *ProdukHandler) UpdateProduk(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	produkID := c.Params("id")

	var req UpdateProdukRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	res, err := h.service.UpdateProduk(c.Context(), userID, produkID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Product updated successfully",
		Data:    res,
	})
}

func (h *ProdukHandler) GetProdukByID(c *fiber.Ctx) error {
	produkID := c.Params("id")

	res, err := h.service.GetProdukByID(c.Context(), produkID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.WebResponse{
			Code:    fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Product retrieved successfully",
		Data:    res,
	})
}

func (h *ProdukHandler) GetAllProduk(c *fiber.Ctx) error {
	var filter ProductFilterDTO
	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid query parameters",
		})
	}

	res, err := h.service.GetAllProduk(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.WebResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	// Because PagedResponse is already nicely formatted with Meta, we return it directly,
	// or we can wrap it. Let's return it directly as standard HTTP response.
	return c.Status(fiber.StatusOK).JSON(res)
}
