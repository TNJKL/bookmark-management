package handler

import (
	"errors"
	"net/http"

	"github.com/TNJKL/bookmark-management/internal/repository/urlstorage"
	"github.com/TNJKL/bookmark-management/internal/service"
	"github.com/TNJKL/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ShortenURL represents the handler interface for shorten URL related operations,
// including creating and redirecting shortened links.
type ShortenURL interface {
	ShortenLink(ctx *gin.Context)
	Redirect(ctx *gin.Context)
}

// shortenURL is the default implementation of the ShortenURL interface.
type shortenURL struct {
	service service.ShortenURL
}

// shortenInputBody is the expected JSON request body for creating a shortened URL.
type shortenInputBody struct {
	Url string `json:"url" binding:"required,url"`
	Exp int64  `json:"exp" binding:"required,gte=300"`
}

// NewShortenURL creates a new ShortenURL handler backed by the given service.
func NewShortenURL(svc service.ShortenURL) ShortenURL {
	return &shortenURL{
		service: svc,
	}
}

// ShortenLink Generate shorten link
// @Sumary Generate shorten url based on original url that last upto 7 days
// @Description Generate shorten url based on original url that last upto 7 days
// @Tags link
// @Accept application/json
// @Produce application/json
// @Param input body shortenInputBody true "Input required"
// @Success 200 {object} map[string]string
// @Router /v1/links/shorten [post]
func (s *shortenURL) ShortenLink(ctx *gin.Context) {
	//read input
	// Dùng ctx để bind request vào input (struct shortenInputBody)
	// Và struct đó sẽ parse dữ liệu đúng như vậy
	input := &shortenInputBody{}
	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	//call service to create shorten url
	//tuyệt đối : ko nên để "error" : err.Error() trả về cho client vì nó sẽ làm lộ thông tin của Server
	code, err := s.service.CreateShortenLink(ctx, input.Url, input.Exp)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.shortenURL.ShortenLink").Msg("Can't create shorten url")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	//return respone
	ctx.JSON(http.StatusOK, gin.H{"code": code, "message": "Shorten URL generated successfully!"})
}

// Redirect Forward the request to the original url
// @Tags link
// @Accept application/json
// @Produce application/json
// @Param code path string true "Shorten code"
// @Success 302
// @Router /v1/links/redirect/{code} [get]
func (s *shortenURL) Redirect(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, response.InputErrResponse)
		return
	}

	url, err := s.service.GetLinkFromCode(ctx, code)
	if err != nil {
		if errors.Is(err, urlstorage.ErrorCodeNotFound) {
			ctx.JSON(http.StatusNotFound, response.InputErrResponse)
			return
		}

		log.Error().Err(err).Str("from", "handler.shortenURL.Redirect").Msg("Can't get url from code")
		ctx.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	ctx.Redirect(http.StatusFound, url)

}
