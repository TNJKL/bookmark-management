package user

import (
	"errors"
	"net/http"

	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/TNJKL/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// registerInputBody defines the JSON request payload for user registration
type registerInputBody struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,strong_password"`
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
}

// registerResponse defines the JSON response structure for a successful registration
type registerResponse struct {
	Message string      `json:"message"`
	Data    *model.User `json:"data"`
}

// Register handles user registration requests.
// It validates the JSON input, delegates to the service layer for user creation,
// and returns the created user or an appropriate error response.
//
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags User
// @Accept json
// @Produce json
// @Param body body registerInputBody true "User registration details"
// @Success 200 {object} registerResponse "Register an user successfully!"
// @Router /v1/users/register [post]
func (h *userHandler) Register(ctx *gin.Context) {
	input := &registerInputBody{}
	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	res, err := h.svc.CreateUser(ctx, input.Username, input.Password, input.DisplayName, input.Email)
	switch {
	case errors.Is(err, dbutils.ErrDuplicationUsername):
		ctx.AbortWithStatusJSON(http.StatusConflict, response.Message{
			Message: "Username already taken",
		})
		return
	case errors.Is(err, dbutils.ErrDuplicationEmail):
		ctx.AbortWithStatusJSON(http.StatusConflict, response.Message{
			Message: "Email already taken",
		})
		return
	case err == nil:
	default:
		log.Error().Err(err).Str("from", "handler.userHandler.Register").Msg("Failed to register user")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	ctx.JSON(http.StatusOK, registerResponse{
		Message: "Register an user successfully!",
		Data:    res,
	})

}
