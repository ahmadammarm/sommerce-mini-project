package category

import (
	"github.com/ahmadammarm/sommerce-mini-project/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	service CategoryService
}

func NewCategoryHandler(service CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	res, err := h.service.CreateCategory(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.WebResponse{
		Code:    fiber.StatusCreated,
		Message: "Category created successfully",
		Data:    res,
	})
}

func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	res, err := h.service.GetAllCategories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.WebResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Categories retrieved successfully",
		Data:    res,
	})
}

func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := h.service.GetCategoryByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.WebResponse{
			Code:    fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Category retrieved successfully",
		Data:    res,
	})
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	var req CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	res, err := h.service.UpdateCategory(c.Context(), id, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Category updated successfully",
		Data:    res,
	})
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteCategory(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Category deleted successfully",
	})
}
