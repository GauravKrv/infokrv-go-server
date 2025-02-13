// internal/repository/repositories.go
package repository

import (
    "go.mongodb.org/mongo-driver/mongo"
    "myapp/internal/repository/mongodb"
)

type Repositories struct {
    User          UserRepository
    SectionDetail SectionDetailRepository
}

func NewRepositories(db *mongo.Database) (*Repositories, error) {
    userRepo := mongodb.NewUserRepository(db)
    sectionRepo := mongodb.NewSectionDetailRepository(db)

    return &Repositories{
        User:          userRepo,
        SectionDetail: sectionRepo,
    }, nil
}

