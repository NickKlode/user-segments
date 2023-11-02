package entity

import "time"

type Operation struct {
	Id            int       `db:"id" json:"id"`
	UserID        int       `db:"user_id" json:"user_id"`
	SegmentName   string    `db:"segment_name" json:"segment_name"`
	OperationType string    `db:"operation_type" json:"operation_type"`
	OperationDate time.Time `db:"operation_date" json:"operation_date"`
}

const (
	OperationTypeAdd    = "Add"
	OperationTypeDelete = "Delete"
)
