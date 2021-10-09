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

func Memory() {
	userA := account.Account{ID: "UserA"}
	userB := account.Account{ID: "UserB"}
	userC := account.Account{ID: "UserC"}
	fmt.Printf("%+v\n", userA.ID)
	fmt.Printf("%+v\n", userB.ID)
	fmt.Printf("%+v\n\n", userC.ID)

	chatA := chats.Chat{ID: chats.ChatID("ChatA"), UsersIDs: []account.AccountID{userA.ID, userB.ID}}
	chatB := chats.Chat{ID: chats.ChatID("ChatB"), UsersIDs: []account.AccountID{userC.ID, userB.ID}}
	chatC := chats.Chat{ID: chats.ChatID("ChatC"), UsersIDs: []account.AccountID{userA.ID, userC.ID}}
	fmt.Printf("%+v\n", chatA)
	fmt.Printf("%+v\n", chatB)
	fmt.Printf("%+v\n\n", chatC)

	messageA := chats.Message{ID: chats.MessageID("Message A to B"), ChatID: chatA.ID, Content: "hello :)", SenderID: userA.ID.String()}
	messageB := chats.Message{ID: chats.MessageID("Message C to B"), ChatID: chatB.ID, Content: "hello :)", SenderID: userC.ID.String()}
	messageC := chats.Message{ID: chats.MessageID("Message A to C"), ChatID: chatC.ID, Content: "hello :)", SenderID: userA.ID.String()}
	fmt.Printf("%+v\n", messageA)
	fmt.Printf("%+v\n", messageB)
	fmt.Printf("%+v\n\n", messageC)

	memory := database.MockDatabase{}
	memory.CreateChat(chatA)
	memory.CreateChat(chatB)
	memory.CreateChat(chatC)
	memory.CreateMessage(messageA)
	memory.CreateMessage(messageB)
	memory.CreateMessage(messageC)
	fmt.Printf("%+v\n%+v\n\n", memory.Chats, memory.Messages)

	getChat, _ := memory.GetChat("ChatB")
	fmt.Printf("%+v\n\n", getChat)

	getMsg, err := memory.GetMessage("Message C to B")
	if err != nil {
		fmt.Printf("There was an error!")
	} else {
		fmt.Printf("%+v\n\n", getMsg)
	}

	memory.DeleteMessage("Message C to B")
	getMsg, err = memory.GetMessage("Message C to B")
	if err != nil {
		fmt.Printf("There was an error!")
	} else {
		fmt.Printf("%+v\n\n", getMsg)
	}

	getAllChats, err := memory.GetAllChatsOfUser(userA.ID)
	if err != nil {
		fmt.Printf("There was an error!")
	} else {
		fmt.Printf("%+v\n", getAllChats)
	}

	getAllMessages, err := memory.GetAllMessagesOfChat(chatB.ID)
	if err != nil {
		fmt.Printf("There was an error!")
	} else {
		fmt.Printf("%+v\n", getAllMessages)
	}
}

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
	}).Info("Setting up Ads Database ...")
	ChatsStorage, err := database.NewFirebaseRealTimeDB(ctx, database.ChatsDefaultConf)
	_ = ChatsStorage
	if err != nil {
		log.Fatal(err)
	}

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
