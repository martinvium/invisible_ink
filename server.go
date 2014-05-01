package main

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"html/template"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("show.html"))
var validDrawingPath = regexp.MustCompile("^/drawing/([a-zA-Z0-9]+)$")

func showAction(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "show.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getDrawingId(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validDrawingPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid drawing id")
	}
	return m[1], nil // The title is the second subexpression.
}

func drawingAction(w http.ResponseWriter, r *http.Request, session *gocql.Session) {
	id, err := getDrawingId(w, r)
	if err != nil {
		return
	}

	image := NewImage(session, id)
	fullpath, err := image.create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/"+fullpath, http.StatusFound)
}

func saveListener(ws *websocket.Conn, session *gocql.Session) {
	writer := NewProtocol(session)

	for {
		var in []byte
		err := websocket.Message.Receive(ws, &in)
		if err != nil {
			fmt.Print("Socket closed")
			break
		}

		if err := writer.execute(string(in)); err != nil {
			fmt.Printf("ERROR parsing: %s\n", string(in))
		} else {
			fmt.Printf("Received: %s\n", string(in))
		}

	}

	writer.end()
}

func NewCassandraSession(keyspace string) *gocql.Session {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	return session
}

func main() {
	session := NewCassandraSession("mykeyspace")
	defer session.Close()

	http.HandleFunc("/drawing/", func(w http.ResponseWriter, r *http.Request) {
		drawingAction(w, r, session)
	})

	http.Handle("/save", websocket.Handler(func(ws *websocket.Conn) {
		saveListener(ws, session)
	}))

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/drawings/", http.StripPrefix("/drawings/", http.FileServer(http.Dir("drawings"))))
	http.HandleFunc("/", showAction)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
