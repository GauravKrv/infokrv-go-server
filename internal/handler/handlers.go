// internal/handler/handlers.go
package handler

import (
    "myapp/internal/service"
)

type Handlers struct {
    User          *UserHandler
    SectionDetail *SectionDetailHandler
}

func NewHandlers(services *service.Services) *Handlers {
    return &Handlers{
        User:          NewUserHandler(services.User),
        SectionDetail: NewSectionDetailHandler(services.SectionDetail),
    }
}