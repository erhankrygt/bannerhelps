package handlers

import (
	"bannerhelps"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	PDF = ".pdf"
)

var Handlers = struct {
	PDFtoVoiceHandler func(s bannerhelps.Service, c *gin.Context)
}{
	PDFtoVoiceHandler: func(s bannerhelps.Service, c *gin.Context) {
		ctx := context.Background()

		pdf, err := c.FormFile("file")

		if err != nil {
			d := gin.H{
				"message": "No file is received",
			}
			c.HTML(http.StatusOK, "Index.html", d)
			return
		}

		extension := filepath.Ext(pdf.Filename)
		if extension != PDF {
			d := gin.H{
				"message": "Extension should be pdf",
			}
			c.HTML(http.StatusOK, "Index.html", d)
			return
		}

		fHeader, _ := pdf.Open()
		file, err := ioutil.ReadAll(fHeader)
		if err != nil {
			d := gin.H{
				"message": "No file is received",
			}
			c.HTML(http.StatusOK, "Index.html", d)
			return
		}

		fName := strings.ReplaceAll(pdf.Filename, extension, "")

		req := bannerhelps.PdfToVoiceRequest{
			CurrentLanguage: c.PostForm("documentLanguage"),
			File:            file,
			FileName:        fName,
		}
		res := s.PdfToVoice(ctx, req)

		f, _ := os.Open(res.FilePath)
		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fName+".mp3")
		c.Writer.Header().Set("Content-Type", c.Request.Header.Get("Content-Type"))

		_, _ = io.Copy(c.Writer, f)

		d := gin.H{
			"success": res.IsSuccessfully,
		}
		c.HTML(http.StatusOK, "Index.html", d)
	},
}
