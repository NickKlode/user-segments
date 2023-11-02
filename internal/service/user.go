package service

import (
	"context"
	"errors"
	"time"
	"usersegments/internal/entity"
	"usersegments/internal/repository/repoerrors"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=user.go -destination=mocks/user_mock.go
type UserRepo interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	AddUserToSegment(ctx context.Context, user_id, timeout int, segment string) error
	DeleteUserFromSegment(ctx context.Context, user_id int, segment string) error
	GetAllUserSegments(ctx context.Context, user_id int) ([]entity.Segment, error)
}

type UserService struct {
	userRepo UserRepo
}

func NewUserService(userRepo UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type UserCreateInput struct {
	Name string
}

var userCounter int

func (u *UserService) CreateUser(ctx context.Context, input UserCreateInput) (int, error) {
	user := entity.User{
		Name: input.Name,
	}

	userId, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			return 0, ErrUserAlreadyExists
		}
		logrus.Errorf("service - u.userRepo.CreateUser. %v", err)
		return 0, ErrCannotCreateUser
	}
	if userCounter == 100 {
		userCounter = 0
	}
	userCounter++

	go func() {
		for _, v := range SegmentsWithPercent {
			if userCounter%int(100/v.Percent) == 0 {
				err := u.userRepo.AddUserToSegment(context.Background(), userId, 0, v.Name)
				if err != nil {
					logrus.Errorf("service - CreateUser - u.userRepo.AddUserToSegment. %v", err)
					return
				}
			}
		}
	}()

	return userId, nil
}

type AddUserToSegmentsInput struct {
	UserID   int
	Segments []string
	Timeout  int
}

func (u *UserService) AddUserToSegment(ctx context.Context, input AddUserToSegmentsInput) error {
	var resErr error
	var isUserNotInAllSegments bool

	for _, seg := range input.Segments {
		err := u.userRepo.AddUserToSegment(ctx, input.UserID, input.Timeout, seg)
		if err != nil {
			if errors.Is(err, repoerrors.ErrAlreadyExists) {
				resErr = ErrUserAlreadyInSegments
			} else if errors.Is(err, repoerrors.ErrUserNotFound) {
				return ErrUserNotFound
			} else if errors.Is(err, repoerrors.ErrSegmentNotFound) {
				return ErrSegmentNotFound
			} else {
				logrus.Errorf("service - u.userRepo.AddUserToSegment. %v", err)
				return ErrCannotAddUser
			}

		} else {
			isUserNotInAllSegments = true
		}

	}
	if !isUserNotInAllSegments {
		return resErr
	}
	if input.Timeout != 0 {
		go func() {
			time.Sleep(time.Second * time.Duration(input.Timeout))
			for _, seg := range input.Segments {
				err := u.userRepo.DeleteUserFromSegment(context.Background(), input.UserID, seg)
				if err != nil {
					logrus.Errorf("service - AddUserToSegment - u.userRepo.DeleteUserFromSegment. %v", err)
					return
				}

			}
		}()
	}

	return nil
}

type DeleteUserFromSegmentsInput struct {
	UserID   int
	Segments []string
}

func (u *UserService) DeleteUserFromSegment(ctx context.Context, input DeleteUserFromSegmentsInput) error {
	var resErr error
	var isUserNotInAllSegments bool
	for _, seg := range input.Segments {
		err := u.userRepo.DeleteUserFromSegment(ctx, input.UserID, seg)
		if err != nil {
			if errors.Is(err, repoerrors.ErrUserNotInSegment) {
				resErr = ErrUserNotInSegment
			} else {
				logrus.Errorf("service - u.userRepo.DeleteUserFromSegment. %v", err)
				return ErrCannotDeleteUserFromSegment
			}
		} else {
			isUserNotInAllSegments = true
		}
	}
	if !isUserNotInAllSegments {
		return resErr
	}
	return nil
}

type AllUserSegmentsOutput struct {
	Name string `json:"name"`
}

func (u *UserService) GetAllUserSegments(ctx context.Context, user_id int) ([]AllUserSegmentsOutput, error) {
	segments, err := u.userRepo.GetAllUserSegments(ctx, user_id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrUserNotFound) {
			return nil, ErrNoSegments
		}
		logrus.Errorf("service - u.userRepo.GetAllUserSegments. %v", err)
		return nil, ErrCannotGetUserSegmets
	}

	output := make([]AllUserSegmentsOutput, 0, len(segments))
	for _, segment := range segments {
		output = append(output, AllUserSegmentsOutput{
			Name: segment.Name,
		})
	}
	if output == nil {
		return nil, ErrNoSegments
	}
	return output, nil
}
