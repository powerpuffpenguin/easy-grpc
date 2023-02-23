package static

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed favicon.ico
var favicon embed.FS

func Favicon(c *gin.Context) {
	c.Header("Cache-Control", "max-age=2419200")
	c.FileFromFS(`favicon.ico`, http.FS(favicon))
}

//go:embed document/*
var document embed.FS

func Document() http.FileSystem {
	f, e := fs.Sub(document, `document`)
	if e != nil {
		panic(e)
	}
	return http.FS(f)
}
