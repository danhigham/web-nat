package main

import (
//	"net/http"
	"github.com/gin-gonic/gin"
)


func main() {
	r := gin.Default()
	r.Static("/assets", "./assets")

	r.GET("/shipyard_containers", func(c *gin.Context) {
		res := TestData{}.Containers()
		writeStringAsJSON(c, 200, res)
	})

	r.GET("/ip_table", func(c *gin.Context) {
		res := TestData{}.IPTable()
		writeStringAsJSON(c, 200, res)
	})

	r.Run(":8080")
}

func writeStringAsJSON(c *gin.Context, status int, body string) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(status)
	c.Writer.Write([]byte(body))
	c.Writer.Flush()
}
