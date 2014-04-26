package main

import (
  "fmt"
  "html/template"
  // "io/ioutil"
  "net/http"
  "code.google.com/p/go.net/websocket"
  "time"
  "os"
  "bufio"
)

var templates = template.Must(template.ParseFiles("show.html"))

func showAction(w http.ResponseWriter, r *http.Request) {
    err := templates.ExecuteTemplate(w, "show.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func drawingFilename() string {
    now := time.Now().UTC()
    return "drawings/" + now.Format("20060102150405") + ".offsets"
}

func saveListener(ws *websocket.Conn) {
    file, err := os.Create(drawingFilename())
    check(err)
    buffer := bufio.NewWriter(file)

    for {
        var in []byte
        err = websocket.Message.Receive(ws, &in)
        if err != nil {
            break
        }

        msg := string(in)
        if msg == "ENDDRAW" {
            buffer.Flush()
            file.Close()
            file, err = os.Create(drawingFilename())
            check(err)
            buffer = bufio.NewWriter(file)
        } else {
            buffer.WriteString(msg + "\n")
        }

        fmt.Printf("Received: %s\n", string(in))
    }

    buffer.Flush()
    file.Close()
}

func main() {
    http.HandleFunc("/", showAction)
    http.Handle("/save", websocket.Handler(saveListener))
    http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
    
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}