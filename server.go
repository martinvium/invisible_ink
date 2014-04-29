package main

import (
  "fmt"
  "html/template"
  "net/http"
  "code.google.com/p/go.net/websocket"
  "image" 
  "image/color" 
  "errors"
  "os"
  "regexp"
  "image/png"
  "github.com/gocql/gocql"
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

    canvas := image.NewRGBA(image.Rect(0, 0, 300, 300))

    var x, y int;
    iter := session.Query(`SELECT x, y FROM coordinates WHERE drawing_id = ?`, id).Iter()
    for iter.Scan(&x, &y) {
        canvas.Set(x, y, color.Black)
    }

    if err := iter.Close(); err != nil {
        panic(err)
    }

    file, err := os.Create("drawings/" + id + ".png")
    if err != nil {
        panic(err)
    }

    defer file.Close()

    png.Encode(file, canvas)

    http.Redirect(w, r, "/drawings/"+id+".png", http.StatusFound)
}

func saveListener(ws *websocket.Conn, session *gocql.Session) {
    writer := NewOffsetWriter(session)

    for {
        var in []byte
        err := websocket.Message.Receive(ws, &in)
        if err != nil {
            fmt.Print("Socket closed")
            break
        }

        writer.push(string(in))
        fmt.Printf("Received: %s\n", string(in))
    }

    writer.end()
}

func getCassandraSession() *gocql.Session {
    cluster := gocql.NewCluster("127.0.0.1")
    cluster.Keyspace = "mykeyspace"
    cluster.Consistency = gocql.Quorum
    session, err := cluster.CreateSession()
    if err != nil {
        panic(err)
    }

    return session
}

func main() {
    session := getCassandraSession()
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