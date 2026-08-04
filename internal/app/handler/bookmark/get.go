package bookmark

import (
	"errors"
	"net/http"

	"github.com/TNJKL/bookmark-management/internal/app/handler/dto"
	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/pkg/requestutils"
	"github.com/TNJKL/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// getBookmarkInput defines the query parameters for retrieving paginated bookmarks
type getBookmarkInput struct {
	Page  int `form:"page" validate:"gte=1"`
	Limit int `form:"limit" validate:"gte=1"`
}

// GetBookmarks     Get a list of bookmarks of the current user
// @Summary        Get a list of bookmarks of the current user
// @Description    Get a list of bookmarks of the current user
// @Tags           bookmark
// @Security       BearerAuth
// @Accept         application/json
// @Produce        application/json
// @Param          page query integer false "page" Format(int32) default(1)
// @Param          limit query integer false "limit" Format(int32) default(20)
// @Success        200 {object} object{data=[]model.Bookmark} "Success"
// @Router         /v1/bookmarks [get]
func (h *bookmarkHandler) GetBookmarks(ctx *gin.Context) {
	input, uid, err := requestutils.BindInputFromRequestWithAuth[getBookmarkInput](ctx)
	if err != nil {
		return
	}
	res, err := h.svc.GetBookmarks(ctx, uid, input.Page, input.Limit)
	switch {
	case errors.Is(err, nil):
		break
	default:
		log.Err(err).Str("userID", uid).Msg("GetBookmark error")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	ctx.JSON(http.StatusOK, &dto.SuccessResponse[[]*model.Bookmark]{
		Data: res.Bookmarks,
		Pagination: &dto.Pagination{
			Page:  input.Page,
			Limit: input.Limit,
			Total: res.Total,
		},
	})
}
