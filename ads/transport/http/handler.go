package http

import (
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/transport/api"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	InternalServerError string = "Internal server error"
	MissingIDParam string = "Missing id parameter"
	MissingQueryContent string = "Missing query content"
	FieldsBadRequest string = "One or more fields are not correct"
)

type Handler struct {
	logger zerolog.Logger
}

func NewAdHandler() *Handler {
	return &Handler{logger: log.With().Str("service", "Http Handler").Logger()}
}

func ResponseError(c *gin.Context, code int, message interface {}) {
	c.JSON(code, gin.H {
		"success": false,
		"message": message,
	})
}

func ResponseSuccess(c *gin.Context, code int, message interface {}) {
	c.JSON(code, gin.H {
		"success": true,
		"data": message,
	})
}

// GetAdsHandler return the handler responsible for fetching all ads
func (a Handler) GetAdsHandler(cmd usecase.GetAllAdsCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := cmd(c.Request.Context())

		if err != nil {
			a.logger.Error().Msg("Error in GET /ads")
			ResponseError(c, http.StatusInternalServerError, InternalServerError)
			return
		}
		ResponseSuccess(c, http.StatusOK, payload)
	}
}

// GetAdsByIDHandler return the handler responsible for fetching a specific ad.
func (a Handler) GetAdsByIDHandler(cmd usecase.GetAdByIdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("GetAdsByIDHandler: param ID missing.")
			ResponseError(c, http.StatusBadRequest, MissingIDParam)
			return
		}
		payload, err := cmd(c.Request.Context(), domain.AdID(id))

		if err != nil {
			a.logger.Error().Msg("Error in GET by /ads/:id")
			ResponseError(c, http.StatusBadRequest, "No ad found")
			return
		}
		ResponseSuccess(c, http.StatusOK, payload)
	}
}

// CreateAdHandler return the handler responsible for creating an ad.
func (a Handler) CreateAdHandler(cmd usecase.CreateAdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		account, err := api.Authorize(c)

		if err != nil {
			ResponseError(c, http.StatusUnauthorized, "You need to be logged in")
			return
		}
		var ad usecase.CreateAdInput
		err = c.BindJSON(&ad)
		ad.UserId = string(account)
		if err != nil {
			a.logger.Error().Msgf("User CreateAdInput invalid: %v", ad)
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

// UpdateAdHandler return the handler responsible for updating an ad.
func (a Handler) UpdateAdHandler(cmd usecase.UpdateAdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("UpdateAdHandler: param ID missing.")
			ResponseError(c, http.StatusBadRequest, MissingIDParam)
			return
		}
		/* AUTHORIZE */
		_, err := api.Authorize(c)

		if err != nil {
			ResponseError(c, http.StatusUnauthorized, "You need to be logged in")
			return
		}
		/* END AUTHORIZE */

		var ad usecase.UpdateAdInput
		err = c.BindJSON(&ad)
		ad.ID = domain.AdID(id)

		if err != nil {
			a.logger.Error().Msgf("User UpdateAdInput invalid: %v", ad)
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

// DeleteAdHandler return the handler responsible for deleting an ad.
func (a Handler) DeleteAdHandler(cmd usecase.DeleteAdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("DeleteAdHandler: param ID missing.")
			ResponseError(c, http.StatusBadRequest, MissingIDParam)
			return
		}

		/* AUTHORIZE */
		_, err := api.Authorize(c)

		if err != nil {
			ResponseError(c, http.StatusUnauthorized, "You need to be logged in")
			return
		}
		/* END AUTHORIZE */

		payload, err := cmd(c.Request.Context(), usecase.DeleteAdInput{ID: domain.AdID(id)})
		if err != nil {
			a.logger.Error().Msgf("Error in POST delete ad: %v", err)
			ResponseError(c, http.StatusInternalServerError, InternalServerError)
			return
		}
		ResponseSuccess(c, http.StatusAccepted, payload)
	}
}

// DeleteAdHandler return the handler responsible for searcgubg an ad.
func (a Handler) SearchAdHandler(cmd usecase.SearchAdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		content := c.Query("content")
		if content == "" {
			a.logger.Error().Msg("SearchAdHandler: param content missing.")
			ResponseError(c, http.StatusBadRequest, MissingQueryContent)
			return
		}

		payload, err := cmd(c.Request.Context(), usecase.SearchAdInput{Content: content})
		if err != nil {
			a.logger.Error().Msgf("Error in GET search ad: %v", err)
			ResponseError(c, http.StatusInternalServerError, InternalServerError)
			return
		}
		ResponseSuccess(c, http.StatusOK, payload)
	}
}
