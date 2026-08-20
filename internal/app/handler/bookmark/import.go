package bookmark

import (
	"net/http"
	"slices"

	"github.com/TNJKL/bookmark-management/internal/app/service/queue"
	"github.com/TNJKL/bookmark-management/pkg/csv"
	"github.com/TNJKL/bookmark-management/pkg/requestutils"
	"github.com/TNJKL/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	MaxFileSize = 10 << 20 //10MB
)

var allowedFileTypes = []string{"text/csv", "application/csv", "text/plain", "application/octet-stream"}

// ImportBookmarks        handles file uploads
// @Summary               Upload and parse a CSV file of bookmarks
// @Description           Accepts a CSV file and import bookmarks
// @Tags                   bookmark
// @Security               BearerAuth
// @Accept                 multipart/form-data
// @Produce                json
// @Param                  file formData file true "CSV file"
// @Success                200 {object} object{message=string} "Success"
// @Router                 /v1/bookmarks/import [post]
func (h *bookmarkHandler) ImportBookmarks(ctx *gin.Context) {
	//Get uuid
	uid, err := requestutils.GetUserIDFromRequest(ctx)
	if err != nil {
		return
	}

	//Get .csv file in request
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.Message{
			Message: "Invalid file",
		})
		return
	}

	//Check file size
	if file.Size > MaxFileSize {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.Message{
			Message: "File too large",
		})
		return
	}

	//Validate file type
	fileType := file.Header.Get("Content-Type")
	print(fileType)
	if !slices.Contains(allowedFileTypes, fileType) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.Message{
			Message: "Invalid file type",
		})
		return
	}

	//parse file to struct (import message)
	var importInput []*queue.ImportBookmarkInput
	err = csv.ParseFromMultipartFile(file, &importInput)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.Message{
			Message: "Unable to parse file",
		})
		return
	}
	//validate
	err = requestutils.InputValidator.Var(importInput, "dive")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	//sent message to redis queue
	err = h.queue.SendImportBookmarkJob(ctx, uid, importInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send import bookmark job")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	//return success response
	ctx.JSON(http.StatusOK, &response.Message{
		Message: "Imported bookmarks successfully",
	})
}
