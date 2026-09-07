package transaction

import (
	"github.com/ahmadammarm/sommerce-mini-project/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	service TransactionService
}

func NewTransactionHandler(service TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) Checkout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	res, err := h.service.Checkout(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.WebResponse{
		Code:    fiber.StatusCreated,
		Message: "Checkout successful",
		Data:    res,
	})
}

func (h *TransactionHandler) GetMyTransactions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	res, err := h.service.GetMyTransactions(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.WebResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Transactions retrieved successfully",
		Data:    res,
	})
}

func (h *TransactionHandler) GetTransactionDetail(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	trxID := c.Params("id")

	res, err := h.service.GetTransactionDetail(c.Context(), userID, trxID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.WebResponse{
			Code:    fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Code:    fiber.StatusOK,
		Message: "Transaction detail retrieved successfully",
		Data:    res,
	})
}
