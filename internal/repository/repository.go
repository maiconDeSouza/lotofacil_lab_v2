package repository

import (
	"context"
	"errors"
	"time"

	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/errapp"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryInterface interface {
	CreateContest(contest *model.Concurso) error
	CreateGame(game *model.Jogo) error
	FindGame(game *model.Jogo, strategy string, contest int) error
}

type Repository struct {
	collectionContest *mongo.Collection
	collectionGame    *mongo.Collection
}

func NewRepository(collectionContest *mongo.Collection, collectionGame *mongo.Collection) *Repository {
	return &Repository{collectionContest: collectionContest, collectionGame: collectionGame}
}

func (r *Repository) CreateContest(contest *model.Concurso) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	_, err := r.collectionContest.InsertOne(ctx, contest)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errapp.ErrContestAlreadyExists
		}
		return err
	}

	return nil
}

func (r *Repository) CreateGame(game *model.Jogo) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	_, err := r.collectionGame.InsertOne(ctx, game)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindGame(game *model.Jogo, strategy string, contest int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	filter := bson.M{
		"numero":         contest,
		"tipoEstrategia": strategy,
	}

	err := r.collectionGame.FindOne(ctx, filter).Decode(game)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errapp.ErrNoDocuments
		}
		return err
	}
	return nil
}
