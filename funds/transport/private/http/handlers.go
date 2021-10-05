package http

import (
	"fmt"
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/funds/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/funds/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	logger zerolog.Logger
	auth   usecase.AuthService
}

func NewFundsHandler(auth *usecase.AuthService) *Handler {
	return &Handler{logger: log.With().Str("service", "Http Handler").Logger(), auth: *auth}
}

func (a Handler) ValidateToken(c *gin.Context) {
	authorization := c.GetHeader("Authorization")

	if authorization == "" {
		a.logger.Error().Msg("No token")
		_ = c.AbortWithError(http.StatusNetworkAuthenticationRequired, fmt.Errorf("no token provided in Authorization header"))
		return
	}

	err := a.auth.ValidateToken(authorization)

	if err != nil {
		a.logger.Error().Msg(fmt.Sprintf("Error in GET /funds: %s", err))
		_ = c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
		return
	}

	c.Next()
}

func (a Handler) CreateFundsHandler(cmd usecase.CreateFundsCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			a.logger.Error().Msg("GetFundsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{})
			return
		}

		payload, err := cmd(c.Request.Context(), usecase.CreateFundsInput{UserId: id, Balance: 0})

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
		id := c.Param("id")

		if id == "" {
			a.logger.Error().Msg("GetFundsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{})
			return
		}

		payload, err := cmd(c.Request.Context(), domain.FundsID(id))

		if err != nil {
			a.logger.Error().Msg(fmt.Sprintf("Error in GET /funds: %s", err))
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}

func (a Handler) GetFundsByUserIDHandler(cmd usecase.GetByUserIDCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			a.logger.Error().Msg("GetFundsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{})
			return
		}

		payload, err := cmd(c.Request.Context(), id)

		if err != nil {
			a.logger.Error().Msg(fmt.Sprintf("Error in GET /funds: %s", err))
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}

func (a Handler) DeleteFundsByIDHandler(cmd usecase.DeleteByIDCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			a.logger.Error().Msg("GetFundsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{})
			return
		}

		err := cmd(c.Request.Context(), domain.FundsID(id))

		if err != nil {
			a.logger.Error().Msg(fmt.Sprintf("Error in GET /funds: %s", err))
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.Status(http.StatusOK)
	}
}

func (a Handler) DeleteFundsByUserIDHandler(cmd usecase.DeleteByUserIDCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			a.logger.Error().Msg("GetFundsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{})
			return
		}

		err := cmd(c.Request.Context(), id)

		if err != nil {
			a.logger.Error().Msg(fmt.Sprintf("Error in GET /funds: %s", err))
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.Status(http.StatusOK)
	}
}

func (a Handler) IncreaseFundsHandler(cmd usecase.IncreaseCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			a.logger.Error().Msg("GetFundsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{})
			return
		}

		var input usecase.IncreaseDecreaseInput
		err := c.BindJSON(&input)

		if err != nil {
			a.logger.Error().Msgf("IncreaseDecreaseInput invalid: %v", input)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = cmd(c.Request.Context(), domain.FundsID(id), input)

		if err != nil {
			a.logger.Error().Msgf("Internal Server Error", input)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (a Handler) DecreaseFundsHandler(cmd usecase.DecreaseCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			a.logger.Error().Msg("GetFundsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{})
			return
		}

		var input usecase.IncreaseDecreaseInput
		err := c.BindJSON(&input)

		if err != nil {
			a.logger.Error().Msgf("IncreaseDecreaseInput invalid: %v", input)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = cmd(c.Request.Context(), domain.FundsID(id), input)

		if err != nil {
			a.logger.Error().Msgf("Internal Server Error", input)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (a Handler) SetFundsHandler(cmd usecase.SetCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			a.logger.Error().Msg("GetFundsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{})
			return
		}

		var input usecase.SetInput
		err := c.BindJSON(&input)

		if err != nil {
			a.logger.Error().Msgf("SetInput invalid: %v", input)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = cmd(c.Request.Context(), domain.FundsID(id), input)

		if err != nil {
			a.logger.Error().Msgf("Internal Server Error", input)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (a Handler) IncreaseFundsByUserHandler(cmd usecase.IncreaseByUserCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			a.logger.Error().Msg("GetFundsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{})
			return
		}

		var input usecase.IncreaseDecreaseInput
		err := c.BindJSON(&input)

		if err != nil {
			a.logger.Error().Msgf("IncreaseDecreaseInput invalid:", err)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = cmd(c.Request.Context(), id, input)

		if err != nil {
			a.logger.Error().Msgf("Internal Server Error", input)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (a Handler) DecreaseFundsByUserHandler(cmd usecase.DecreaseByUserCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			a.logger.Error().Msg("GetFundsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{})
			return
		}

		var input usecase.IncreaseDecreaseInput
		err := c.BindJSON(&input)

		if err != nil {
			a.logger.Error().Msgf("IncreaseDecreaseInput invalid: %v", input)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = cmd(c.Request.Context(), id, input)

		if err != nil {
			a.logger.Error().Msgf("Internal Server Error", input)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (a Handler) SetFundsByUserHandler(cmd usecase.SetByUserCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			a.logger.Error().Msg("GetFundsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{})
			return
		}

		var input usecase.SetInput
		err := c.BindJSON(&input)

		if err != nil {
			a.logger.Error().Msgf("SetInput invalid: %v", input)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = cmd(c.Request.Context(), id, input)

		if err != nil {
			a.logger.Error().Msgf("Internal Server Error", input)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
