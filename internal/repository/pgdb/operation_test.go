package pgdb

import (
	"context"
	"testing"
	"time"
	"usersegments/internal/entity"
	"usersegments/pkg/postgres"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestOperationRepo_GetOperationHistory(t *testing.T) {
	type args struct {
		ctx     context.Context
		user_id int
		month   string
		year    string
	}

	type MockBehaviour func(m pgxmock.PgxPoolIface, args args)
	tests := []struct {
		name          string
		args          args
		mockBehaviour MockBehaviour
		want          []entity.Operation
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				ctx:     context.Background(),
				user_id: 1,
				month:   "01",
				year:    "1970",
			},
			mockBehaviour: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id", "user_id", "segment_name", "operation_type", "operation_date"}).
					AddRow(1, 1, "test1", "Add", time.UnixMilli(123456)).AddRow(2, 1, "test1", "Delete", time.UnixMilli(123456))
				m.ExpectQuery("SELECT id, user_id, segment_name, operation_type, operation_date FROM operations").
					WithArgs(args.user_id, args.month, args.year).WillReturnRows(rows)
			},
			want: []entity.Operation{
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tt.mockBehaviour(poolMock, tt.args)

			psqlMock := &OperationRepo{
				db: &postgres.Postgres{
					Pool: poolMock,
				},
			}
			got, err := psqlMock.GetOperationHistory(tt.args.ctx, tt.args.user_id, tt.args.month, tt.args.year)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
