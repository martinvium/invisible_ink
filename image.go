package main

import (
	"github.com/gocql/gocql"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Image struct {
	session *gocql.Session
	id      string
}

func NewImage(session *gocql.Session, id string) *Image {
	return &Image{session, id}
}

func (self *Image) create() (string, error) {
	canvas := image.NewRGBA(image.Rect(0, 0, 300, 300))

	var x, y int
	iter := FindAllCoordinatesByDrawingId(self.session, self.id)
	for iter.Scan(&x, &y) {
		canvas.Set(x, y, color.Black)
	}

	if err := iter.Close(); err != nil {
		return "", err
	}

	fullpath := "drawings/" + self.id + ".png"
	file, err := os.Create(fullpath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	png.Encode(file, canvas)

	return fullpath, nil
}
