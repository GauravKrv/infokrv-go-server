package service

import (
     "myapp/internal/repository"
)

type Services struct {
    User          *UserService
    SectionDetail *SectionDetailService
}

func NewServices(repos *repository.Repositories) *Services {
    return &Services{
        User:          NewUserService(repos.User),
        SectionDetail: NewSectionDetailService(repos.SectionDetail),
    }
}