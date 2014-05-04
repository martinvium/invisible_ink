package main

import (
	"github.com/gocql/gocql"
)

type Coordinate struct {
	session   *gocql.Session
	drawingId string
	timestamp gocql.UUID
	x, y      int
}

func NewCoordinate(session *gocql.Session, drawingId string, x int, y int) *Coordinate {
	return &Coordinate{session, drawingId, gocql.TimeUUID(), x, y}
}

func FindAllCoordinatesByDrawingId(session *gocql.Session, id string) *gocql.Iter {
	return session.Query(`SELECT x, y FROM coordinates WHERE drawing_id = ?`, id).Iter()
}

func DeleteAllCoordinates(session *gocql.Session) {
	_ = session.Query(`TRUNCATE coordinates`)
}

func (self *Coordinate) create() error {
	sql := `INSERT INTO coordinates (drawing_id, timestamp, x, y) VALUES (?, ?, ?, ?)`
	return self.session.Query(sql,
		self.drawingId,
		self.timestamp,
		self.x,
		self.y).Exec()
}
