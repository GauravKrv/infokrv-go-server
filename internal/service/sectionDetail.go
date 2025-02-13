// internal/service/sectionDetail.go
package service

import (
    "context"
    "time"

    "myapp/internal/model"
    "myapp/internal/repository"
    "myapp/pkg/validator"
)

type SectionDetailService struct {
    repo repository.SectionDetailRepository
}

func NewSectionDetailService(repo repository.SectionDetailRepository) *SectionDetailService {
    return &SectionDetailService{
        repo: repo,
    }
}

func (s *SectionDetailService) Create(ctx context.Context, sectionDetail *model.SectionDetail) error {
    if err := validator.Validate(sectionDetail); err != nil {
        return err
    }
    sectionDetail.CreatedAt = time.Now()
    sectionDetail.UpdatedAt = time.Now()

    return s.repo.Create(ctx, sectionDetail)
}

func (s *SectionDetailService) Find(ctx context.Context, id string) (*model.SectionDetail, error) {
    return s.repo.Find(ctx, id)
}

func (s *SectionDetailService) FindAll(ctx context.Context) ([]*model.SectionDetail, error) {
    return s.repo.FindAll(ctx)
}

func (s *SectionDetailService) Update(ctx context.Context, sectionDetail *model.SectionDetail) error {
    if err := validator.Validate(sectionDetail); err != nil {
        return err
    }

    sectionDetail.UpdatedAt = time.Now()
    return s.repo.Update(ctx, sectionDetail)
}

func (s *SectionDetailService) Delete(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}

// func (s *SectionDetailService) FindBySectionType(ctx context.Context, sectionType string)  ([]*model.SectionDetail, error) {
//     return s.repo.FindBySectionType(ctx, sectionType)
// }

// Update the service method to return DTOs
func (s *SectionDetailService) FindBySectionType(ctx context.Context, sectionType string) ([]*model.SectionDetailDTO, error) {
    sections, err := s.repo.FindBySectionType(ctx, sectionType)
    if err != nil {
        return nil, err
    }

    // Convert slice of SectionDetail to slice of DTOs
    dtos := make([]*model.SectionDetailDTO, len(sections))
    for i, section := range sections {
        dtos[i] = section.ToDTO()
    }

    return dtos, nil
}