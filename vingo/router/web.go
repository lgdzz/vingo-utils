package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ServeStaticFile(c *gin.Context, fs http.FileSystem) {
	http.FileServer(fs).ServeHTTP(c.Writer, c.Request)
}

func LoadClient(r *gin.Engine) {
	fs := assetFS()

	staticRoutes := []string{
		"/",
		"/assets/*filepath",
		"/css/*filepath",
		"/static/*filepath",
		"/ueditor/*filepath",
		"/favicon.ico",
		"/loading.gif",
	}

	for _, route := range staticRoutes {
		r.GET(route, func(c *gin.Context) {
			ServeStaticFile(c, fs)
		})
	}
}
