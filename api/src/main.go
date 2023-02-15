package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	worker sync.WaitGroup
)

func main() {
	router := gin.Default()
	router.GET("/hi", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello world!")
	})
	err := router.Run(":9888")
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}

// 	queue_download.Add(1)
// 	go download(link, data, &queue_download)
//  defer queue_download.Done()
// 	queue_download.Wait()
