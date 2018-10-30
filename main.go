package main

import (
	"summaryGeneraterApi/summary"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/post", func(c *gin.Context) {
		text := c.PostForm("text")
		delimiter := c.PostForm("delimiter")

		summary.Generate(text, delimiter)
		c.JSON(200, gin.H{
			"summary":   text,
			"delimiter": delimiter,
		})
	})

	router.Run(":8080")
}
