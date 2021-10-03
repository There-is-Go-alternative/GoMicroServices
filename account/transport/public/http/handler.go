package http

import (
	"errors"
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

// GetAllAccountsHandler return the handler responsible for fetching all users account
func (a Handler) GetAllAccountsHandler(cmd usecase.GetAllAccountsCmd) gin.HandlerFunc {
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
func (a Handler) GetAccountsByIDHandler(cmd usecase.GetAccountByIDCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("GetAccountsByIDHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusInternalServerError, xerrors.MissingParam)
			return
		}
		payload, err := cmd(c.Request.Context(), domain.AccountID(id))

		if err != nil {
			a.logger.Error().Err(err).Msg("Error in GET by /account/:id")
			// Todo : APIError with switch case on error Code
			if errors.Is(err, xerrors.AccountNotFound) {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
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
			a.logger.Error().Msgf("CreateAccountInput invalid: %v", account)
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

// PatchAccountHandler return the handler responsible for updating a user account.
func (a Handler) PatchAccountHandler(cmd usecase.PatchAccountCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("PatchAccountHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusInternalServerError, xerrors.MissingParam)
			return
		}
		var account usecase.PatchAccountInput
		err := c.BindJSON(&account)
		if err != nil {
			a.logger.Error().Msgf("PatchAccountInput invalid: %v", account)
			// TODO: better error
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		account.ID = domain.AccountID(id)
		payload, err := cmd(c.Request.Context(), account)
		if err != nil {
			a.logger.Error().Msgf("Error in PATCH account: %v", err)
			// TODO: better error
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusCreated, payload)
	}
}

// PutAccountHandler return the handler responsible for replacing a user account.
func (a Handler) PutAccountHandler(cmd usecase.UpdateAccountCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			a.logger.Error().Msg("PutAccountHandler: param ID missing.")
			_ = c.AbortWithError(http.StatusInternalServerError, xerrors.MissingParam)
			return
		}
		var account usecase.UpdateAccountInput
		err := c.BindJSON(&account)
		if err != nil {
			a.logger.Error().Msgf("PutAccountInput invalid: %v", account)
			// TODO: better error
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		account.ID = domain.AccountID(id)
		payload, err := cmd(c.Request.Context(), account)
		if err != nil {
			a.logger.Error().Msgf("Error in PUT account: %v", err)
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
