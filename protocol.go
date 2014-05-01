package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"regexp"
	"strconv"
)

var startRegex = regexp.MustCompile("^STARTDRAW: ([a-zA-Z0-9]+)$")
var offsetRegex = regexp.MustCompile("^([0-9]+),([0-9]+)$")

type Protocol struct {
	session    *gocql.Session
	drawing_id string
}

func NewProtocol(session *gocql.Session) *Protocol {
	f := Protocol{session, ""}
	return &f
}

func (self *Protocol) start(id string) {
	self.drawing_id = id
}

func (self *Protocol) end() {
	self.drawing_id = ""
}

func (self *Protocol) write(x string, y string) error {
	ix, _ := strconv.ParseInt(x, 10, 32)
	iy, _ := strconv.ParseInt(y, 10, 32)

	coord := NewCoordinate(self.session, self.drawing_id, int(ix), int(iy))
	return coord.create()
}

func (self *Protocol) execute(msg string) error {
	if m := offsetRegex.FindStringSubmatch(msg); m != nil {
		return self.write(m[1], m[2])
	} else if m := startRegex.FindStringSubmatch(msg); m != nil {
		self.start(m[1])
	} else if msg == "ENDDRAW" {
		self.end()
	} else {
		fmt.Printf("Invalid command pushed: %s\n", msg)
	}

	return nil
}
