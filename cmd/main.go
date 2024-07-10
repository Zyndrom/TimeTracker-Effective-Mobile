package main

import (
	"fmt"
	"path"
	"runtime"

	_ "github.com/TimeTracker-Effective-Mobile/internal/model"
	"github.com/TimeTracker-Effective-Mobile/internal/repository"
	"github.com/TimeTracker-Effective-Mobile/internal/router"
	"github.com/TimeTracker-Effective-Mobile/internal/service/task"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

// swagger embed files
//	@title			Time Tracker
//	@version		1.0
//	@description	Time Tracker.

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/

func main() {
	setLogger()
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf(".env file not found.")
	}
	repository := repository.New()
	taskService := task.New(repository)
	router := router.New(taskService)
	router.StartServer()
}

func setLogger() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	formatter := &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}
	logrus.SetFormatter(formatter)
}
