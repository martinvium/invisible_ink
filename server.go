package main

import (
  "fmt"
  "html/template"
  "net/http"
  "code.google.com/p/go.net/websocket"
  "time"
)

var templates = template.Must(template.ParseFiles("show.html"))

func showAction(w http.ResponseWriter, r *http.Request) {
    err := templates.ExecuteTemplate(w, "show.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func drawingAction(w http.ResponseWriter, r *http.Request) {
    
}

func drawingFilename() string {
    now := time.Now().UTC()
    return "offsets/" + now.Format("20060102150405") + ".offsets"
}

func saveListener(ws *websocket.Conn) {
    writer := new(OffsetWriter)

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

func main() {
    http.HandleFunc("/", showAction)
    http.HandleFunc("/drawing", drawingAction)
    http.Handle("/save", websocket.Handler(saveListener))
    http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
    http.Handle("/drawings/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
    
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}