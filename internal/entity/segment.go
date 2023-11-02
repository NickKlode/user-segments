package entity

type Segment struct {
	Id      int    `db:"id"`
	Name    string `db:"name"`
	Percent int    `db:"percent"`
}
