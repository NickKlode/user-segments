package v1

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"
	mocks "usersegments/internal/controller/http/v1/mocks"
	"usersegments/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSegmentHandler_createSegment(t *testing.T) {

	type args struct {
		ctx   context.Context
		input service.SegmentCreateInput
	}

	type MockBehaviour func(m *mocks.MockSegmentService, args args)
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
				input: service.SegmentCreateInput{
					Name: "test",
				},
			},
			input: `{"name": "test"}`,
			mockBehaviour: func(m *mocks.MockSegmentService, args args) {
				m.EXPECT().CreateSegment(args.ctx, args.input).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"data":{"ID":1}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			segmentService := mocks.NewMockSegmentService(c)
			tt.mockBehaviour(segmentService, tt.args)

			services := &SegmentHandler{segmentService}
			handler := Handler{SegmentHandler: services}

			r := gin.New()

			r.POST("/api/v1/segments", handler.SegmentHandler.createSegment)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/segments", bytes.NewBufferString(tt.input))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestSegmentHandler_deleteSegment(t *testing.T) {
	type args struct {
		ctx   context.Context
		input service.SegmentDeleteInput
	}

	type MockBehaviour func(m *mocks.MockSegmentService, args args)
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
				input: service.SegmentDeleteInput{
					Name: "test",
				},
			},
			input: `{"name": "test"}`,
			mockBehaviour: func(m *mocks.MockSegmentService, args args) {
				m.EXPECT().DeleteSegment(args.ctx, args.input).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"data":"Segment deleted"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			segmentService := mocks.NewMockSegmentService(c)
			tt.mockBehaviour(segmentService, tt.args)

			services := &SegmentHandler{segmentService}
			handler := Handler{SegmentHandler: services}

			r := gin.New()

			r.DELETE("/api/v1/segments", handler.SegmentHandler.deleteSegment)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/v1/segments", bytes.NewBufferString(tt.input))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}
