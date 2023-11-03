package pgdb

import (
	"context"
	"testing"
	"usersegments/internal/entity"
	"usersegments/pkg/postgres"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestSegmentRepo_CreateSegment(t *testing.T) {
	type args struct {
		ctx     context.Context
		segment entity.Segment
	}

	type MockBehaviour func(m pgxmock.PgxPoolIface, args args)
	tests := []struct {
		name          string
		args          args
		mockBehaviour MockBehaviour
		want          int
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				segment: entity.Segment{
					Name:    "test1",
					Percent: 0,
				},
			},
			mockBehaviour: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery("INSERT INTO segments").WithArgs(args.segment.Name, args.segment.Percent).WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tt.mockBehaviour(poolMock, tt.args)

			psqlMock := &SegmentRepo{
				db: &postgres.Postgres{
					Pool: poolMock,
				},
			}
			got, err := psqlMock.CreateSegment(tt.args.ctx, tt.args.segment)
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

func TestSegmentRepo_DeleteSegment(t *testing.T) {
	type args struct {
		ctx     context.Context
		segment string
	}

	type MockBehaviour func(m pgxmock.PgxPoolIface, args args)
	tests := []struct {
		name          string
		args          args
		mockBehaviour MockBehaviour
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				ctx:     context.Background(),
				segment: "test",
			},
			mockBehaviour: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec("DELETE FROM segments").WithArgs(args.segment).WillReturnResult(pgxmock.NewResult("", 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tt.mockBehaviour(poolMock, tt.args)

			psqlMock := &SegmentRepo{
				db: &postgres.Postgres{
					Pool: poolMock,
				},
			}
			err := psqlMock.DeleteSegment(tt.args.ctx, tt.args.segment)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
