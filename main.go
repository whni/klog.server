package main

import (
	"context"
	"fmt"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var serverStartTime time.Time
var serverConfigFile = "server_config.json"

func main() {
	serverStartTime = time.Now()

	// global serverConfig variable
	var scErr error
	serverConfig, scErr = readServerConfig(serverConfigFile)
	if scErr != nil {
		fmt.Printf("Could not load server configuration\n")
		panic(scErr)
	}
	initDefaultServerConfig(serverConfig)

	// logging setup
	var loggingErr error
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
	var dbErr error
	dbPool, dbErr = dbPoolInit(serverConfig)
	if dbErr != nil {
		logging.Panicmf(logModMain, "Unable to load DB - Error msg: %s", dbErr.Error())
	}
	logging.Infomln(logModMain, "DB module loaded.")
	defer func() {
		if dbPool != nil && dbPool.Client() != nil {
			if dbDisconnectErr := dbPool.Client().Disconnect(context.TODO()); dbDisconnectErr != nil {
				logging.Errormln(logModMain, "Unable to disconnected DB: ", dbDisconnectErr)
			} else {
				logging.Infomln(logModMain, "Disconnected from DB")
			}
		}
	}()

	// azure storage setup
	var azureContainerErr error
	azMediaContainerURL, azureContainerErr = azureStorageInit(serverConfig)
	if azureContainerErr != nil {
		logging.Panicmf(logModMain, "Unable to load azure storage container - Error msg: %s", azureContainerErr.Error())
	}
	logging.Infomln(logModMain, "Azure storage container loaded.")

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

	// regiester config api handlers
	for apiURL, apiHandlerTable := range ginConfigAPITable {
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
	// register workflow api handlers
	for apiURL, apiHandler := range ginWorkflowAPITable {
		r.POST(apiURL, apiHandler)
	}

	if serverConfig.RunHTTPS {
		logging.Infomf(logModMain, "HTTPS Server is listening on port %d", serverConfig.ServerHTTPSecurePort)
		r.RunTLS(fmt.Sprintf(":%d", serverConfig.ServerHTTPSecurePort), serverConfig.SSLCertPath, serverConfig.SSLKeyPath)
	} else {
		logging.Infomf(logModMain, "HTTP Server is listening on port %d", serverConfig.ServerHTTPPort)
		r.Run(fmt.Sprintf(":%d", serverConfig.ServerHTTPPort))
	}
	logging.Warnmln(logModMain, "Server existed unexpectedly -> please check server config")
}
