package pgdb

import (
	"context"
	"fmt"
	"usersegments/internal/entity"
	"usersegments/internal/repository/repoerrors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPoolOperation interface {
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
type OperationRepo struct {
	db PgxPoolOperation
}

func NewOperationRepo(db PgxPoolOperation) *OperationRepo {
	return &OperationRepo{db: db}
}

func (o *OperationRepo) GetOperationHistory(ctx context.Context, user_id int, month, year string) ([]entity.Operation, error) {
	query := fmt.Sprintf(`SELECT id, user_id, segment_name, operation_type, operation_date FROM %s 
	WHERE user_id=$1 
	AND date_part('month', operation_date) = $2 
	AND date_part('year', operation_date) = $3`, operationsTable)
	rows, err := o.db.Query(ctx, query, user_id, month, year)
	if err != nil {
		return nil, fmt.Errorf("pgdb - GetOperationHistory - o.db.Query. %v", err)
	}
	defer rows.Close()

	var operations []entity.Operation
	for rows.Next() {
		var operation entity.Operation
		err := rows.Scan(
			&operation.Id,
			&operation.UserID,
			&operation.SegmentName,
			&operation.OperationType,
			&operation.OperationDate,
		)
		if err != nil {
			return nil, fmt.Errorf("pgdb - GetOperationHistory - rows.Scan. %v", err)
		}
		operations = append(operations, operation)
	}
	if operations == nil {
		return nil, repoerrors.ErrOperationHistoryNotFound
	}

	return operations, nil
}
