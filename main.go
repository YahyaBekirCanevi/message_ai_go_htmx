package main

import (
	"fmt"
	"text/template"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles("simple.html"))
		err := tmpl.Execute(c.Writer, gin.H{})
		if err != nil {
			panic(err)
		}
	})
	r.Run(":8080")
	fmt.Println("Server is running on port 8080")
}
