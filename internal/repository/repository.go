package repository

import (
	"usersegments/internal/repository/pgdb"
	"usersegments/pkg/postgres"
)

type Repository struct {
	UserRepo      *pgdb.UserRepo
	SegmentRepo   *pgdb.SegmentRepo
	OperationRepo *pgdb.OperationRepo
}

func NewRepository(pg *postgres.Postgres) *Repository {
	return &Repository{
		UserRepo:      pgdb.NewUserRepo(pg),
		SegmentRepo:   pgdb.NewSegmentRepo(pg),
		OperationRepo: pgdb.NewOperationRepo(pg),
	}
}
