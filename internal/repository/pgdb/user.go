package pgdb

import (
	"context"
	"errors"
	"fmt"
	"usersegments/internal/entity"
	"usersegments/internal/repository/repoerrors"
	"usersegments/pkg/postgres"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepo struct {
	db *postgres.Postgres
}

func NewUserRepo(db *postgres.Postgres) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateUser(ctx context.Context, user entity.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", usersTable)

	var id int

	err := u.db.Pool.QueryRow(ctx, query, user.Name).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrors.ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("pgdb - CreateUser - u.db.QueryRow. %v", err)
	}
	return id, nil
}

func (o *UserRepo) AddUserToSegment(ctx context.Context, user_id, timeout int, segment string) error {

	tx, err := o.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pgdb - AddUserToSegment - o.db.Begin. %v", err)
	}
	defer tx.Rollback(ctx)
	query := fmt.Sprintf("INSERT INTO %s(user_id, segment_id, timeout) VALUES($1, (SELECT id FROM %s WHERE name=$2), $3)", usersSegmentsTable, segmentsTable)

	_, err = tx.Exec(ctx, query, user_id, segment, timeout)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return repoerrors.ErrAlreadyExists
			}
		}
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23503" {
				return repoerrors.ErrUserNotFound
			}
		}
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23502" {
				return repoerrors.ErrSegmentNotFound
			}
		}

		return fmt.Errorf("pgdb - AddUserToSegment - tx.Exec 1. %v", err)
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, segment_name, operation_type) VALUES ($1, $2, $3)", operationsTable)
	_, err = tx.Exec(ctx, query, user_id, segment, entity.OperationTypeAdd)
	if err != nil {
		return fmt.Errorf("pgdb - AddUserToSegment - tx.Exec 2. %v", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("pgdb - AddUserToSegment - tx.Commit. %v", err)
	}

	return nil
}

func (o *UserRepo) DeleteUserFromSegment(ctx context.Context, user_id int, segment string) error {
	tx, err := o.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pgdb - DeleteUserFromSegment - o.db.Begin. %v", err)
	}
	defer tx.Rollback(ctx)

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1 AND segment_id=(SELECT id FROM %s WHERE name=$2)", usersSegmentsTable, segmentsTable)

	coTag, err := tx.Exec(ctx, query, user_id, segment)
	if err != nil {
		return fmt.Errorf("pgdb - DeleteUserFromSegment - tx.Exec 1. %v", err)
	}
	if coTag.RowsAffected() == 0 {
		return repoerrors.ErrUserNotInSegment
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, segment_name, operation_type) VALUES ($1, $2, $3)", operationsTable)
	_, err = tx.Exec(ctx, query, user_id, segment, entity.OperationTypeDelete)
	if err != nil {
		return fmt.Errorf("pgdb - DeleteUserFromSegment - tx.Exec 2. %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("pgdb - DeleteUserFromSegment - tx.Commit. %v", err)
	}
	return nil
}

func (s *UserRepo) GetAllUserSegments(ctx context.Context, user_id int) ([]entity.Segment, error) {
	query := fmt.Sprintf(`SELECT id, name, percent FROM %s 
	JOIN %s 
	ON %s.id=%s.segment_id 
	AND user_id=$1`, segmentsTable, usersSegmentsTable, segmentsTable, usersSegmentsTable)

	rows, err := s.db.Pool.Query(ctx, query, user_id)
	if err != nil {
		return nil, fmt.Errorf("pgdb - GetAllUserSegments - s.db.Query. %v", err)
	}
	defer rows.Close()

	var segments []entity.Segment
	for rows.Next() {
		var segment entity.Segment
		err := rows.Scan(
			&segment.Id,
			&segment.Name,
			&segment.Percent,
		)
		if err != nil {
			return nil, fmt.Errorf("pgdb - GetAllUserSegments - rows.Scan. %v", err)
		}
		segments = append(segments, segment)
	}
	if segments == nil {
		return nil, repoerrors.ErrUserNotFound
	}
	return segments, nil
}
