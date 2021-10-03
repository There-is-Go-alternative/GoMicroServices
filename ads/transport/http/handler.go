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

// GetAdsHandler return the handler responsible for fetching all ads
func (a Handler) GetAdsHandler(cmd usecase.GetAllAdsCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := cmd(c.Request.Context())

		if err != nil {
			a.logger.Error().Msg("Error in GET /ads")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}

// GetAdsByIDHandler return the handler responsible for fetching a specific ad.
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
			a.logger.Error().Msg("Error in GET by /ads/:id")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}

// CreateAdHandler return the handler responsible for creating an ad.
func (a Handler) CreateAdHandler(cmd usecase.CreateAdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ad usecase.CreateAdInput
		err := c.BindJSON(&ad)
		if err != nil {
			a.logger.Error().Msgf("User CreateAdInput invalid: %v", ad)
			// TODO: better error
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		payload, err := cmd(c.Request.Context(), ad)
		if err != nil {
			a.logger.Error().Msgf("Error in POST create ad: %v", err)
			// TODO: better error
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusCreated, payload)
	}
}

// UpdateAdHandler return the handler responsible for updating an ad.
func (a Handler) UpdateAdHandler(cmd usecase.UpdateAdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("UpdateAdHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusInternalServerError, xerrors.MissingParam)
			return
		}
		var ad usecase.UpdateAdInput
		err := c.BindJSON(&ad)
		ad.ID = domain.AdID(id)

		if err != nil {
			a.logger.Error().Msgf("User UpdateAdInput invalid: %v", ad)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		payload, err := cmd(c.Request.Context(), ad)
		if err != nil {
			a.logger.Error().Msgf("Error in PATCH update ad: %v", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusCreated, payload)
	}
}

// DeleteAdHandler return the handler responsible for deleting an ad.
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
			a.logger.Error().Msgf("Error in POST delete ad: %v", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}

// DeleteAdHandler return the handler responsible for searcgubg an ad.
func (a Handler) SearchAdHandler(cmd usecase.SearchAdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		content := c.Query("content")
		if content == "" {
			a.logger.Error().Msg("SearchAdHandler: param content missing.")
			_ = c.AbortWithError(http.StatusInternalServerError, xerrors.MissingParam)
			return
		}

		payload, err := cmd(c.Request.Context(), usecase.SearchAdInput{Content: content})
		if err != nil {
			a.logger.Error().Msgf("Error in GET search ad: %v", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}