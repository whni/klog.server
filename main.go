package main

import (
	_ "net/http/pprof"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2" // go get -u gopkg.in/natefinch/lumberjack.v2
)

var serverStartTime time.Time

const (
	stdLogFile = "/tmp/vlog_server.log"
	errLogFile = "/tmp/vlog_crash.log"
)

var ljLogger = &lumberjack.Logger{
	Filename:   stdLogFile,
	MaxSize:    80, // max file size is 80M
	MaxBackups: 10,
}

var ljErrLogger = &lumberjack.Logger{
	Filename:   errLogFile,
	MaxSize:    80, // max file size is 80M
	MaxBackups: 10,
}

func main() {
	serverStartTime = time.Now()

	// global serverConfig variable
	serverConfig = ReadServerConfig("config.json")
	InitDefServerConfig()

	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = ljLogger
	gin.DefaultErrorWriter = ljErrLogger

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	//	log.SetOutput(gin.DefaultWriter)
	//	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	debug.SetTraceback("crash")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	/*
			for api, handler := range Api_maps {
				func_map := handler
				r.POST(api, func_map["Append"])
				r.GET(api, func_map["Get"])
				r.PUT(api, func_map["Change"])
				r.DELETE(api, func_map["Delete"])
			}

			for api, handler := range Api_controller_post {
				r.POST(api, handler)
			}

			for api, handler := range Api_history_post {
				r.POST(api, handler)
			}

			for api, handler := range Api_cwconf_post {
				r.POST(api, handler)
			}

			for api, handler := range Api_status_get {
				r.GET(api, handler)
			}

			for api, handler := range Api_sys_get {
				r.GET(api, handler)
			}

		DbInit()
	*/

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
