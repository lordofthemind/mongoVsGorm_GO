package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/mongoVsGorm_GO/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAuthorRepository struct {
	collection *mongo.Collection
}

func NewMongoAuthorRepository(db *mongo.Database) *MongoAuthorRepository {
	return &MongoAuthorRepository{collection: db.Collection("authors")}
}

func (repo *MongoAuthorRepository) CreateAuthor(ctx context.Context, name string, bio string, email string, dateOfBirth *time.Time) (uuid.UUID, error) {
	id := uuid.New()
	author := types.Author{
		ID:          id,
		Name:        name,
		Bio:         bio,
		Email:       email,
		DateOfBirth: dateOfBirth,
	}
	_, err := repo.collection.InsertOne(ctx, author)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (repo *MongoAuthorRepository) GetAuthor(ctx context.Context, id uuid.UUID) (types.Author, error) {
	var author types.Author
	err := repo.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&author)
	if err != nil {
		return types.Author{}, err
	}
	return author, nil
}

func (repo *MongoAuthorRepository) ListAuthors(ctx context.Context) ([]types.Author, error) {
	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var authors []types.Author
	if err := cursor.All(ctx, &authors); err != nil {
		return nil, err
	}
	return authors, nil
}

func (repo *MongoAuthorRepository) DeleteAuthor(ctx context.Context, id uuid.UUID) error {
	_, err := repo.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (repo *MongoAuthorRepository) UpdateAuthor(ctx context.Context, id uuid.UUID, name string, bio string, email string, dateOfBirth *time.Time) error {
	update := bson.M{
		"$set": bson.M{
			"name":          name,
			"bio":           bio,
			"email":         email,
			"date_of_birth": dateOfBirth,
		},
	}
	_, err := repo.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (repo *MongoAuthorRepository) GetAuthorsByBirthdateRange(ctx context.Context, startDate, endDate time.Time) ([]types.Author, error) {
	cursor, err := repo.collection.Find(ctx, bson.M{
		"date_of_birth": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var authors []types.Author
	if err := cursor.All(ctx, &authors); err != nil {
		return nil, err
	}
	return authors, nil
}
