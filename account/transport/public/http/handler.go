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
	"strings"
)

const AccountIDKey = "accountID"

type Handler struct {
	APIKey      string
	AuthService usecase.AuthService
	logger      zerolog.Logger
}

func NewAccountHandler(APIKey string, authService usecase.AuthService) *Handler {
	return &Handler{
		APIKey:      APIKey,
		AuthService: authService,
		logger:      log.With().Str("service", "Http Handler").Logger(),
	}
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
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    payload,
		})
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
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    payload,
		})
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
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data":    payload,
		})
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
		if token := c.GetString(AccountIDKey); token != id && token != a.APIKey {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		var account usecase.PatchAccountInput
		err := c.BindJSON(&account)
		if err != nil {
			a.logger.Error().Msgf("PatchAccountInput invalid: %v", account)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		account.ID = domain.AccountID(id)
		payload, err := cmd(c.Request.Context(), account)
		if err != nil {
			a.logger.Error().Msgf("Error in PATCH account: %v", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data":    payload,
		})
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
		if token := c.GetString(AccountIDKey); token != id && token != a.APIKey {
			c.Status(http.StatusUnauthorized)
			c.Abort()
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
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data":    payload,
		})
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
		if token := c.GetString(AccountIDKey); token != id && token != a.APIKey {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		payload, err := cmd(c.Request.Context(), usecase.DeleteAccountInput{ID: domain.AccountID(id)})
		if err != nil {
			a.logger.Error().Msgf("Error in POST delete account: %v", err)
			// TODO: better error
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    payload,
		})
	}
}

// IsAdminHandler.
func (a Handler) IsAdminHandler(cmd usecase.IsAdminCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id usecase.IsAdminInput
		err := c.BindJSON(&id)

		payload, err := cmd(c.Request.Context(), id)
		if err != nil {
			a.logger.Error().Msgf("Error in Is Admin : %v", err)
			// TODO: better error
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    payload,
		})
	}
}

func (a Handler) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("auth token missing"))
			return
		}
		if token == a.APIKey {
			c.Set(AccountIDKey, a.APIKey)
			return
		}
		tokenSplitted := strings.Split(token, "Bearer ")
		if len(tokenSplitted) != 2 {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		accountID, err := a.AuthService.Authorize(tokenSplitted[1])
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("auth token is not valid"))
			return
		}
		c.Set(AccountIDKey, accountID)
	}
}
