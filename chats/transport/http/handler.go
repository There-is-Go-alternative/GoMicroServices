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

// GetChatsHandler return the handler responsible for fetching all chats
func (a Handler) GetChatsHandler(cmd usecase.GetAllChatsCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := cmd(c.Request.Context())

		if err != nil {
			a.logger.Error().Msg("Error in GET /chats")
			ResponseError(c, http.StatusInternalServerError, InternalServerError)
			return
		}
		ResponseSuccess(c, http.StatusOK, payload)
	}
}

// GetChatsByIDHandler return the handler responsible for fetching a specific ad.
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
			ResponseError(c, http.StatusBadRequest, "No ad found")
			return
		}
		ResponseSuccess(c, http.StatusOK, payload)
	}
}

// CreateChatHandler return the handler responsible for creating an ad.
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

		var ad usecase.CreateChatInput
		err = c.BindJSON(&ad)
		ad.UserId = string(account.ID)
		if err != nil {
			a.logger.Error().Msgf("User CreateChatInput invalid: %v", ad)
			ResponseError(c, http.StatusBadRequest, FieldsBadRequest)
			return
		}
		payload, err := cmd(c.Request.Context(), ad)
		if err != nil {
			a.logger.Error().Msgf("Error in POST create ad: %v", err)
			ResponseError(c, http.StatusInternalServerError, InternalServerError)
			return
		}
		ResponseSuccess(c, http.StatusCreated, payload)
	}
}

// UpdateChatHandler return the handler responsible for updating an ad.
func (a Handler) UpdateChatHandler(cmd usecase.UpdateChatCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("UpdateChatHandler: param ID missing.")
			ResponseError(c, http.StatusBadRequest, MissingIDParam)
			return
		}
		var ad usecase.UpdateChatInput
		err := c.BindJSON(&ad)
		ad.ID = domain.ChatID(id)

		if err != nil {
			a.logger.Error().Msgf("User UpdateChatInput invalid: %v", ad)
			ResponseError(c, http.StatusBadRequest, FieldsBadRequest)
			return
		}
		payload, err := cmd(c.Request.Context(), ad)
		if err != nil {
			a.logger.Error().Msgf("Error in PATCH update ad: %v", err)
			ResponseError(c, http.StatusInternalServerError, InternalServerError)
			return
		}
		ResponseSuccess(c, http.StatusOK, payload)
	}
}

// DeleteChatHandler return the handler responsible for deleting an ad.
func (a Handler) DeleteChatHandler(cmd usecase.DeleteChatCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("DeleteChatHandler: param ID missing.")
			ResponseError(c, http.StatusBadRequest, MissingIDParam)
			return
		}

		payload, err := cmd(c.Request.Context(), usecase.DeleteChatInput{ID: domain.ChatID(id)})
		if err != nil {
			a.logger.Error().Msgf("Error in POST delete ad: %v", err)
			ResponseError(c, http.StatusInternalServerError, InternalServerError)
			return
		}
		ResponseSuccess(c, http.StatusAccepted, payload)
	}
}

// DeleteChatHandler return the handler responsible for searcgubg an ad.
func (a Handler) SearchChatHandler(cmd usecase.SearchChatCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		content := c.Query("content")
		if content == "" {
			a.logger.Error().Msg("SearchChatHandler: param content missing.")
			ResponseError(c, http.StatusBadRequest, MissingQueryContent)
			return
		}

		payload, err := cmd(c.Request.Context(), usecase.SearchChatInput{Content: content})
		if err != nil {
			a.logger.Error().Msgf("Error in GET search ad: %v", err)
			ResponseError(c, http.StatusInternalServerError, InternalServerError)
			return
		}
		ResponseSuccess(c, http.StatusOK, payload)
	}
}
