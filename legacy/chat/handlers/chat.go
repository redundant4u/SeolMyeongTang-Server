package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redundant4u/SeolMyeongTang-Server/chat/models"
)

func FindChats(c *gin.Context) {
	var chats []models.Chat
	models.DB.Limit(100).Find(&chats)

	c.JSON(http.StatusOK, gin.H{"data": chats})
}
