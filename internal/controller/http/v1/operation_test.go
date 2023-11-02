package v1

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"
	mocks "usersegments/internal/controller/http/v1/mocks"
	"usersegments/internal/entity"
	"usersegments/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOperationHandler_getOperationHistory(t *testing.T) {
	type args struct {
		ctx     context.Context
		input   service.GetOperationHistoryInput
		user_id int
	}

	type MockBehaviour func(m *mocks.MockOperationService, args args, output []entity.Operation)
	tests := []struct {
		name                 string
		args                 args
		output               []entity.Operation
		input                string
		mockBehaviour        MockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: service.GetOperationHistoryInput{
					Month: "10",
					Year:  "2023",
				},
			},
			output: []entity.Operation{
				{
					Id:            1,
					UserID:        1,
					SegmentName:   "test1",
					OperationType: "Add",
					OperationDate: time.UnixMilli(123456),
				},
				{
					Id:            2,
					UserID:        1,
					SegmentName:   "test1",
					OperationType: "Delete",
					OperationDate: time.UnixMilli(123456),
				},
			},
			input: `{"month":"10","year":"2023"}`,
			mockBehaviour: func(m *mocks.MockOperationService, args args, output []entity.Operation) {
				m.EXPECT().GetOperationHistory(args.ctx, args.user_id, args.input).Return(output, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"data":[{"id":1,"user_id":1,"segment_name":"test1","operation_type":"Add","operation_date":"1970-01-01T03:02:03.456+03:00"},{"id":2,"user_id":1,"segment_name":"test1","operation_type":"Delete","operation_date":"1970-01-01T03:02:03.456+03:00"}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			operationService := mocks.NewMockOperationService(c)
			tt.mockBehaviour(operationService, tt.args, tt.output)

			services := &OperationHandler{operationService}
			handler := Handler{OperationHandler: services}

			r := gin.New()

			r.POST("/api/v1/operations/:id", handler.OperationHandler.getOperationHistory)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/operations/%d", tt.args.user_id), bytes.NewBufferString(tt.input))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}
