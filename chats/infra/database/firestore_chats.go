package database

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	firebaseDB "firebase.google.com/go/db"
	"github.com/There-is-Go-alternative/GoMicroServices/chats/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/chats/internal/xerrors"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type FirebaseRealTimeDB struct {
	App  *firebase.App
	Conf *FirebaseConfig
	DB   *firebaseDB.Client
}

var ChatsDefaultConf = &FirebaseConfig{
	CollectionName:    "chats",
	ServiceAdsKeyPath: "infra/database/FirebaseCredentials.json",
	BaseConfig: &firebase.Config{
		DatabaseURL: "https://gomicroservicedatabase-default-rtdb.firebaseio.com/",
	},
}

type FirebaseConfig struct {
	CollectionName    string
	ServiceAdsKeyPath string
	BaseConfig        *firebase.Config
}

func NewFirebaseRealTimeDB(ctx context.Context, conf *FirebaseConfig) (*FirebaseRealTimeDB, error) {
	opt := option.WithCredentialsFile(conf.ServiceAdsKeyPath)
	opt2 := option.WithEndpoint(conf.BaseConfig.DatabaseURL)

	app, err := firebase.NewApp(ctx, conf.BaseConfig, opt, opt2)
	if err != nil {
		return nil, err
	}
	db, err := app.DatabaseWithURL(ctx, conf.BaseConfig.DatabaseURL)
	if err != nil {
		return nil, err
	}
	return &FirebaseRealTimeDB{
		App:  app,
		Conf: conf,
		DB:   db,
	}, nil
}

// Create add a chat to the Firestore realtime database
func (m *FirebaseRealTimeDB) Create(ctx context.Context, chat domain.Chat) error {

	err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, chat.ID.String())).Set(ctx, chat)
	if err != nil {
		return err
	}
	return nil
}

// ByID Retrieve the info that match "id".
// Strict: As ID is the key of the map, return an error if not found
func (m *FirebaseRealTimeDB) ByID(ctx context.Context, ID domain.ChatID) (*domain.Chat, error) {
	var chat domain.Chat
	if err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, ID)).Get(ctx, &chat); err != nil {
		return nil, err
	}

	if chat.ID.String() == "" {
		return nil, errors.Wrapf(
			xerrors.ErrorWithCode{Code: xerrors.ResourceNotFound, Err: xerrors.ChatNotFound}, "ID {%v}", ID,
		)
	}
	return &chat, nil
}

// All return all domain.Ad in the Firestore realtime database
func (m *FirebaseRealTimeDB) All(ctx context.Context) ([]*domain.Chat, error) {
	var chat map[string]*domain.Chat
	if err := m.DB.NewRef(m.Conf.CollectionName).OrderByChild("id").Get(ctx, &chat); err != nil {
		return nil, err
	}
	lst := make([]*domain.Chat, 0, len(chat))
	for _, a := range chat {
		lst = append(lst, a)
	}
	return lst, nil
}

// Remove a domain.Ad from the Firestore realtime database
// func (m *FirebaseRealTimeDB) Remove(ctx context.Context, ads ...*domain.Ad) error {
// 	if len(ads) <= 0 {
// 		return nil
// 	}

// 	errs := xerrors.ErrList{}
// 	for _, ad := range ads {
// 		err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, ad.ID.String())).Delete(ctx)
// 		if err != nil {
// 			errs.Add(err)
// 		}
// 	}

// 	if !errs.Nil() {
// 		return errs
// 	}
// 	return nil
// }
