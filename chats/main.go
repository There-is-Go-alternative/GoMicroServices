package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	account "github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	chats "github.com/There-is-Go-alternative/GoMicroServices/chats/domain"
	database "github.com/There-is-Go-alternative/GoMicroServices/chats/infra/database"
	log "github.com/sirupsen/logrus"
)

var (
	confPath        = flag.String("conf", os.Getenv("CONF_PATH"), "path to the json config file.")
	shutdownTimeOut = flag.Int("timeout", 2, "Time out for graceful shutdown, in seconds.")
)

func Firebase() {
	// Setup context to be notified when the program receive a signal
	log.WithFields(log.Fields{
		"stage": "setup",
	}).Info("Setting up context ...")
	signalCtx, _ := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	ctx, ctxCancel := context.WithCancel(signalCtx)
	_ = ctxCancel

	// Initialising Chats Database
	log.WithFields(log.Fields{
		"stage": "setup",
	}).Info("Setting up Chats Database ...")
	ChatsStorage, err := database.NewFirebaseRealTimeDB(ctx, database.ChatsDefaultConf)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(42)
	}

	// Initialising Messages Database
	log.WithFields(log.Fields{
		"stage": "setup",
	}).Info("Setting up Messages Database ...")
	MessagesStorage, err := database.NewFirebaseRealTimeDB(ctx, database.MessagesDefaultConf)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(42)
	}

	userA := account.Account{ID: "UserA"}
	userB := account.Account{ID: "UserB"}
	userC := account.Account{ID: "UserC"}
	chatA := chats.Chat{ID: chats.ChatID("ChatA"), UsersIDs: []string{userA.ID.String(), userB.ID.String()}}
	chatB := chats.Chat{ID: chats.ChatID("ChatB"), UsersIDs: []string{userC.ID.String(), userB.ID.String()}}
	chatC := chats.Chat{ID: chats.ChatID("ChatC"), UsersIDs: []string{userA.ID.String(), userC.ID.String()}}
	chatD := chats.Chat{ID: chats.ChatID("ChatD"), UsersIDs: []string{userA.ID.String(), userB.ID.String(), userC.ID.String()}}
	messageA := chats.Message{ID: chats.MessageID("Message A to B"), ChatID: chatA.ID.String(), Content: "hello :)", SenderID: userA.ID.String()}
	messageB := chats.Message{ID: chats.MessageID("Message C to B"), ChatID: chatB.ID.String(), Content: "hello :)", SenderID: userC.ID.String()}
	messageC := chats.Message{ID: chats.MessageID("Message A to C"), ChatID: chatC.ID.String(), Content: "hello :)", SenderID: userA.ID.String()}

	//Create Chats
	err = ChatsStorage.CreateChat(ctx, chatA)
	if err != nil {
		fmt.Printf("Create ChatA error:\n%+v\nEnd\n", err)
	}
	err = ChatsStorage.CreateChat(ctx, chatB)
	if err != nil {
		fmt.Printf("Create ChatB error:\n%+v\nEnd\n", err)
	}
	err = ChatsStorage.CreateChat(ctx, chatC)
	if err != nil {
		fmt.Printf("Create ChatC error:\n%+v\nEnd\n", err)
	}
	err = ChatsStorage.CreateChat(ctx, chatD)
	if err != nil {
		fmt.Printf("Create ChatD error:\n%+v\nEnd\n", err)
	}

	// Get all chats of userA
	chatsOfUser, err := ChatsStorage.GetAllChatsOfUser(ctx, userA.ID.String())
	if err != nil {
		fmt.Printf("Chats of user error:\n%+v\nEnd\n", err)
	} else {
		fmt.Printf("Chats of user:\n%+v\nEnd\n", chatsOfUser)
	}

	//Send messages
	err = MessagesStorage.CreateMessage(ctx, messageA)
	if err != nil {
		fmt.Printf("Create messageA error:\n%+v\nEnd\n", err)
	}
	err = MessagesStorage.CreateMessage(ctx, messageB)
	if err != nil {
		fmt.Printf("Create messageB error:\n%+v\nEnd\n", err)
	}
	err = MessagesStorage.CreateMessage(ctx, messageC)
	if err != nil {
		fmt.Printf("Create messageC error:\n%+v\nEnd\n", err)
	}

	//Get all messages of chatA
	// messagesOfChatA, err := ChatsStorage.AllMessagesOfChat()

	// // Setup an error channel
	// errChan := make(chan error)

	// // Waiting for a channel to receive something
	// select {
	// case <-ctx.Done():
	// 	log.WithFields(log.Fields{
	// 		"stage": "runner",
	// 	}).Info("Context Canceled. Shutdown ...")
	// 	time.Sleep(time.Second * time.Duration(*shutdownTimeOut))
	// 	return
	// case err := <-errChan:
	// 	log.WithFields(log.Fields{
	// 		"stage": "runner",
	// 	}).Errorf("An Error happend in a service: %s", err)
	// 	// Cancel context to shut down blocking services.
	// 	ctxCancel()
	// 	time.Sleep(time.Second * time.Duration(*shutdownTimeOut))
	// 	os.Exit(1)
	// }
}

func main() {
	Firebase()
}
