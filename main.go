package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"gopkg.in/natefinch/lumberjack.v2" // go get -u gopkg.in/natefinch/lumberjack.v2
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/debug"
	"time"
)

var counter uint64
var startTime time.Time

const LOG_FILE_NAME string = "/tmp/vlog_server.log"

var LJ_LOGGER = &lumberjack.Logger{
	Filename:   LOG_FILE_NAME,
	MaxSize:    80, // max file size is 80M
	MaxBackups: 10,
}
var PROFILE_FILE_NAME string = "/tmp/cpuprofile"
var flagCpuprofile = &PROFILE_FILE_NAME
var MEMPROFILE_FILE_NAME string = "/tmp/memprofile"
var EMPTY_STR string = ""
var flagMemprofile = &EMPTY_STR

func Wrap(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request)
	}
}

func checkFileSize(file string) int64 {
	fi, e := os.Stat(file)
	if e != nil {
		return 0
	}
	// get the size
	size := fi.Size()
	return size
}

func initCrashLog(file string) {
	if checkFileSize(file) > 1024*1024 {
		os.Rename(file, file+".old")
	}
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		log.Errorf("Failed to seek log file to end: %v", err)
	}
	// redirect stderr
	err = unix.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
}

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	// global serverConfig variable
	if len(os.Args) == 1 {
		serverConfig = ReadServerConfig("config.json")
	} else {
		serverConfig = ReadServerConfig("test.json")
	}

	InitDefServerConfig()

	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = LJ_LOGGER
	initCrashLog("/tmp/fcl_crash.log")

	startTime = time.Now()
	//r := gin.Default()
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
	r.Run() // listen and serve on 0.0.0.0:8080
}
