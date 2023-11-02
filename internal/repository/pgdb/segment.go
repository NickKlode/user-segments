package pgdb

import (
	"context"
	"errors"
	"fmt"
	"usersegments/internal/entity"
	"usersegments/internal/repository/repoerrors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPoolSegment interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

type SegmentRepo struct {
	db PgxPoolSegment
}

func NewSegmentRepo(db PgxPoolSegment) *SegmentRepo {
	return &SegmentRepo{db: db}
}

func (s *SegmentRepo) CreateSegment(ctx context.Context, segment entity.Segment) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, percent) VALUES ($1, $2) RETURNING id", segmentsTable)

	var id int

	err := s.db.QueryRow(ctx, query, segment.Name, segment.Percent).Scan(&id)
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

	coTag, err := s.db.Exec(ctx, query, name)
	if err != nil {
		return fmt.Errorf("pgdb - DeleteSegment - s.db.Exec. %v", err)
	}
	if coTag.RowsAffected() == 0 {
		return repoerrors.ErrSegmentNotFound
	}

	return nil
}
