// internal/repository/mongodb/user.go
package mongodb

import (
    "context"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "myapp/internal/model"
)

type UserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
    collection := db.Collection("users")
    
    // Create indexes
    indexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
    }

    ctx := context.Background()
    if _, err := collection.Indexes().CreateMany(ctx, indexes); err != nil {
        log.Printf("Failed to create indexes: %v", err)
    }

    return &UserRepository{
        collection: collection,
    }
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
    _, err := r.collection.InsertOne(ctx, user)
    return err
}

func (r *UserRepository) Find(ctx context.Context, id string) (*model.User, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var user model.User
    err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*model.User, error) {
    cursor, err := r.collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var users []*model.User
    if err = cursor.All(ctx, &users); err != nil {
        return nil, err
    }

    return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
    _, err := r.collection.UpdateOne(
        ctx,
        bson.M{"_id": user.ID},
        bson.D{{"$set", user}},
    )
    return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    _, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
    return err
}