package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func OK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, data)
}

func NoContent(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}

func NotFound(c echo.Context) error {
	return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
}

func BadRequest(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Bad Request"})
}

func InternalError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, echo.Map{"msg": "Internal Server Error"})
}
