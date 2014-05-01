package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"regexp"
	"strconv"
)

var startRegex = regexp.MustCompile("^STARTDRAW: ([a-zA-Z0-9]+)$")
var offsetRegex = regexp.MustCompile("^([0-9]+),([0-9]+)$")

type OffsetWriter struct {
	session    *gocql.Session
	drawing_id string
}

func NewOffsetWriter(session *gocql.Session) *OffsetWriter {
	f := OffsetWriter{session, ""}
	return &f
}

func (self *OffsetWriter) start(id string) {
	self.drawing_id = id
}

func (self *OffsetWriter) end() {
	self.drawing_id = ""
}

func (self *OffsetWriter) write(x string, y string) {
	ix, _ := strconv.ParseInt(x, 10, 32)
	iy, _ := strconv.ParseInt(y, 10, 32)

	coord := NewCoordinate(self.session, self.drawing_id, int(ix), int(iy))
	if err := coord.create(); err != nil {
		panic(err)
	}
}

func (self *OffsetWriter) push(msg string) {
	if m := offsetRegex.FindStringSubmatch(msg); m != nil {
		self.write(m[1], m[2])
	} else if m := startRegex.FindStringSubmatch(msg); m != nil {
		self.start(m[1])
	} else if msg == "ENDDRAW" {
		self.end()
	} else {
		fmt.Printf("Invalid command pushed: %s\n", msg)
	}
}
