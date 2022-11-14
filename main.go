package main

import (
	"bannerhelps/internal"
	"bannerhelps/internal/client/pdf"
	"bannerhelps/internal/client/voice"
	ep "bannerhelps/internal/handlers"
	"bannerhelps/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-cleanhttp"
	"net/http"
	"time"
)

func main() {
	fmt.Println(time.Now())

	var vo voice.Client
	{
		vo = voice.NewClient(cleanhttp.DefaultPooledClient())
	}

	var pf pdf.Client
	{
		pf = pdf.NewClient(cleanhttp.DefaultPooledClient())
	}

	var s internal.Service
	{
		s = service.NewService("", vo, pf)
	}

	router := gin.Default()
	router.LoadHTMLGlob("internal/transport/host/*.html")
	router.Static("/static", "./internal/transport/host")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Index.html", gin.H{})
	})

	router.POST("/pdf-to-voice", func(c *gin.Context) {
		ep.Handlers.PDFtoVoiceHandler(s, c)
	})

	err := router.Run(":9090")
	if err != nil {
		return
	}
}
