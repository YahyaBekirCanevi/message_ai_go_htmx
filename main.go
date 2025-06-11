package main

import (
	"fmt"
	"text/template"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		users := GetUsers()
		tmpl := template.Must(template.ParseFiles("users.html"))
		err := tmpl.Execute(c.Writer, gin.H{"Title": "User List", "Users": users})
		if err != nil {
			panic(err)
		}
	})
	r.POST("/users", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles("user_row.html"))

		name := c.PostForm("name")
		email := c.PostForm("email")

		user := User{Name: name, Email: email}
		err := tmpl.Execute(c.Writer, user)
		if err != nil {
			panic(err)
		}
	})
	r.DELETE("/users/:name", func(c *gin.Context) {
		name := c.Param("name")
		fmt.Println("Delete user with name:", name)
	})
	r.Run(":8080")
	fmt.Println("Server is running on port 8080")
}

type User struct {
	Name  string
	Email string
}

func GetUsers() []User {
	return []User{
		{Name: "John Doe", Email: "johndoe@example.com"},
		{Name: "Alice Smith", Email: "alicesmith@example.com"},
	}
}
