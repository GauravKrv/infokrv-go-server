// internal/model/sectionDetail.go
package model

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type SectionDetail struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Order       int               `bson:"order" json:"order"`
    Title       string            `bson:"title" json:"title" validate:"required,min=2,max=100"`
    SectionType string            `bson:"sectionType" json:"sectionType" validate:"required,min=2,max=100"`
    Description string            `bson:"description" json:"description"`
    CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}

// Add this DTO struct to your model package
type SectionDetailDTO struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    SectionType string    `json:"sectionType"`
    Description string    `json:"description"`
}

// Add this helper method to convert SectionDetail to DTO
func (s *SectionDetail) ToDTO() *SectionDetailDTO {
    return &SectionDetailDTO{
        ID:          s.ID.Hex(),
        Title:       s.Title,
        SectionType: s.SectionType,
        Description: s.Description,
    }
}

