package main

import (
  "fmt"
  "regexp"
  "github.com/gocql/gocql"
  "strconv"
)

var startRegex = regexp.MustCompile("^STARTDRAW: ([a-zA-Z0-9]+)$")
var offsetRegex = regexp.MustCompile("^([0-9]+),([0-9]+)$")

type OffsetWriter struct {
  session *gocql.Session
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
  uuid := gocql.TimeUUID()
  ix, _ := strconv.ParseInt(x, 10, 32)
  iy, _ := strconv.ParseInt(y, 10, 32)

  fmt.Printf("Args: %s, %s, %d, %d\n", uuid, self.drawing_id, int(ix), int(iy))

  if err := self.session.Query(`INSERT INTO coordinates (id, drawing_id, x, y) VALUES (?, ?, ?, ?)`,
      uuid, self.drawing_id, int(ix), int(iy)).Exec(); err != nil {
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