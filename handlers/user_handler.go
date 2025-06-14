package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/YahyaBekirCanevi/message_ai_go_htmx/models"
	"github.com/gin-gonic/gin"
)

func ListUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := models.GetAllUsers(db)
		if err != nil {
			log.Printf("Failed to retrieve users: %v", err)
			c.String(http.StatusInternalServerError, "Error retrieving users")
			return
		}
		tmpl := template.Must(template.ParseFiles("templates/users.html"))
		tmpl.Execute(c.Writer, gin.H{"Title": "User List", "Users": users})
	}
}

func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		user := models.User{Name: name, Email: email}

		_, err := user.Insert(db)
		if err != nil {
			log.Printf("Failed to insert user: %v", err)
			c.String(http.StatusInternalServerError, "Error inserting user")
			return
		}
		tmpl := template.Must(template.ParseFiles("templates/user_row.html"))
		c.Header("Content-Type", "text/html")
		c.Status(http.StatusCreated)
		c.Writer.Write([]byte(fmt.Sprintf(`<tr hx-swap-oob="true" id="user-%d">`, user.ID)))
		_ = tmpl.Execute(c.Writer, user)
		c.Writer.Write([]byte(`</tr>`))
	}
}

func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Invalid user ID: %v", err)
			c.String(http.StatusBadRequest, "Invalid user ID")
			return
		}

		name := c.PostForm("name")
		email := c.PostForm("email")
		user := models.User{ID: id, Name: name, Email: email}

		_, err = user.Update(db)
		if err != nil {
			log.Printf("Failed to update user: %v", err)
			c.String(http.StatusInternalServerError, "Error updating user")
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/user_row.html"))
		c.Header("Content-Type", "text/html")
		c.Status(http.StatusCreated)
		c.Writer.Write([]byte(fmt.Sprintf(`<tr hx-swap-oob="true" id="user-%d">`, user.ID)))
		_ = tmpl.Execute(c.Writer, user)
		c.Writer.Write([]byte(`</tr>`))
	}
}

func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")
		user, err := models.FindByEmail(db, email)
		if err != nil {
			log.Printf("User not found: %v", err)
			c.String(http.StatusNotFound, "User not found")
			return
		}
		_, err = user.Delete(db)
		if err != nil {
			log.Printf("Failed to delete user: %v", err)
			c.String(http.StatusInternalServerError, "Error deleting user")
			return
		}
		c.String(http.StatusOK, "")
	}
}

func RenderNewUserModal() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := struct {
			IsEdit bool
			User   models.User
		}{IsEdit: false, User: models.User{}}

		tmpl := template.Must(template.ParseFiles("templates/user_modal.html"))
		tmpl.Execute(c.Writer, data)
	}
}

func RenderEditUserModal(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Invalid user ID: %v", err)
			c.String(http.StatusBadRequest, "Invalid user ID")
			return
		}

		user, err := models.FindByID(db, id)
		if err != nil {
			log.Printf("User not found: %v", err)
			c.String(http.StatusNotFound, "User not found")
			return
		}

		data := struct {
			IsEdit bool
			User   models.User
		}{IsEdit: true, User: *user}

		tmpl := template.Must(template.ParseFiles("templates/user_modal.html"))
		tmpl.Execute(c.Writer, data)
	}
}
