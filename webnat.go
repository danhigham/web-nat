package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


func main() {
	r := gin.Default()
	r.Static("/assets", "./assets")

	r.GET("/ping", func(c *gin.Context) {
	    c.String(http.StatusOK, "pong")
	})

	r.Run(":8080")
}
