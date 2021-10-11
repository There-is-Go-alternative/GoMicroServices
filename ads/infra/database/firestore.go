package database

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	firebaseDB "firebase.google.com/go/db"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal/xerrors"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type FirebaseRealTimeDB struct {
	App  *firebase.App
	Conf *FirebaseConfig
	DB   *firebaseDB.Client
	Storage *storage.Client
	Client *firestore.Client
}

var DefaultConf = &FirebaseConfig{
	CollectionName:    "ads",
	ServiceAdsKeyPath: "FirebaseCredentials.json",
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
	firestore_client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	storage, err := storage.NewClient(ctx, opt)
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
		Storage: storage,
		Client: firestore_client,
	}, nil
}

type ImageStructure struct {
	ImageName string `json:"imageName"`
	URL       string `json:"url"`
}

func DownloadImage() io.Reader {
	url := "http://i.imgur.com/m1UIjW1.jpg"

    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }

	return response.Body
}

func UploadImage(m *FirebaseRealTimeDB, ctx context.Context) error {
	wc := m.Storage.Bucket("gomicroservicedatabase.appspot.com").Object("Aucun").NewWriter(ctx)
	body := DownloadImage()

	id := uuid.New() 
	wc.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()} 
	_, err := io.Copy(wc, body)
	if err != nil {
		fmt.Print(err)
		fmt.Printf("ERROR")
		return err
	}
	if err := wc.Close(); err != nil {
		fmt.Printf("ERROR3")
		return err
	}
	imageStructure := ImageStructure{
		ImageName: "Aucune",
		URL:       "https://storage.cloud.google.com/gomicroservicedatabase.appspot.com/Aucune",
	}
	_, _, err = m.Client.Collection("image").Add(ctx, imageStructure)

	if err != nil {
		fmt.Printf("ERROR2")
		return err
	}
	return nil
} 

// Create add list of domain.Ad to the Firestore realtime database
func (m *FirebaseRealTimeDB) Create(ctx context.Context, ads ...*domain.Ad) error {
	UploadImage(m, ctx)
	if len(ads) == 0 {
		return nil
	}

	errs := xerrors.ErrList{}
	for _, ad := range ads {
		err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, ad.ID.String())).Set(ctx, ad)
		if err != nil {
			errs.Add(err)
		}
	}

	if !errs.Nil() {
		return errs
	}
	return nil
}

// Update a list of domain.Ad to the Firestore realtime database
func (m *FirebaseRealTimeDB) Update(ctx context.Context, ads ...*domain.Ad) error {
	adTransaction := func(ad *domain.Ad) func(transaction firebaseDB.TransactionNode) (interface{}, error) {
		return func(transaction firebaseDB.TransactionNode) (interface{}, error) {
			var new_ad domain.Ad

			if err := transaction.Unmarshal(&new_ad); err != nil {
				return nil, err
			}

			new_ad = *ad
			return new_ad, nil
		}
	}
	errs := xerrors.ErrList{}
	for _, ad := range ads {
		err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, ad.ID.String())).Transaction(ctx, adTransaction(ad))
		if err != nil {
			errs.Add(err)
		}
	}

	if !errs.Nil() {
		return errs
	}
	return nil
}

// ByID Retrieve the info that match "id".
// Strict: As ID is the key of the map, return an error if not found
func (m *FirebaseRealTimeDB) ByID(ctx context.Context, ID domain.AdID) (*domain.Ad, error) {
	var ad domain.Ad
	if err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, ID)).Get(ctx, &ad); err != nil {
		return nil, err
	}

	if ad.ID == "" {
		return nil, errors.Wrapf(
			xerrors.ErrorWithCode{Code: xerrors.ResourceNotFound, Err: xerrors.AdNotFound}, "ID {%v}", ID,
		)
	}
	return &ad, nil
}

// All return all domain.Ad in the Firestore realtime database
func (m *FirebaseRealTimeDB) All(ctx context.Context) ([]*domain.Ad, error) {
	var ads map[string]*domain.Ad
	if err := m.DB.NewRef(m.Conf.CollectionName).OrderByChild("id").Get(ctx, &ads); err != nil {
		return nil, err
	}
	lst := make([]*domain.Ad, 0, len(ads))
	for _, a := range ads {
		lst = append(lst, a)
	}
	return lst, nil
}

// Remove a domain.Ad from the Firestore realtime database
func (m *FirebaseRealTimeDB) Remove(ctx context.Context, ads ...*domain.Ad) error {
	if len(ads) <= 0 {
		return nil
	}

	errs := xerrors.ErrList{}
	for _, ad := range ads {
		err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, ad.ID.String())).Delete(ctx)
		if err != nil {
			errs.Add(err)
		}
	}

	if !errs.Nil() {
		return errs
	}
	return nil
}
