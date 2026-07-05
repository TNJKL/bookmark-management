package handler

import (
	"net/http"
	"time"

	"github.com/TNJKL/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

// Shorten URl represents the Handler for shorten url related operations
// including creating and redirecting shortened links.
type ShortenURL interface {
	ShortenLink(ctx *gin.Context)
	Redirect(ctx *gin.Context)
}

// shortenURL is the default implementation of the ShortenURL interface
type shortenURL struct {
	service service.ShortenURL
}

// shortenInputBody is the expected JSON request body for creating a shortened URL
type shortenInputBody struct {
	Url string `json:"url"`
	Exp int64  `json:"exp"`
}

// NewShortenURL creates a new ShortenURL handler backed by the given service
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	//call service to create shorten url
	//tuyệt đối : ko nên để "error" : err.Error() trả về cho client vì nó sẽ làm lộ thông tin của Server
	code, err := s.service.CreateShortenLink(ctx, input.Url, time.Duration(input.Exp))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	//return respone
	ctx.JSON(http.StatusOK, gin.H{"code": code})
}

func (s *shortenURL) Redirect(ctx *gin.Context) {}
