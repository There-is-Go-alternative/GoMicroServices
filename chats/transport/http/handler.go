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
	return &Handler{logger: log.With().Str("service", "Chat Http Handler").Logger()}
}

func NewMessageHandler() *Handler {
	return &Handler{logger: log.With().Str("service", "Message Http Handler").Logger()}
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

// GetChatByIDHandler return the handler responsible for fetching a specific chat.
func (a Handler) GetChatByIDHandler(cmd usecase.GetChatByIdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := api.Authorize(c)
		if err != nil {
			ResponseError(c, http.StatusUnauthorized, "You need to be logged in.")
			return
		}
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("GetChatByIDHandler: param ID missing.")
			ResponseError(c, http.StatusBadRequest, MissingIDParam)
			return
		}
		payload, err := cmd(c.Request.Context(), domain.ChatID(id))
		if err != nil {
			a.logger.Error().Msg("Error in GET by /chats/:id")
			ResponseError(c, http.StatusBadRequest, "No chat found")
			return
		}
		checker := false
		for _, current_user := range payload.UsersIDs {
			if current_user == id {
				checker = true
			}
		}
		if checker {
			ResponseSuccess(c, http.StatusOK, payload)
			return
		} else {
			a.logger.Error().Msg("Error in GET by /chats/:id")
			ResponseError(c, http.StatusBadRequest, "User is not in this chat.")
			return
		}
	}
}

// CreateChatHandler return the handler responsible for creating a chat.
func (a Handler) CreateChatHandler(cmd usecase.CreateChatCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := api.Authorize(c)
		if err != nil {
			ResponseError(c, http.StatusUnauthorized, "You need to be logged in.")
			return
		}

		var chat usecase.CreateChatInput
		err = c.BindJSON(&chat)
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
func (a Handler) GetAllChatsOfUserHandler(cmd usecase.GetAllChatsOfUserCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		account, err := api.Authorize(c)
		if err != nil {
			ResponseError(c, http.StatusUnauthorized, "You need to be logged in.")
			return
		}
		id := c.Param("user_id")
		if id == "" {
			a.logger.Error().Msg("GetChatsOfUserHandler: param user_id missing.")
			ResponseError(c, http.StatusBadRequest, MissingIDParam)
			return
		}
		if account.ID.String() != id {
			ResponseError(c, http.StatusUnauthorized, "The user ID given does not correspond to the token's.")
			return
		}
		payload, err := cmd(c.Request.Context(), id)

		if err != nil {
			a.logger.Error().Msg("Error in GET by /chats/:user_id")
			ResponseError(c, http.StatusBadRequest, "No chats found for this user")
			return
		}
		ResponseSuccess(c, http.StatusOK, payload)
	}
}

// GetChatByIDHandler return the handler responsible for fetching a specific chat.
func (a Handler) GetMessagesByChatIDHandler(cmd usecase.GetMessagesByChatIDCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := api.Authorize(c)
		if err != nil {
			ResponseError(c, http.StatusUnauthorized, "You need to be logged in.")
			return
		}
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("GetMessagesByChatID: param ID missing.")
			ResponseError(c, http.StatusBadRequest, MissingIDParam)
			return
		}
		payload, err := cmd(c.Request.Context(), domain.ChatID(id))

		if err != nil {
			a.logger.Error().Msg("Error in GET by /messages/:id")
			ResponseError(c, http.StatusBadRequest, "No message found")
			return
		}
		ResponseSuccess(c, http.StatusOK, payload)
	}
}

// CreateChatHandler return the handler responsible for creating a chat.
func (a Handler) CreateMessageHandler(cmd usecase.CreateMessageCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := api.Authorize(c)
		if err != nil {
			ResponseError(c, http.StatusUnauthorized, "You need to be logged in.")
			return
		}

		var message usecase.CreateMessageInput
		err = c.BindJSON(&message)
		if err != nil {
			a.logger.Error().Msgf("User CreateMessageInput invalid: %v", message)
			ResponseError(c, http.StatusBadRequest, FieldsBadRequest)
			return
		}
		payload, err := cmd(c.Request.Context(), message)
		if err != nil {
			a.logger.Error().Msgf("Error in POST create chat: %v", err)
			ResponseError(c, http.StatusInternalServerError, InternalServerError)
			return
		}
		ResponseSuccess(c, http.StatusCreated, payload)
	}
}
