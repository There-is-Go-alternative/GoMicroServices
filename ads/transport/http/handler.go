package http

import (
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal/xerrors"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	logger zerolog.Logger
}

func NewAdHandler() *Handler {
	return &Handler{logger: log.With().Str("service", "Http Handler").Logger()}
}

// GetAdsHandler return the handler responsible for fetching all users account
func (a Handler) GetAdsHandler(cmd usecase.GetAllAdsCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := cmd(c.Request.Context())

		if err != nil {
			a.logger.Error().Msg("Error in GET /accounts")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}

// GetAdsByIDHandler return the handler responsible for fetching a specific account.
func (a Handler) GetAdsByIDHandler(cmd usecase.GetAdByIdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("GetAdsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusInternalServerError, xerrors.MissingParam)
			return
		}
		payload, err := cmd(c.Request.Context(), domain.AdID(id))

		if err != nil {
			a.logger.Error().Msg("Error in GET by /account/:id")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}

// CreateAdHandler return the handler responsible for creating a user account.
func (a Handler) CreateAdHandler(cmd usecase.CreateAdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		var account usecase.CreateAdInput
		err := c.BindJSON(&account)
		if err != nil {
			a.logger.Error().Msgf("User CreateAdInput invalid: %v", account)
			// TODO: better error
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		payload, err := cmd(c.Request.Context(), account)
		if err != nil {
			a.logger.Error().Msgf("Error in POST create account: %v", err)
			// TODO: better error
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusCreated, payload)
	}
}

// DeleteAdHandler return the handler responsible for deleting a user account.
func (a Handler) DeleteAdHandler(cmd usecase.DeleteAdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("DeleteAdHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusInternalServerError, xerrors.MissingParam)
			return
		}

		payload, err := cmd(c.Request.Context(), usecase.DeleteAdInput{ID: domain.AdID(id)})
		if err != nil {
			a.logger.Error().Msgf("Error in POST delete account: %v", err)
			// TODO: better error
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}
