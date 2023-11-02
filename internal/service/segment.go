package service

import (
	"context"
	"errors"
	"usersegments/internal/entity"
	"usersegments/internal/repository/repoerrors"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=segment.go -destination=mocks/segment_mock.go
type SegmentRepo interface {
	CreateSegment(ctx context.Context, segment entity.Segment) (int, error)
	DeleteSegment(ctx context.Context, name string) error
}
type SegmentService struct {
	segmentRepo SegmentRepo
}

func NewSegmentService(segmentRepo SegmentRepo) *SegmentService {
	return &SegmentService{segmentRepo: segmentRepo}
}

type SegmentCreateInput struct {
	Name    string
	Percent int
}

type SegmentWithPercent struct {
	Name    string
	Percent int
}

var SegmentsWithPercent []SegmentWithPercent

func (s *SegmentService) CreateSegment(ctx context.Context, input SegmentCreateInput) (int, error) {
	segment := entity.Segment{
		Name:    input.Name,
		Percent: input.Percent,
	}

	segmentId, err := s.segmentRepo.CreateSegment(ctx, segment)
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			return 0, ErrSegmentAlreadyExists
		}
		logrus.Errorf("service - s.segmentRepo.CreateSegment. %v", err)
		return 0, ErrCannotCreateSegment
	}
	if input.Percent != 0 {
		SegmentsWithPercent = append(SegmentsWithPercent, SegmentWithPercent{
			Name:    input.Name,
			Percent: input.Percent,
		})
	}

	return segmentId, nil
}

type SegmentDeleteInput struct {
	Name string
}

func (s *SegmentService) DeleteSegment(ctx context.Context, input SegmentDeleteInput) error {
	err := s.segmentRepo.DeleteSegment(ctx, input.Name)
	if err != nil {
		if errors.Is(err, repoerrors.ErrSegmentNotFound) {
			return ErrSegmentNotFound
		}
		logrus.Errorf("service - s.segmentRepo.DeleteSegment. %v", err)
		return ErrCannotDeleteSegment
	}
	return nil
}
