package v1

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	mocks "usersegments/internal/controller/http/v1/mocks"
	"usersegments/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_createUser(t *testing.T) {

	type args struct {
		ctx   context.Context
		input service.UserCreateInput
	}
	type MockBehaviour func(m *mocks.MockUserService, args args)
	tests := []struct {
		name                 string
		args                 args
		input                string
		mockBehaviour        MockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: service.UserCreateInput{
					Name: "test",
				},
			},
			input: `{"name": "test"}`,
			mockBehaviour: func(m *mocks.MockUserService, args args) {
				m.EXPECT().CreateUser(args.ctx, args.input).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"data":{"ID":1}}`,
		},
		{
			name: "User Already Exists",
			args: args{
				ctx: context.Background(),
				input: service.UserCreateInput{
					Name: "test",
				},
			},
			input: `{"name": "test"}`,
			mockBehaviour: func(m *mocks.MockUserService, args args) {
				m.EXPECT().CreateUser(args.ctx, args.input).Return(0, service.ErrUserAlreadyExists)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"user already exists"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userService := mocks.NewMockUserService(c)
			tt.mockBehaviour(userService, tt.args)

			services := &UserHandler{userService}
			handler := Handler{UserHandler: services}

			r := gin.New()

			r.POST("/api/v1/users", handler.UserHandler.createUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(tt.input))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestUserHandler_addUserToSegment(t *testing.T) {

	type args struct {
		ctx   context.Context
		input service.AddUserToSegmentsInput
	}
	type MockBehaviour func(m *mocks.MockUserService, args args)
	tests := []struct {
		name                 string
		args                 args
		input                string
		mockBehaviour        MockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: service.AddUserToSegmentsInput{
					UserID:   1,
					Segments: []string{"test1", "test2"},
					Timeout:  0,
				},
			},
			input: `{"segments":["test1", "test2"]}`,
			mockBehaviour: func(m *mocks.MockUserService, args args) {
				m.EXPECT().AddUserToSegment(args.ctx, args.input).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"data":"User Added"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userService := mocks.NewMockUserService(c)
			tt.mockBehaviour(userService, tt.args)

			services := &UserHandler{userService}
			handler := Handler{UserHandler: services}

			r := gin.New()

			r.POST("/api/v1/users/:id", handler.UserHandler.addUserToSegment)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/users/%d", tt.args.input.UserID), bytes.NewBufferString(tt.input))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestUserHandler_deleteUserFromSegment(t *testing.T) {
	type args struct {
		ctx   context.Context
		input service.DeleteUserFromSegmentsInput
	}
	type MockBehaviour func(m *mocks.MockUserService, args args)
	tests := []struct {
		name                 string
		args                 args
		input                string
		mockBehaviour        MockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: service.DeleteUserFromSegmentsInput{
					UserID:   1,
					Segments: []string{"test1", "test2"},
				},
			},
			input: `{"segments":["test1", "test2"]}`,
			mockBehaviour: func(m *mocks.MockUserService, args args) {
				m.EXPECT().DeleteUserFromSegment(args.ctx, args.input).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"data":"User Deleted"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userService := mocks.NewMockUserService(c)
			tt.mockBehaviour(userService, tt.args)

			services := &UserHandler{userService}
			handler := Handler{UserHandler: services}

			r := gin.New()

			r.POST("/api/v1/users/:id", handler.UserHandler.deleteUserFromSegment)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/users/%d", tt.args.input.UserID), bytes.NewBufferString(tt.input))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestUserHandler_getAllUserSegments(t *testing.T) {
	type args struct {
		ctx   context.Context
		input int
	}

	type MockBehaviour func(m *mocks.MockUserService, args args, output []service.AllUserSegmentsOutput)
	tests := []struct {
		name                 string
		args                 args
		output               []service.AllUserSegmentsOutput
		mockBehaviour        MockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			args: args{
				ctx:   context.Background(),
				input: 1,
			},
			output: []service.AllUserSegmentsOutput{
				{
					Name: "test1",
				},
				{
					Name: "test2",
				},
			},
			mockBehaviour: func(m *mocks.MockUserService, args args, output []service.AllUserSegmentsOutput) {
				m.EXPECT().GetAllUserSegments(args.ctx, args.input).Return(output, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"data":[{"name":"test1"},{"name":"test2"}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userService := mocks.NewMockUserService(c)
			tt.mockBehaviour(userService, tt.args, tt.output)

			services := &UserHandler{userService}
			handler := Handler{UserHandler: services}

			r := gin.New()

			r.POST("/api/v1/users/:id", handler.UserHandler.getAllUserSegments)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/users/%d", tt.args.input), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}
