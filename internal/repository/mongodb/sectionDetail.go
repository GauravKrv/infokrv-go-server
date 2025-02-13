// internal/repository/mongodb/sectionDetail.go
package mongodb

import (
    "context"
    "errors"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "myapp/internal/model"
)

type SectionDetailRepository struct {
    collection *mongo.Collection
}

func NewSectionDetailRepository(db *mongo.Database) *SectionDetailRepository {
    collection := db.Collection("sectionDetails")
    return &SectionDetailRepository{
        collection: collection,
    }
}

func (r *SectionDetailRepository) Create(ctx context.Context, sectionDetail *model.SectionDetail) error {
    if sectionDetail.ID.IsZero() {
        sectionDetail.ID = primitive.NewObjectID()
    }
    
    sectionDetail.CreatedAt = primitive.NewDateTimeFromTime(time.Now()).Time()
    sectionDetail.UpdatedAt = sectionDetail.CreatedAt
    
    _, err := r.collection.InsertOne(ctx, sectionDetail)
    if mongo.IsDuplicateKeyError(err) {
        return errors.New("section detail already exists")
    }
    return err
}

func (r *SectionDetailRepository) Find(ctx context.Context, id string) (*model.SectionDetail, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, errors.New("invalid id format")
    }

    var sectionDetail model.SectionDetail
    err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&sectionDetail)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("section detail not found")
        }
        return nil, err
    }

    return &sectionDetail, nil
}

func (r *SectionDetailRepository) FindAll(ctx context.Context, opts ...*options.FindOptions) ([]*model.SectionDetail, error) {
    cursor, err := r.collection.Find(ctx, bson.M{}, opts...)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var sectionDetails []*model.SectionDetail
    if err = cursor.All(ctx, &sectionDetails); err != nil {
        return nil, err
    }

    if sectionDetails == nil {
        sectionDetails = make([]*model.SectionDetail, 0)
    }

    return sectionDetails, nil
}

func (r *SectionDetailRepository) Update(ctx context.Context, sectionDetail *model.SectionDetail) error {
    if sectionDetail.ID.IsZero() {
        return errors.New("invalid section detail id")
    }

    sectionDetail.UpdatedAt = primitive.NewDateTimeFromTime(time.Now()).Time()

    update := bson.M{
        "$set": bson.M{
            "order":       sectionDetail.Order,
            "title":       sectionDetail.Title,
            "sectionType": sectionDetail.SectionType,
            "description": sectionDetail.Description,
            "updated_at":  sectionDetail.UpdatedAt,
        },
    }

    result, err := r.collection.UpdateOne(
        ctx,
        bson.M{"_id": sectionDetail.ID},
        update,
    )
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return errors.New("section detail not found")
    }

    return nil
}

func (r *SectionDetailRepository) Delete(ctx context.Context, id string) error {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return errors.New("invalid id format")
    }

    result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
    if err != nil {
        return err
    }

    if result.DeletedCount == 0 {
        return errors.New("section detail not found")
    }

    return nil
}

//new method impl
func (r *SectionDetailRepository) FindBySectionType(ctx context.Context, sectionType string) ([]*model.SectionDetail, error) {
    filter := bson.M{"sectionType": sectionType}
    
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, errors.New("error finding sections by type: %v")
    }
    defer cursor.Close(ctx)

    var sections []*model.SectionDetail
    if err = cursor.All(ctx, &sections); err != nil {
        return nil, errors.New("error decoding sections: %v")
    }

    return sections, nil
}
