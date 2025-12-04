package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redundant4u/SeolMyeongTang-Server/chat/handlers"
)

func InitRouter(r *gin.Engine) {
	r.GET("/", handlers.FindChats)
}
