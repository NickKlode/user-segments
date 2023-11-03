package pgdb

import (
	"context"
	"fmt"
	"usersegments/internal/entity"
	"usersegments/internal/repository/repoerrors"
	"usersegments/pkg/postgres"
)

type OperationRepo struct {
	db *postgres.Postgres
}

func NewOperationRepo(db *postgres.Postgres) *OperationRepo {
	return &OperationRepo{db: db}
}

func (o *OperationRepo) GetOperationHistory(ctx context.Context, user_id int, month, year string) ([]entity.Operation, error) {
	query := fmt.Sprintf(`SELECT id, user_id, segment_name, operation_type, operation_date FROM %s 
	WHERE user_id=$1 
	AND date_part('month', operation_date) = $2 
	AND date_part('year', operation_date) = $3`, operationsTable)
	rows, err := o.db.Pool.Query(ctx, query, user_id, month, year)

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
