package requestutils

import (
	"errors"
	"io"
	"net/http"

	"github.com/TNJKL/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	InputValidator = validator.New(validator.WithRequiredStructEnabled())
)

// BindInputFromRequest binds and validates the request input into the target struct
func BindInputFromRequest[T any](c *gin.Context) (*T, error) {
	reqInput := new(T)

	// Skip JSON binding for GET requests to avoid EOF error on empty body
	if c.Request.Method != http.MethodGet {
		if err := c.ShouldBindJSON(reqInput); err != nil && !errors.Is(err, io.EOF) {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
			return nil, err
		}
	}

	if err := c.ShouldBindUri(reqInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	if err := c.ShouldBindQuery(reqInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	if err := c.ShouldBindHeader(reqInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	if err := InputValidator.Struct(reqInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}
	return reqInput, nil

}
