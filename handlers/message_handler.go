package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/YahyaBekirCanevi/message_ai_go_htmx/models"
	"github.com/gin-gonic/gin"
)

func SendMessage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		conversationIDStr := c.PostForm("conversation_id")
		content := strings.TrimSpace(c.PostForm("message"))

		if content == "" || conversationIDStr == "" {
			c.String(http.StatusBadRequest, "Message or conversation ID missing")
			return
		}

		conversationID, err := strconv.Atoi(conversationIDStr)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid conversation ID")
			return
		}

		// Create user message
		userMsg, err := models.CreateMessage(db, conversationID, "user", content)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to save user message")
			return
		}

		// Create dummy AI reply (can be replaced with real model later)
		aiReply := "This is a placeholder AI response."
		aiMsg, err := models.CreateMessage(db, conversationID, "ai", aiReply)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to save AI message")
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/message_bubble.html"))

		c.Header("Content-Type", "text/html")
		tmpl.Execute(c.Writer, userMsg)
		tmpl.Execute(c.Writer, aiMsg)
	}
}

func ListConversations(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		conversations, err := models.GetAllChats(db)
		if err != nil {
			log.Printf("Failed to retrieve conversations: %v", err)
			c.String(http.StatusInternalServerError, "Error retrieving conversations")
			return
		}

		var messages []models.Message

		tmpl := template.Must(template.ParseGlob("templates/*.html"))

		err = tmpl.ExecuteTemplate(c.Writer, "layout.html", gin.H{
			"Conversations":  conversations,
			"Messages":       messages,
		})
		if err != nil {
			log.Printf("Template execution error: %v", err)
		}
	}
}
