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

// /chat/new
func NewChatForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("NewChatForm handler called")

		funcMap := template.FuncMap{"markdown": markdownHTML}
		tmpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(
			"templates/chat_window.html",
			"templates/chat_input.html",
			"templates/message_bubble.html",
		))

		err := tmpl.ExecuteTemplate(c.Writer, "chat_window.html", gin.H{"IsNew": true})
		if err != nil {
			log.Printf("Template execution error: %v", err)
		}
		c.Status(http.StatusOK)
	}
}

// /chat/start
func StartChat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		topic := strings.TrimSpace(c.PostForm("topic"))
		if topic == "" {
			c.String(http.StatusBadRequest, "Topic cannot be empty")
			return
		}

		conv := &models.Conversation{}
		_, err := conv.CreateChat(db, topic)
		if err != nil {
			log.Printf("Failed to start conversation: %v", err)
			c.String(http.StatusInternalServerError, "Error creating chat")
			return
		}
		conversations, err := models.GetAllChats(db)
		if err != nil {
			log.Printf("Failed to retrieve conversations: %v", err)
			c.String(http.StatusInternalServerError, "Error retrieving conversations")
			return
		}

		funcMap := template.FuncMap{"markdown": markdownHTML}
		tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

		err = tmpl.ExecuteTemplate(c.Writer, "body.html", gin.H{
			"ConversationID": conv.ID,
			"Conversations":  conversations,
			"Messages":       []models.Message{},
		})

		if err != nil {
			log.Printf("Template execution error: %v", err)
		}
	}
}

// /chat/:id
func LoadChat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		conv, err := models.FindConversationByID(db, id)
		if err != nil {
			c.String(http.StatusNotFound, "Conversation not found")
			return
		}

		conversations, err := models.GetAllChats(db)
		if err != nil {
			log.Printf("Failed to retrieve conversations: %v", err)
			c.String(http.StatusInternalServerError, "Error retrieving conversations")
			return
		}

		msgs, err := models.GetMessagesByConversation(db, id)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error retrieving messages")
			return
		}

		funcMap := template.FuncMap{"markdown": markdownHTML}
		tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

		err = tmpl.ExecuteTemplate(c.Writer, "body.html", gin.H{
			"ConversationID": id,
			"Conversations":  conversations,
			"Messages":       msgs,
			"Title":          conv.Title,
		})

		if err != nil {
			log.Printf("Template execution error: %v", err)
		}
	}
}
