package user

import (
	"errors"
	"net/http"

	"github.com/TNJKL/bookmark-management/internal/app/service/user"
	"github.com/TNJKL/bookmark-management/pkg/requestutils"
	"github.com/TNJKL/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// loginInputBody defines the JSON request payload for user login
type loginInputBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,strong_password"`
}

// loginResponse defines the JSON response structure for a successful login
type loginResponse struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

// Login handles user authentication requests
// Login        Authentication endpoint
// @Summary     Return a jwt token if the input is correct
// @Description Return a jwt token if the input is correct
// @Tags        user
// @Accept      application/json
// @Produce     application/json
// @Param       input body loginInputBody true "Input required"
// @Success     200 {object} loginResponse "Success"
// @Router      /v1/users/login [post]
func (h *userHandler) Login(ctx *gin.Context) {
	//doc input
	input, err := requestutils.BindInputFromRequest[loginInputBody](ctx)
	if err != nil {
		return
	}

	//call service -- service tra ve token
	token, err := h.svc.Login(ctx, input.Username, input.Password)
	switch {
	case errors.Is(err, user.ErrInvalidCredentials):
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Message{Message: "invalid credentials"})
		return
	case err == nil:
	default:
		log.Err(err).Msg("Failed to login")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	//return token
	ctx.JSON(http.StatusOK, loginResponse{
		Message: "Logged in successfully",
		Data:    token,
	})
}
