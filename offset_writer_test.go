package main

import(
  "testing"
  "os"
  "io/ioutil"
)

func TestPush_StartDraw(t *testing.T) {
  writer := new(OffsetWriter)
  writer.push("STARTDRAW: 1234")
  writer.end()

  if _, err := os.Stat("offsets/1234.offsets"); os.IsNotExist(err) {
    t.Error("Offset file not created")
  }
}

func TestPush_Write(t *testing.T) {
  writer := new(OffsetWriter)
  writer.start("2345")
  writer.push("123,456")
  writer.end()

  content, err := ioutil.ReadFile("offsets/2345.offsets")
  if err != nil {
    t.Error("Error reading offset file")
  }

  if string(content) != "123,456\n" {
    t.Errorf("Offset file content not valid: %s", string(content))
  }
}