package transport

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/There-is-Go-alternative/GoMicroServices/transactions/usecase"
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
	uc     *usecase.UseCase
}

func NewTransactionHandler(uc *usecase.UseCase) *Handler {
	return &Handler{logger: log.With().Str("service", "Http Handler").Logger(), uc: uc}
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

func (h Handler) getUserId(g *gin.Context) (*string, error) {
	token := g.GetHeader("Authorization")

	if token == "" {
		return nil, fmt.Errorf("no token provided")
	}
	hasPrefix := strings.HasPrefix(token, "Bearer ")

	if !hasPrefix {
		return nil, fmt.Errorf("wrong token format")
	}
	token = token[7:]
	userId, err := h.uc.Auth.GetUser(g, token)

	if err != nil {
		return nil, err
	}
	return &userId.Id, nil
}

func (h Handler) isAdmin(g *gin.Context) bool {
	userId, err := h.getUserId(g)

	if err != nil {
		return false
	}

	isAdmin, err := h.uc.Auth.IsAdmin(*userId)

	if err != nil {
		return false
	}

	return isAdmin
}

func (h Handler) hasAccess(g *gin.Context, id string) bool {
	userId, err := h.getUserId(g)

	if err != nil {
		return false
	}

	isAdmin, err := h.uc.Auth.IsAdmin(*userId)

	if err != nil {
		return false
	}

	return isAdmin || *userId == id
}

func (h Handler) Register(cmd usecase.RegiterCmd) gin.HandlerFunc {
	return func(g *gin.Context) {
		id := g.Param("id")
		userId, err := h.getUserId(g)

		if err != nil {
			h.logger.Error().Msg(fmt.Sprintf("Error in POST /register: %s", err))
			ResponseError(g, http.StatusInternalServerError, err)
			return
		}

		err = cmd(g, id, *userId)

		if err != nil {
			h.logger.Error().Msg(fmt.Sprintf("Error in POST /register: %s", err))
			ResponseError(g, http.StatusInternalServerError, err)
			return
		}

		ResponseSuccess(g, http.StatusOK, nil)
	}
}

func (h Handler) GetAll(cmd usecase.GetAllCmd) gin.HandlerFunc {
	return func(g *gin.Context) {
		if !h.isAdmin(g) {
			ResponseError(g, http.StatusNetworkAuthenticationRequired, fmt.Errorf("you are not admin"))
			return
		}
		payload, err := cmd(g.Request.Context())

		if err != nil {
			h.logger.Error().Msg(fmt.Sprintf("Error in GET /: %s", err))
			ResponseError(g, http.StatusInternalServerError, err)
			return
		}
		ResponseSuccess(g, http.StatusOK, payload)
	}
}

func (h Handler) GetById(cmd usecase.GetByIdCmd) gin.HandlerFunc {
	return func(g *gin.Context) {
		id := g.Param("id")

		if !h.isAdmin(g) {
			ResponseError(g, http.StatusInternalServerError, fmt.Errorf("you are not admin"))
			return
		}

		payload, err := cmd(g.Request.Context(), id)

		if err != nil {
			h.logger.Error().Msg(fmt.Sprintf("Error in GET /: %s", err))
			ResponseError(g, http.StatusInternalServerError, err)
			return
		}
		ResponseSuccess(g, http.StatusOK, payload)
	}
}

func (h Handler) GetByDateRange(cmd usecase.GetByDateRangeCmd) gin.HandlerFunc {
	return func(g *gin.Context) {
		if !h.isAdmin(g) {
			ResponseError(g, http.StatusInternalServerError, fmt.Errorf("you are not admin"))
			return
		}

		var input usecase.DateRangeInput
		err := g.BindJSON(&input)

		if err != nil {
			h.logger.Error().Msgf("DateRange invalid: %v", input)
			ResponseError(g, http.StatusBadRequest, err)
			return
		}

		payload, err := cmd(g.Request.Context(), input.Start, input.End)

		if err != nil {
			h.logger.Error().Msg(fmt.Sprintf("Error in GET /: %s", err))
			ResponseError(g, http.StatusInternalServerError, err)
			return
		}

		ResponseSuccess(g, http.StatusOK, payload)
	}
}

func (h Handler) GetAllBuyedByUser(cmd usecase.GetByBuyerIdCmd) gin.HandlerFunc {
	return func(g *gin.Context) {
		id := g.Param("user_id")
		payload, err := cmd(g.Request.Context(), id)

		if !h.hasAccess(g, id) {
			ResponseError(g, http.StatusNetworkAuthenticationRequired, fmt.Errorf("you don't have access to this ressource"))
			return
		}

		if err != nil {
			h.logger.Error().Msg(fmt.Sprintf("Error in GET /: %s", err))
			ResponseError(g, http.StatusInternalServerError, err)
			return
		}
		ResponseSuccess(g, http.StatusOK, payload)
	}
}

func (h Handler) GetBuyedByDateRangeUser(cmd usecase.GetByBuyerIdDateRangeCmd) gin.HandlerFunc {
	return func(g *gin.Context) {
		id := g.Param("user_id")

		if !h.hasAccess(g, id) {
			ResponseError(g, http.StatusNetworkAuthenticationRequired, fmt.Errorf("you don't have access to this ressource"))
			return
		}

		var input usecase.DateRangeInput
		err := g.BindJSON(&input)

		if err != nil {
			h.logger.Error().Msgf("DateRange invalid: %v", input)
			ResponseError(g, http.StatusBadRequest, err)
			return
		}

		payload, err := cmd(g.Request.Context(), id, input.Start, input.End)

		if err != nil {
			h.logger.Error().Msg(fmt.Sprintf("Error in GET /: %s", err))
			ResponseError(g, http.StatusInternalServerError, err)
			return
		}

		ResponseSuccess(g, http.StatusOK, payload)
	}
}

func (h Handler) GetAllSelledByUser(cmd usecase.GetBySellerIdCmd) gin.HandlerFunc {
	return func(g *gin.Context) {
		id := g.Param("user_id")

		if !h.hasAccess(g, id) {
			ResponseError(g, http.StatusNetworkAuthenticationRequired, fmt.Errorf("you don't have access to this ressource"))
			return
		}

		payload, err := cmd(g.Request.Context(), id)

		if err != nil {
			h.logger.Error().Msg(fmt.Sprintf("Error in GET /: %s", err))
			ResponseError(g, http.StatusInternalServerError, err)
			return
		}
		ResponseSuccess(g, http.StatusOK, payload)
	}
}

func (h Handler) GetSelledByDateRangeUser(cmd usecase.GetBySellerIdDateRangeCmd) gin.HandlerFunc {
	return func(g *gin.Context) {
		id := g.Param("user_id")

		if !h.hasAccess(g, id) {
			ResponseError(g, http.StatusNetworkAuthenticationRequired, fmt.Errorf("you don't have access to this ressource"))
			return
		}

		var input usecase.DateRangeInput
		err := g.BindJSON(&input)

		if err != nil {
			h.logger.Error().Msgf("DateRange invalid: %v", input)
			ResponseError(g, http.StatusBadRequest, err)
			return
		}

		payload, err := cmd(g.Request.Context(), id, input.Start, input.End)

		if err != nil {
			h.logger.Error().Msg(fmt.Sprintf("Error in GET /: %s", err))
			ResponseError(g, http.StatusInternalServerError, err)
			return
		}

		ResponseSuccess(g, http.StatusOK, payload)
	}
}
