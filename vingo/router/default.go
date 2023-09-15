package router

import (
	"fmt"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"github.com/lgdzz/vingo-utils/vingo"
	"github.com/lgdzz/vingo-utils/vingo/config"
	"net/http"
	"time"
)

type Hook struct {
	Config         config.Config
	RegisterRouter func(r *gin.Engine)
	BaseMiddle     func(c *gin.Context)
	LoadWeb        []WebItem
}

type WebItem struct {
	Route string
	FS    *assetfs.AssetFS
}

// 初始化路由
func InitRouter(hook *Hook) {
	var option = hook.Config
	if option.System.Service.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	_ = r.SetTrustedProxies(nil)

	// 加载web前端
	for _, item := range hook.LoadWeb {
		currentItem := item
		r.GET(currentItem.Route, func(c *gin.Context) {
			http.FileServer(currentItem.FS).ServeHTTP(c.Writer, c.Request)
		})
	}

	// 屏蔽搜索引擎爬虫
	vingo.ShieldRobots(r)

	// 404捕捉
	r.NoRoute(func(c *gin.Context) {
		context := vingo.Context{c}
		context.Response(&vingo.ResponseData{Message: "404:Not Found", Error: 1})
	})

	// 注册异常处理、基础中间件
	r.Use(vingo.ExceptionHandler, BaseMiddle(hook))

	// 注册路由
	hook.RegisterRouter(r)

	fmt.Println("+------------------------------------------------------------+")
	fmt.Println(fmt.Sprintf("+ 项目名称：%v", option.System.Service.Name))
	fmt.Println(fmt.Sprintf("+ 服务端口：%d", option.System.Service.Port))
	fmt.Println(fmt.Sprintf("+ 调试模式：%v", option.System.Service.Debug))
	fmt.Println(fmt.Sprintf("+ Mysql：%v:%v db:%v", option.Database.Host, option.Database.Port, option.Database.Dbname))
	fmt.Println(fmt.Sprintf("+ Redis：%v:%v db:%v", option.Redis.Host, option.Redis.Port, option.Redis.Select))
	vingo.ApiAddress(option.System.Service.Port)
	fmt.Println(fmt.Sprintf("+ 启动时间：%v", time.Now().Format(vingo.DatetimeFormatChinese)))
	fmt.Println(fmt.Sprintf("+ 技术支持：", option.System.Service.Copyright))
	fmt.Println("+------------------------------------------------------------+")

	// 开启服务
	_ = r.Run(fmt.Sprintf(":%d", option.System.Service.Port))
}
