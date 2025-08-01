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
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
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

		// Fetch previous messages for context
		prevMessages, err := models.GetMessagesByConversation(db, conversationID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to fetch previous messages")
			return
		}

		// Create user message
		err = models.CreateMessage(db, conversationID, "user", content)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to save user message")
			return
		}

		// Get AI reply and updated context prompt
		aiReply, err := GetGeminiAIResponse(prevMessages, content)
		if err != nil {
			log.Printf("Gemini API error: %v", err)
			aiReply = "[AI Error: could not get response]"
		}
		err = models.CreateMessage(db, conversationID, "ai", aiReply)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to save AI message")
			return
		}

		// Fetch all messages for the conversation
		allMessages, err := models.GetMessagesByConversation(db, conversationID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to fetch messages")
			return
		}

		funcMap := template.FuncMap{"markdown": markdownHTML}
		tmpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(
			"templates/chat_body.html",
			"templates/message_bubble.html",
		))

		c.Header("Content-Type", "text/html")
		tmpl.ExecuteTemplate(c.Writer, "chat_body.html", gin.H{
			"Messages": allMessages,
		})
	}
}

func markdownHTML(input string) template.HTML {
	output := markdown.ToHTML([]byte(input), nil, html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank}))
	return template.HTML(output)
}

func ListConversations(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		conversations, err := models.GetAllChats(db)
		if err != nil {
			log.Printf("Failed to retrieve conversations: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving conversations"})
			return
		}

		funcMap := template.FuncMap{"markdown": markdownHTML}
		tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

		var cId *int
		if len(conversations) > 0 {
			cId = &conversations[0].ID
		} else {
			cId = nil
		}

		var messages []models.Message
		if cId != nil {
			messages, err = models.GetMessagesByConversation(db, *cId)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error retrieving messages")
				return
			}
		} else {
			messages = []models.Message{}
		}

		var conversationID interface{}
		if cId != nil {
			conversationID = *cId
		} else {
			conversationID = nil
		}

		err = tmpl.ExecuteTemplate(c.Writer, "layout.html", gin.H{
			"IsNew":          cId == nil,
			"ConversationID": conversationID,
			"Conversations":  conversations,
			"Messages":       messages,
		})
		if err != nil {
			log.Printf("Template execution error: %v", err)
		}
	}
}
