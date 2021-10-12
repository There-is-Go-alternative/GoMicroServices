package http

import (
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/chats/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/chats/transport/api"
	"github.com/There-is-Go-alternative/GoMicroServices/chats/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	InternalServerError string = "Internal server error"
	MissingIDParam      string = "Missing id parameter"
	MissingQueryContent string = "Missing query content"
	FieldsBadRequest    string = "One or more fields are not correct"
)

type Handler struct {
	logger zerolog.Logger
}

func NewChatHandler() *Handler {
	return &Handler{logger: log.With().Str("service", "Http Handler").Logger()}
}

func ResponseError(c *gin.Context, code int, message interface{}) {
	c.JSON(code, gin.H{
		"success": false,
		"message": message,
	})
}

func ResponseSuccess(c *gin.Context, code int, message interface{}) {
	c.JSON(code, gin.H{
		"success": true,
		"data":    message,
	})
}

// GetChatsByIDHandler return the handler responsible for fetching a specific chat.
func (a Handler) GetChatsByIDHandler(cmd usecase.GetChatByIdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("GetChatsByIDHandler: param ID missing.")
			ResponseError(c, http.StatusBadRequest, MissingIDParam)
			return
		}
		payload, err := cmd(c.Request.Context(), domain.ChatID(id))

		if err != nil {
			a.logger.Error().Msg("Error in GET by /chats/:id")
			ResponseError(c, http.StatusBadRequest, "No chat found")
			return
		}
		ResponseSuccess(c, http.StatusOK, payload)
	}
}

// CreateChatHandler return the handler responsible for creating a chat.
func (a Handler) CreateChatHandler(cmd usecase.CreateChatCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		account, err := api.Authorize(c)
		//TODO fix error
		if err != nil {
			//TODO encapsulate function
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "You need to be logged in",
			})
			return
		}

		var chat usecase.CreateChatInput
		err = c.BindJSON(&chat)
		chat.UserId = string(account.ID)
		if err != nil {
			a.logger.Error().Msgf("User CreateChatInput invalid: %v", chat)
			ResponseError(c, http.StatusBadRequest, FieldsBadRequest)
			return
		}
		payload, err := cmd(c.Request.Context(), chat)
		if err != nil {
			a.logger.Error().Msgf("Error in POST create chat: %v", err)
			ResponseError(c, http.StatusInternalServerError, InternalServerError)
			return
		}
		ResponseSuccess(c, http.StatusCreated, payload)
	}
}

// GetChatsOfUseHandler return the handler responsible for fetching a user's chats.
func (a Handler) GetChatsOfUserHandler(cmd usecase.GetChatByIdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("user_id")
		if id == "" {
			a.logger.Error().Msg("GetChatsOfUserHandler: param user_id missing.")
			ResponseError(c, http.StatusBadRequest, MissingIDParam)
			return
		}
		payload, err := cmd(c.Request.Context(), domain.ChatID(id))

		if err != nil {
			a.logger.Error().Msg("Error in GET by /chats/:user_id")
			ResponseError(c, http.StatusBadRequest, "No chats found for this user")
			return
		}
		ResponseSuccess(c, http.StatusOK, payload)
	}
}
