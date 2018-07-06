package main

import (
    "fmt"
    "log"
	//"bytes"
    "strings"
    "html/template"
    "net/http"
	"io/ioutil"
	//"encoding/xml"

    "github.com/gorilla/websocket"
)

type msg struct {
	Num int
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
    render(w, "welcome.html")
}

func boseSoundtouch(w http.ResponseWriter, req *http.Request) {
    fmt.Println("SoundTouch")
    render(w, "bosesoundtouch.html")
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	go echo(conn)
}

func echo(conn *websocket.Conn) {

    boseurl := "http://192.168.1.185:8090/key"

    pressxml := "<?xml version='1.0' ?><key state=\"press\" sender=\"Gabbo\">" 
    releasexml := "<?xml version='1.0' ?><key state=\"release\" sender=\"Gabbo\">"

	for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }

        var data = string(p[:])

        /* # ---------------------------------- 
           #   Process bose connection here     
           # ---------------------------------- */
        if strings.Contains(data, "bose") {
            fmt.Println(boseurl + " " + data) 

            parts := strings.Split(data, ",")
            pressxml += parts[1] + "</key>"
            releasexml += parts[1] + "</key>"
        
            fmt.Println("Press: " + pressxml)
            fmt.Println("Release: " + releasexml)
            //buf, _ := xml.Marshal(pressxml)
            //body := bytes.NewBuffer(buf)
    
            rpxml, _ := http.Post(boseurl, "text/xml", strings.NewReader(pressxml))
            resp, _ := ioutil.ReadAll(rpxml.Body)
            print(string(resp))

            rlxml, _ := http.Post(boseurl, "text/xml", strings.NewReader(releasexml))
            resprlxml, _ := ioutil.ReadAll(rlxml.Body)
            print(string(resprlxml))
        }

        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }
	}
}

func render(w http.ResponseWriter, tmpl string) {
    tmpl = fmt.Sprintf("templates/%s", tmpl)
    t, err := template.ParseFiles(tmpl)

    if err != nil {
        log.Print("template parsing error: ", err)
    }
    err = t.Execute(w, "test")
    if err != nil {
        log.Print("template executing error: ", err)
    }
}

func main() {
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/", helloWorld)
    http.HandleFunc("/bosesoundtouch", boseSoundtouch)
    http.HandleFunc("/websocket", websocketHandler)

    err := http.ListenAndServe(":80", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
