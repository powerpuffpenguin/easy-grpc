package register

import (
	"net/http"
	"server/m/helper"
	"server/m/web"
	"server/m/web/api"
	"server/static"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func HTTP(cc *grpc.ClientConn, engine *gin.Engine, gateway *runtime.ServeMux, swagger bool) {
	var w web.Helper
	engine.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusOK)
		if c.Request.Method == `GET` || c.Request.Method == `HEAD` {
			c.Request.Header.Set(`Method`, c.Request.Method)
		}
		if c.Request.Header.Get(`Authorization`) == `` {
			var query struct {
				AccessToken string `form:"access_token"`
			}
			c.ShouldBindQuery(&query)
			if query.AccessToken != `` {
				c.Request.Header.Set(`Authorization`, `Bearer `+query.AccessToken)
			}
			if query.AccessToken == `` {
				query.AccessToken, _ = c.Cookie(helper.CookieName)
				if query.AccessToken != `` {
					c.Request.Header.Set(`Authorization`, `Bearer `+query.AccessToken)
				}
			}
		}
		gateway.ServeHTTP(c.Writer, c.Request)
	})
	if swagger {
		r := engine.Group(`document`)
		r.Use(w.Compression())
		r.StaticFS(``, static.Document())
	}
	// favicon.ico
	engine.GET(`favicon.ico`, static.Favicon)
	engine.HEAD(`favicon.ico`, static.Favicon)
	// other gin route
	var apis api.Helper
	apis.Register(cc, &engine.RouterGroup)
}
