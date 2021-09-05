package http

import (
	"axie-notify/models"
	"axie-notify/services"
	"context"

	"log"
	"time"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

type HTTPCallBackHanlder struct {
	Bot          *linebot.Client
	ServicesInfo *models.ServicesInfo
}

// NewServiceHTTPHandler provide the inititail set up service path to handle request
func NewServiceHTTPHandler(e *echo.Echo, linebot *linebot.Client, servicesInfo *models.ServicesInfo) {

	hanlders := &HTTPCallBackHanlder{Bot: linebot, ServicesInfo: servicesInfo}
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Service is online")
	})
	e.POST("/callback", hanlders.Callback)

}

// Callback provides the function to handle request from line
func (handler *HTTPCallBackHanlder) Callback(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	events, err := handler.Bot.ParseRequest(c.Request())
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.String(400, linebot.ErrInvalidSignature.Error())
		} else {
			c.String(500, "internal")
		}
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				messageFromPing := services.PingService(message.Text, handler.ServicesInfo, time.Second*5)
				if _, err = handler.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageFromPing)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
	return c.JSON(200, "")
}
