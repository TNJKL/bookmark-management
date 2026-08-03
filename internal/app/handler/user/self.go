package user

import (
	"errors"
	"net/http"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/TNJKL/bookmark-management/pkg/requestutils"
	"github.com/TNJKL/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// getSelfInfoResponse defines the JSON response structure for getting current user info
type getSelfInfoResponse struct {
	Data *model.User `json:"data"`
}

// updateSelfInfoBody defines the JSON request payload for updating current user info
type updateSelfInfoBody struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email" binding:"required,email"`
}

// updateSelfInfoResponse defines the JSON response structure for updating current user info
type updateSelfInfoResponse struct {
	Message string `json:"message"`
}

// GetSelfInfo    get your current information
// @Summary      get your current information
// @Description  get your current information
// @Tags         user
// @Security     BearerAuth
// @Accept       application/json
// @Produce      application/json
// @Success      200 {object} getSelfInfoResponse "Success"
// @Router       /v1/self/info [get]
func (h *userHandler) GetSelfInfo(ctx *gin.Context) {
	uid, err := requestutils.GetUserIDFromRequest(ctx)
	if err != nil {
		return
	}
	currentUser, err := h.svc.GetSelfInfo(ctx, uid)
	if err != nil {
		log.Err(err).Msg("Failed to get self info")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	ctx.JSON(http.StatusOK, getSelfInfoResponse{Data: currentUser})
}

// UpdateSelfInfo  edit your current information
// @Summary      edit your current information
// @Description  edit your current information
// @Tags         user
// @Security     BearerAuth
// @Accept       application/json
// @Produce      application/json
// @Param        input body updateSelfInfoBody true "Input required"
// @Success      200 {object} updateSelfInfoResponse "Success"
// @Router       /v1/self/info [put]
func (h *userHandler) UpdateSelfInfo(ctx *gin.Context) {
	uid, err := requestutils.GetUserIDFromRequest(ctx)
	if err != nil {
		return
	}
	input, err := requestutils.BindInputFromRequest[updateSelfInfoBody](ctx)
	if err != nil {
		return
	}
	err = h.svc.UpdateSelfInfo(ctx, uid, input.DisplayName, input.Email)
	switch {
	case errors.Is(err, dbutils.ErrDuplicationEmail):
		ctx.AbortWithStatusJSON(http.StatusConflict, response.Message{Message: "Email already taken"})
		return
	case err == nil:
	default:
		log.Err(err).Msg("Failed to update self info")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	ctx.JSON(http.StatusOK, updateSelfInfoResponse{
		Message: "Edit current user successfully!",
	})
}
