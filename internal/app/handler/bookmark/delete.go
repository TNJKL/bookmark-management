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

// deleteBookmarkInput defines the URI path parameter for deleting a bookmark
type deleteBookmarkInput struct {
	ID string `uri:"id" validate:"required,uuid"`
}

// DeleteBookmark     Delete a bookmark of the current user
// @Summary        	  Delete a bookmark of the current user
// @Description       Delete a bookmark of the current user
// @Tags              bookmark
// @Security          BearerAuth
// @Accept            application/json
// @Produce           application/json
// @Param             id path string true "Bookmark ID" Format(uuid)
// @Success           200 {object} response.Message "Success"
// @Router            /v1/bookmarks/{id} [delete]
func (h *bookmarkHandler) DeleteBookmark(ctx *gin.Context) {
	input, uid, err := requestutils.BindInputFromRequestWithAuth[deleteBookmarkInput](ctx)
	if err != nil {
		return
	}

	err = h.svc.DeleteBookmark(ctx, input.ID, uid)
	switch {
	case errors.Is(err, dbutils.ErrRecordNotFound):
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Message{Message: "Bookmark not found"})
		return
	case errors.Is(err, nil):
		break
	default:
		log.Err(err).Str("userID", uid).Str("bookmarkID", input.ID).Msg("DeleteBookmark error")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	ctx.JSON(http.StatusOK, response.Message{Message: "Success"})
}
