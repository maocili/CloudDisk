package main

import (
	"CloudDisk/internal/router"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func disk() http.Handler {
	f, _ := os.Create(fmt.Sprintf("disk_%v.log", time.Now().Format("2006-01-02")))
	gin.DefaultWriter = io.MultiWriter(f)
	gin.DefaultErrorWriter = io.MultiWriter(f)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	//注册路由
	router.DiskRouter(r)
	return r
}

func main() {

	serverDisk := &http.Server{
		Addr:         ":9000",
		Handler:      disk(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return serverDisk.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
