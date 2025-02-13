// internal/service/user.go
package service

import (
    "context"
    "time"

    "golang.org/x/crypto/bcrypt"

    "myapp/internal/model"
    "myapp/internal/repository"
    "myapp/pkg/validator"
)

type UserService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
    return &UserService{
        repo: repo,
    }
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
    if err := validator.Validate(user); err != nil {
        return err
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user.Password = string(hashedPassword)
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    return s.repo.Create(ctx, user)
}

func (s *UserService) Find(ctx context.Context, id string) (*model.User, error) {
    return s.repo.Find(ctx, id)
}

func (s *UserService) FindAll(ctx context.Context) ([]*model.User, error) {
    return s.repo.FindAll(ctx)
}

func (s *UserService) Update(ctx context.Context, user *model.User) error {
    if err := validator.Validate(user); err != nil {
        return err
    }

    user.UpdatedAt = time.Now()
    return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}