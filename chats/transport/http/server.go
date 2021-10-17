package http

import (
	"context"
	"fmt"
	netHTTP "net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/chats/internal/config"
	"github.com/There-is-Go-alternative/GoMicroServices/chats/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Engine *netHTTP.Server
	logger zerolog.Logger
}

type useCase interface {
	CreateChat() usecase.CreateChatCmd
	GetChatById() usecase.GetChatByIdCmd
	GetAllChatsOfUser() usecase.GetAllChatsOfUserCmd
	CreateMessage() usecase.CreateMessageCmd
	GetMessagesByChatID() usecase.GetMessagesByChatIDCmd
}

// TODO: change database by future Database interface
func NewChatsHttpServer(uc useCase, conf *config.Config) *Server {
	router := gin.Default()
	router.Use(cors.Default())

	chatHandler := NewChatHandler()
	// Grouping Chat routes with url specified in config (I.E: 'chat')
	chat := router.Group(fmt.Sprintf("/%s", conf.ChatsEndpoint))
	{
		chat.GET("/health", func(c *gin.Context) {
			c.Status(netHTTP.StatusOK)
		})
		chat.POST("/", chatHandler.CreateChatHandler(uc.CreateChat()))
		chat.GET("/:id", chatHandler.GetChatByIDHandler(uc.GetChatById()))
		chat.GET("/user/:user_id", chatHandler.GetAllChatsOfUserHandler(uc.GetAllChatsOfUser()))
	}
	return &Server{
		Engine: &netHTTP.Server{
			Addr:    fmt.Sprintf("%s:%s", conf.ChatsHost, conf.ChatsPort),
			Handler: router,
		},
		logger: log.With().Str("service", "Chats HTTP gin server").Logger(),
	}
}

func NewMessagesHttpServer(uc useCase, conf *config.Config) *Server {
	router := gin.Default()
	router.Use(cors.Default())

	messageHandler := NewMessageHandler()
	// Grouping Chat routes with url specified in config (I.E: 'chat')
	message := router.Group(fmt.Sprintf("/%s", conf.MessagesEndpoint))
	{
		message.GET("/health", func(c *gin.Context) {
			c.Status(netHTTP.StatusOK)
		})
		message.POST("/", messageHandler.CreateMessageHandler(uc.CreateMessage()))
		message.GET("/:id", messageHandler.GetMessagesByChatIDHandler(uc.GetMessagesByChatID()))
	}
	return &Server{
		Engine: &netHTTP.Server{
			Addr:    fmt.Sprintf("%s:%s", conf.MessageHost, conf.MessagePort),
			Handler: router,
		},
		logger: log.With().Str("service", "Messages HTTP gin server").Logger(),
	}
}

func (s Server) Run(ctx context.Context) (err error) {
	s.logger.Info().Msg("Running gin HTTP server ...")
	errc := make(chan error)
	go func() {
		errc <- s.Engine.ListenAndServe()
	}()
	select {
	case err = <-errc:
		return
	case <-ctx.Done():
		if err = s.Engine.Shutdown(ctx); err != nil && err != context.Canceled {
			s.logger.Error().Msgf("Error happened when server forced to shutdown: %v", err)
			return
		}
	}
	return
}
