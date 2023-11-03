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

type SegmentRepo struct {
	db *postgres.Postgres
}

func NewSegmentRepo(db *postgres.Postgres) *SegmentRepo {
	return &SegmentRepo{db: db}
}

func (s *SegmentRepo) CreateSegment(ctx context.Context, segment entity.Segment) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, percent) VALUES ($1, $2) RETURNING id", segmentsTable)

	var id int

	err := s.db.Pool.QueryRow(ctx, query, segment.Name, segment.Percent).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrors.ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("pgdb - CreateSegment - s.db.QueryRow. %v", err)
	}
	return id, nil
}

func (s *SegmentRepo) DeleteSegment(ctx context.Context, name string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE name=$1", segmentsTable)

	coTag, err := s.db.Pool.Exec(ctx, query, name)
	if err != nil {
		return fmt.Errorf("pgdb - DeleteSegment - s.db.Exec. %v", err)
	}
	if coTag.RowsAffected() == 0 {
		return repoerrors.ErrSegmentNotFound
	}

	return nil
}
