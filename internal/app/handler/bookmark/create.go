package bookmark

import (
	"errors"
	"net/http"

	"github.com/TNJKL/bookmark-management/internal/app/handler/dto"
	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/TNJKL/bookmark-management/pkg/requestutils"
	"github.com/TNJKL/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// createBookmarkInput defines the request body payload for creating a new bookmark
type createBookmarkInput struct {
	Description string `json:"description" example:"A bookmark description" validate:"lte=255"`
	Url         string `json:"url" example:"http://google.com" validate:"required,url,lte=2048"`
}

// CreateBookmark     Create a new bookmark for a bookmark group
// @Summary        	  Create a new bookmark for a bookmark group
// @Description       Create a new bookmark for a bookmark group
// @Tags              bookmark
// @Security          BearerAuth
// @Accept            application/json
// @Produce           application/json
// @Param             input body createBookmarkInput true "Input required"
// @Success           200 {object} object{data=model.Bookmark,message=string} "Success"
// @Router            /v1/bookmarks [post]
func (h *bookmarkHandler) CreateBookmark(ctx *gin.Context) {
	//get input
	input, uid, err := requestutils.BindInputFromRequestWithAuth[createBookmarkInput](ctx)
	if err != nil {
		return
	}

	///call service
	res, err := h.svc.CreateBookmark(ctx, input.Description, input.Url, uid)
	switch {
	case errors.Is(err, dbutils.ErrUniqueConstraint):
		ctx.AbortWithStatusJSON(http.StatusConflict, response.Message{Message: "Bookmark already exists"})
		return
	case errors.Is(err, nil):
		break
	default:
		log.Err(err).Str("userID", uid).Msg("AddBookmark error")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	ctx.JSON(http.StatusOK, dto.SuccessResponse[*model.Bookmark]{
		Message: "Create a bookmark successfully!",
		Data:    res,
	})
}
