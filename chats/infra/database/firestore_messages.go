package database

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/There-is-Go-alternative/GoMicroServices/chats/domain"
)

var MessagesDefaultConf = &FirebaseConfig{
	CollectionName:    "messages",
	ServiceAdsKeyPath: "infra/database/FirebaseCredentials.json",
	BaseConfig: &firebase.Config{
		DatabaseURL: "https://gomicroservicedatabase-default-rtdb.firebaseio.com/",
	},
}

// Create adds a message to the Firestore realtime database
func (m *FirebaseRealTimeDB) CreateMessage(ctx context.Context, message domain.Message) error {

	err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, message.ID.String())).Set(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

// All return all domain.Message in the Firestore realtime database
func (m *FirebaseRealTimeDB) AllMessagesOfChat(ctx context.Context) ([]*domain.Message, error) {
	var message map[string]*domain.Message
	if err := m.DB.NewRef(m.Conf.CollectionName).Get(ctx, &message); err != nil {
		return nil, err
	}
	lst := make([]*domain.Message, 0, len(message))
	for _, a := range message {
		lst = append(lst, a)
	}
	return lst, nil
}
