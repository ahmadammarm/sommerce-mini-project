package user

import (
	"github.com/ahmadammarm/sommerce-mini-project/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	res, err := h.service.Register(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.WebResponse{
		Code:    fiber.StatusCreated,
		Message: "User registered successfully",
		Data:    res,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	token, err := h.service.Login(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.WebResponse{
			Code:    fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Login successful",
		Data:    fiber.Map{"token": token},
	})
}

func (h *UserHandler) GetMyProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	res, err := h.service.GetMyProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.WebResponse{
			Code:    fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Profile retrieved successfully",
		Data:    res,
	})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	res, err := h.service.UpdateProfile(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Profile updated successfully",
		Data:    res,
	})
}
