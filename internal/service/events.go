package service

import (
	"awesomeAPI/internal/storage"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"regexp"
	"strings"
)

type EventService struct {
	eventStorage storage.EventStorage
}

type jsonBody struct {
	Type string `json:"type"`
}

func newEventStorage(eventStorage storage.EventStorage) *EventService {
	return &EventService{
		eventStorage: eventStorage,
	}
}

func (s *EventService) StartEvent(c *gin.Context) {

	body, err := s.parseJSON(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if !s.validateRequestType(body.Type) {
		c.JSON(200, gin.H{
			"error": "Invalid type",
		})
		return
	}
	storageDB := newEventStorage(s.eventStorage)
	err = storageDB.eventStorage.StartEvent(body.Type)
	if err != nil {
		c.JSON(200, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "OK",
	})

}

func (s *EventService) EndEvent(c *gin.Context) {

	body, err := s.parseJSON(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !s.validateRequestType(body.Type) {
		c.JSON(200, gin.H{
			"error": "Invalid type",
		})
	}

	storageDB := newEventStorage(s.eventStorage)
	err = storageDB.eventStorage.EndEvent(body.Type)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			c.JSON(404, gin.H{
				"error": "Not found",
			})
			return
		default:
			c.JSON(200, gin.H{
				"error": "Invalid type",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"status": "Event finished",
	})

}

func (s *EventService) parseJSON(c *gin.Context) (*jsonBody, error) {
	body := jsonBody{}
	err := c.BindJSON(&body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func (s *EventService) validateRequestType(req string) bool {
	if req != "" {
		return len(regexp.MustCompile("\\w+").FindAllString(req, -1)) == len(strings.Fields(req))
	}
	return false
}
