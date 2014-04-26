package main

import (
  "fmt"
  "os"
  "bufio"
  "regexp"
)

var startRegex = regexp.MustCompile("^STARTDRAW: ([a-zA-Z0-9]+)$")
var offsetRegex = regexp.MustCompile("^[0-9]+,[0-9]+$")

type OffsetWriter struct {
  file *os.File
  buffer *bufio.Writer
}

func (self *OffsetWriter) start(id string) {
  file, err := os.Create("offsets/" + id + ".offsets")
  if err != nil {
      panic(err)
  }

  self.buffer = bufio.NewWriter(file)
  self.file = file
}

func (self *OffsetWriter) end() {

  if self.buffer != nil {
    self.buffer.Flush()
  }

  if self.file != nil {
    self.file.Close()
  }
}

func (self *OffsetWriter) write(msg string) {
  if self.buffer != nil {
    self.buffer.WriteString(msg + "\n")
  } else {
    fmt.Print("Failed to write msg - drawing not started")
  }
}

func (self *OffsetWriter) push(msg string) {
  if m := offsetRegex.FindStringSubmatch(msg); m != nil {
    self.write(msg)
  } else if m := startRegex.FindStringSubmatch(msg); m != nil {
    self.start(m[1])
  } else if msg == "ENDDRAW" {
    self.end()
  } else {
    fmt.Printf("Invalid command pushed: %s\n", msg)
  }
}