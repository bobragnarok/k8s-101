package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/api/ping", ping)
	e.GET("/api/health", health)
	e.GET("/api/heat", heat)
	e.GET("/api/file", file)

	go startServer(e)

	waitForGracefulShutdown(e)
}

func ping(c echo.Context) error {
	env := os.Getenv("ENVIRONMENT")
	return c.String(http.StatusOK, env)
}

func health(c echo.Context) error {
	os.Getenv("ENVIRONMENT")
	return c.String(http.StatusOK, "OK")
}

func heat(c echo.Context) error {
	mapI := make(map[int]string)
	for i := 0; i <= 1000000; i++ {
		mapI[i] = strconv.Itoa(i)
	}
	return c.String(http.StatusOK, "OK")
}

func file(c echo.Context) error {
	dateStr := time.Now().Format("20060102150405")
	if _, err := os.Stat("file"); os.IsNotExist(err) {
		if err := os.Mkdir("file", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	f, err := os.Create(fmt.Sprintf("file/%s.txt", dateStr))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	return c.String(http.StatusOK, "OK")
}

func startServer(e *echo.Echo) {
	if err := e.Start(":8080"); err != nil {
		log.Info("shutting down the server")
	}
}

func waitForGracefulShutdown(e *echo.Echo) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
