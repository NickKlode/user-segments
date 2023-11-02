package repository

import (
	"usersegments/internal/repository/pgdb"
)

type Deps struct {
	UserPgx      pgdb.PgxPoolUser
	SegmentPgx   pgdb.PgxPoolSegment
	OperationPgx pgdb.PgxPoolOperation
}

type Repository struct {
	UserRepo      *pgdb.UserRepo
	SegmentRepo   *pgdb.SegmentRepo
	OperationRepo *pgdb.OperationRepo
}

func NewRepository(deps Deps) *Repository {
	return &Repository{
		UserRepo:      pgdb.NewUserRepo(deps.UserPgx),
		SegmentRepo:   pgdb.NewSegmentRepo(deps.SegmentPgx),
		OperationRepo: pgdb.NewOperationRepo(deps.OperationPgx),
	}
}
