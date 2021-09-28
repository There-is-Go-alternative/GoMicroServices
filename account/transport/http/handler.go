package http

import (
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/xerrors"
	"github.com/There-is-Go-alternative/GoMicroServices/account/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Handler struct {
	logger zerolog.Logger
}

func NewAccountHandler() *Handler {
	return &Handler{logger: log.With().Str("service", "Http Handler").Logger()}
}

// GetAccountsHandler return the handler responsible for fetching all users account
func (a Handler) GetAccountsHandler(cmd usecase.GetAllAccountsCmd) gin.HandlerFunc {
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

// GetAccountsByIDHandler return the handler responsible for fetching a specific account.
func (a Handler) GetAccountsByIDHandler(cmd usecase.GetAccountByIdCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("GetAccountsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusInternalServerError, xerrors.MissingParam)
			return
		}
		payload, err := cmd(c.Request.Context(), domain.AccountID(id))

		if err != nil {
			a.logger.Error().Msg("Error in GET by /account/:id")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}

// CreateAccountHandler return the handler responsible for creating a user account.
func (a Handler) CreateAccountHandler(cmd usecase.CreateAccountCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		var account usecase.CreateAccountInput
		err := c.BindJSON(&account)
		if err != nil {
			a.logger.Error().Msgf("User CreateAccountInput invalid: %v", account)
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

// DeleteAccountHandler return the handler responsible for deleting a user account.
func (a Handler) DeleteAccountHandler(cmd usecase.DeleteAccountCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("DeleteAccountHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusInternalServerError, xerrors.MissingParam)
			return
		}

		payload, err := cmd(c.Request.Context(), usecase.DeleteAccountInput{ID: domain.AccountID(id)})
		if err != nil {
			a.logger.Error().Msgf("Error in POST delete account: %v", err)
			// TODO: better error
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}
