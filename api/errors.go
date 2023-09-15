package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiErr, ok := err.(Error); ok {
		return c.Status(apiErr.Code).JSON(apiErr)
	}
	apiError := NewError(fiber.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrUnauthorized() Error {
	return NewError(http.StatusUnauthorized, "unauthorized")
}

func ErrInvalidID() Error {
	return NewError(http.StatusBadRequest, "invalid id")
}

func ErrBadRequest() Error {
	return NewError(http.StatusBadRequest, "bad request")
}

func ErrResourceNotFound(res string) Error {
	return NewError(http.StatusNotFound, res+" not found")
}
