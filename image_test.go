package main

import (
	"github.com/gocql/gocql"
	"os"
	"testing"
)

func TestImageCreate(t *testing.T) {
	session := NewCassandraSession("testkeyspace")
	defer session.Close()
	CreateImageTestData(session, "exampledrawingid")

	image := NewImage(session, "exampledrawingid")
	fullpath, err := image.create()

	AssertNotError(t, err)
	AssertNotEqual(t, "", fullpath)
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		t.Error("Image file not found")
	}

	os.Remove(fullpath)
	DeleteAllCoordinates(session)
}

func CreateImageTestData(session *gocql.Session, id string) {
	protocol := NewProtocol(session)
	protocol.execute("STARTDRAW: " + id)
	protocol.execute("123,234")
	protocol.execute("432,524")
	protocol.execute("ENDDRAW")
}
