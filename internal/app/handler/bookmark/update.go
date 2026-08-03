package bookmark

import (
	"errors"
	"net/http"

	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/TNJKL/bookmark-management/pkg/requestutils"
	"github.com/TNJKL/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// updateBookmarkBody defines the request body payload for updating a bookmark
type updateBookmarkBody struct {
	Description string `json:"description" example:"Google" validate:"lte=255"`
	Url         string `json:"url" example:"https://www.google.com" validate:"required,url,lte=2048"`
}

// updateBookmarkInput combines path URI parameters and request body payload for bookmark update
type updateBookmarkInput struct {
	ID string `uri:"id" validate:"required,uuid"`
	updateBookmarkBody
}

// UpdateBookmark     Update an existing bookmark of the current user
// @Summary        	  Update an existing bookmark of the current user
// @Description       Update an existing bookmark of the current user
// @Tags              bookmark
// @Security          BearerAuth
// @Accept            application/json
// @Produce           application/json
// @Param             id path string true "Bookmark ID" Format(uuid)
// @Param             input body updateBookmarkBody true "Input required"
// @Success           200 {object} response.Message "Success"
// @Router            /v1/bookmarks/{id} [put]
func (h *bookmarkHandler) UpdateBookmark(ctx *gin.Context) {
	input, uid, err := requestutils.BindInputFromRequestWithAuth[updateBookmarkInput](ctx)
	if err != nil {
		return
	}

	err = h.svc.UpdateBookmark(ctx, input.ID, uid, input.Description, input.Url)
	switch {
	case errors.Is(err, dbutils.ErrRecordNotFound):
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Message{Message: "Bookmark not found"})
		return
	case errors.Is(err, dbutils.ErrUniqueConstraint):
		ctx.AbortWithStatusJSON(http.StatusConflict, response.Message{Message: "Bookmark already exists"})
		return
	case errors.Is(err, nil):
		break
	default:
		log.Err(err).Str("userID", uid).Str("bookmarkID", input.ID).Msg("UpdateBookmark error")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	ctx.JSON(http.StatusOK, response.Message{Message: "Success"})
}
