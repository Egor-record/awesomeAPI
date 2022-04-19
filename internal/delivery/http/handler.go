package http

import (
	"awesomeAPI/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	eventService service.EventService
}

func newHandler(eventService service.EventService) *Handler {
	return &Handler{
		eventService: eventService,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "200",
		})
	})

	h.InitAPI(router)

	return router
}

func (h *Handler) InitAPI(router *gin.Engine) {
	handler := newHandler(h.eventService)
	v1 := router.Group("/v1")
	{
		v1.POST("/start", handler.eventService.StartEvent)
		v1.POST("/finish", handler.eventService.EndEvent)
	}
}
