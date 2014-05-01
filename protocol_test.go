package main

import (
	"testing"
)

func TestProtocolExecuteStartDraw(t *testing.T) {
	session := NewCassandraSession("testkeyspace")
	defer session.Close()
	protocol := NewProtocol(session)

	protocol.execute("STARTDRAW: exampledrawingid")

	AssertEqual(t, "exampledrawingid", protocol.drawing_id)
}

func TestProtocolExecuteEndDrawClearsId(t *testing.T) {
	session := NewCassandraSession("testkeyspace")
	defer session.Close()
	protocol := NewProtocol(session)

	protocol.execute("STARTDRAW: exampledrawingid")
	protocol.execute("ENDDRAW")

	AssertEqual(t, "", protocol.drawing_id)
}

func TestProtocolExecuteStartAndCoordinates(t *testing.T) {
	session := NewCassandraSession("testkeyspace")
	defer session.Close()
	protocol := NewProtocol(session)

	protocol.execute("STARTDRAW: 1234")
	protocol.execute("123,543")
	protocol.execute("ENDDRAW")

	var x, y int
	iter := FindAllCoordinatesByDrawingId(session, "1234")
	iter.Scan(&x, &y)
	AssertEqual(t, x, 123)
	AssertEqual(t, y, 543)

	DeleteAllCoordinates(session)
}
