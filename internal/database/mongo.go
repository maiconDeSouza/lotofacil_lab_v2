package database

import (
	"context"
	"fmt"
	"time"

	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/env"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/errapp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collections struct {
	CollectionContest *mongo.Collection
	CollectionGame    *mongo.Collection
}

func NewConnectMongo(env *env.Env) (*Collections, *errapp.ErrApp) {
	uri := fmt.Sprintf(
		"mongodb+srv://%s:%s@cluster0.xhbyyve.mongodb.net/",
		env.MONGO_ROOT_USER,
		env.MONGO_ROOT_PASSWORD,
	)
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		errapp := errapp.NewErrApp(err, "❌​ Erro ao ligar ao MongoDB.")
		return nil, errapp
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		errapp := errapp.NewErrApp(err, "❌​ Não foi possível alcançar o MongoDB.")
		return nil, errapp
	}

	db := client.Database(env.MONGO_DB_NAME)
	collectionContest := db.Collection(env.MONGO_COLLECTION_CONTESTS)
	collectionGame := db.Collection(env.MONGO_COLLECTION_GAMERS)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "numero", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	if _, err := collectionContest.Indexes().CreateOne(ctx, indexModel); err != nil {
		return nil, errapp.NewErrApp(err, "❌ Erro ao criar índice único na collection de concursos.")
	}

	fmt.Println("✅ Conectado ao banco de dados com sucesso!")

	collections := &Collections{
		CollectionContest: collectionContest,
		CollectionGame:    collectionGame,
	}

	return collections, nil
}
