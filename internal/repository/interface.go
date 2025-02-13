// internal/repository/interface.go
package repository

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo/options"
    "myapp/internal/model"
)

type SectionDetailRepository interface {
    Create(ctx context.Context, sectionDetail *model.SectionDetail) error
    Find(ctx context.Context, id string) (*model.SectionDetail, error)
    FindAll(ctx context.Context, opts ...*options.FindOptions) ([]*model.SectionDetail, error)
    Update(ctx context.Context, sectionDetail *model.SectionDetail) error
    Delete(ctx context.Context, id string) error
    FindBySectionType(ctx context.Context, sectionType string) ([]*model.SectionDetail, error)
}

type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    Find(ctx context.Context, id string) (*model.User, error)
    FindAll(ctx context.Context) ([]*model.User, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id string) error
}

