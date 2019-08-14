package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	"time"
)

var serverStartTime time.Time
var serverConfigFile = "server_config.json"

func main() {
	serverStartTime = time.Now()

	// global serverConfig variable
	var scErr error = nil
	serverConfig, scErr = readServerConfig(serverConfigFile)
	if scErr != nil {
		fmt.Printf("Could not load server configuration\n")
		panic(scErr)
	}
	initDefaultServerConfig(serverConfig)

	// logging setup
	var loggingErr error = nil
	logging = loggingInitSetup(serverConfig)
	if loggingErr = loggingRegisterModules(logging, logModEnabledTable); loggingErr != nil {
		fmt.Printf("Could not register logging modules\n")
		panic(loggingErr)
	}
	if loggingErr = loggingErrRedirect(errLogFile); loggingErr != nil {
		fmt.Printf("Could not redirect logging error to %s\n", errLogFile)
		panic(loggingErr)
	}
	logging.Infomln(logModMain, "Logging module loaded.")

	// db setup
	var dbErr error = nil
	dbPool, dbErr = dbPoolInit(serverConfig)
	if dbErr != nil {
		logging.Panicmf(logModMain, "Unable to load DB - Error msg: %s", dbErr.Error())
	}
	defer func() {
		if dbPool != nil && dbPool.Client() != nil {
			if dbDisconnectErr := dbPool.Client().Disconnect(context.TODO()); dbDisconnectErr != nil {
				logging.Errormln(logModMain, "Unable to disconnected DB: ", dbDisconnectErr)
			} else {
				logging.Infomln(logModMain, "Disconnected from DB")
			}
		}
	}()

	// gin web framework
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = ljGinLogger
	//gin.DefaultErrorWriter = ljErrLogger

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

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

	r.POST("/ping", func(c *gin.Context) {
		var response = GinResponse{
			Status: http.StatusOK,
		}
		defer func() {
			ginContextProcessResponse(c, &response)
		}()

		var params = ginContextRequestParameter(c)
		response.Payload = params
		logging.Infoln(params)
	})

	for apiURL, apiHandlerTable := range ginAPITable {
		for apiMethod, apiHandler := range apiHandlerTable {
			switch apiMethod {
			case "get":
				r.GET(apiURL, apiHandler)
			case "post":
				r.POST(apiURL, apiHandler)
			case "put":
				r.PUT(apiURL, apiHandler)
			case "delete":
				r.DELETE(apiURL, apiHandler)
			}
		}
	}

	logging.Infomln(logModMain, "Server is listening and serving on 0.0.0.0:8080")
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
	logging.Warnmln(logModMain, "Server existed unexpectedly :(")

}
