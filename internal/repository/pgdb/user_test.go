package pgdb

import (
	"context"
	"testing"
	"usersegments/internal/entity"
	"usersegments/pkg/postgres"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_CreateUser(t *testing.T) {

	type args struct {
		ctx  context.Context
		user entity.User
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
				user: entity.User{
					Name: "test",
				},
			},
			mockBehaviour: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

				m.ExpectQuery("INSERT INTO users").WithArgs(args.user.Name).WillReturnRows(rows)
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

			psqlMock := &UserRepo{
				db: &postgres.Postgres{
					Pool: poolMock,
				},
			}
			got, err := psqlMock.CreateUser(tt.args.ctx, tt.args.user)
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

func TestUserRepo_AddUserToSegment(t *testing.T) {
	type args struct {
		ctx          context.Context
		user_id      int
		timeout      int
		segment_name string
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
				ctx:          context.Background(),
				user_id:      1,
				timeout:      0,
				segment_name: "test",
			},
			mockBehaviour: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO users_segments").WithArgs(args.user_id, args.segment_name, args.timeout).WillReturnResult(pgxmock.NewResult("", 1))

				m.ExpectExec("INSERT INTO operations").WithArgs(args.user_id, args.segment_name, entity.OperationTypeAdd).WillReturnResult(pgxmock.NewResult("", 1))
				m.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tt.mockBehaviour(poolMock, tt.args)

			psqlMock := &UserRepo{
				db: &postgres.Postgres{
					Pool: poolMock,
				},
			}
			err := psqlMock.AddUserToSegment(tt.args.ctx, tt.args.user_id, tt.args.timeout, tt.args.segment_name)
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

func TestUserRepo_DeleteUserFromSegment(t *testing.T) {
	type args struct {
		ctx          context.Context
		user_id      int
		segment_name string
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
				ctx:          context.Background(),
				user_id:      1,
				segment_name: "test",
			},
			mockBehaviour: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()
				m.ExpectExec("DELETE FROM users_segments").WithArgs(args.user_id, args.segment_name).WillReturnResult(pgxmock.NewResult("", 1))

				m.ExpectExec("INSERT INTO operations").WithArgs(args.user_id, args.segment_name, entity.OperationTypeDelete).WillReturnResult(pgxmock.NewResult("", 1))
				m.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tt.mockBehaviour(poolMock, tt.args)

			psqlMock := &UserRepo{
				db: &postgres.Postgres{
					Pool: poolMock,
				},
			}
			err := psqlMock.DeleteUserFromSegment(tt.args.ctx, tt.args.user_id, tt.args.segment_name)
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

func TestUserRepo_GetAllUserSegments(t *testing.T) {
	type args struct {
		ctx     context.Context
		user_id int
	}
	type MockBehaviour func(m pgxmock.PgxPoolIface, args args)
	tests := []struct {
		name          string
		args          args
		mockBehaviour MockBehaviour
		want          []entity.Segment
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				ctx:     context.Background(),
				user_id: 1,
			},
			mockBehaviour: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id", "name", "percent"}).AddRow(1, "test1", 0).AddRow(2, "test2", 10)
				m.ExpectQuery("SELECT id, name, percent FROM segments").WithArgs(args.user_id).WillReturnRows(rows)
			},
			want: []entity.Segment{
				{
					Id:      1,
					Name:    "test1",
					Percent: 0,
				},
				{
					Id:      2,
					Name:    "test2",
					Percent: 10,
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

			psqlMock := &UserRepo{
				db: &postgres.Postgres{
					Pool: poolMock,
				},
			}
			got, err := psqlMock.GetAllUserSegments(tt.args.ctx, tt.args.user_id)
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
