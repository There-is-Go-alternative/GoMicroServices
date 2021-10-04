package http

import (
	"fmt"
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/funds/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	logger zerolog.Logger
}

func NewFundsHandler() *Handler {
	return &Handler{logger: log.With().Str("service", "Http Handler").Logger()}
}

func (a Handler) CreateFundsHandler(cmd usecase.CreateFundsCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		a.logger.Debug().Msg("here lol")
		id := c.Param("id")

		a.logger.Debug().Msg(fmt.Sprintf("id: %s", id))

		if id == "" {
			a.logger.Error().Msg("GetFundsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{})
			return
		}

		payload, err := cmd(c.Request.Context(), usecase.CreateFundsInput{UserId: id, Balance: 0})

		a.logger.Debug().Msg(fmt.Sprintf("%v", payload))

		if err != nil {
			a.logger.Error().Msg(fmt.Sprintf("Error in GET /funds: %s", err))
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}

func (a Handler) GetAllFundsHandler(cmd usecase.AllCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := cmd(c.Request.Context())

		if err != nil {
			a.logger.Error().Msg(fmt.Sprintf("Error in GET /funds: %s", err))
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}

func (a Handler) GetFundsByIDHandler(cmd usecase.GetByIDCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusServiceUnavailable)
	}
}

func (a Handler) GetFundsByUserIDHandler(cmd usecase.GetByUserIDCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusServiceUnavailable)
	}
}

func (a Handler) DeleteFundsByIDHandler(cmd usecase.DeleteByIDCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusServiceUnavailable)
	}
}

func (a Handler) DeleteFundsByUserIDHandler(cmd usecase.DeleteByUserIDCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusServiceUnavailable)
	}
}

func (a Handler) IncreaseFundsHandler(cmd usecase.IncreaseCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusServiceUnavailable)
	}
}

func (a Handler) DecreaseFundsHandler(cmd usecase.DecreaseCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusServiceUnavailable)
	}
}
