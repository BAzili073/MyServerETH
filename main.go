package main
//GOOS=linux GOARCH=arm GOARM=6 go build
//scp ~/Work/goHomeServer/* pi@192.168.1.50:GoServer

import (
    "fmt"
    "net/http"
    "net"
    "html/template"
    "log"
    // "time"
    "os"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)

type page struct {
  Title string
  Msg string
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love  %s!", r.URL.Path[1:])
}

func index(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-type","text/html");
  t, err := template.ParseFiles("./index.html")
  if err !=nil {log.Panic(err)}
  t.Execute(w, &page{Title:"Just Page",Msg: "Just Message"});
  log.Printf("hello");
}

func main() {
  // Listen for incoming connections.
 l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
 if err != nil {
     fmt.Println("Error listening:", err.Error())
     os.Exit(1)
 }
 // Close the listener when the application closes.
 defer l.Close()
 fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
 for {
     // Listen for an incoming connection.
     conn, err := l.Accept()
     if err != nil {
         fmt.Println("Error accepting: ", err.Error())
         os.Exit(1)
     }
     // Handle connections in a new goroutine.
     go handleRequest(conn)
 }

  http.HandleFunc("/", index)
  http.ListenAndServe(":8080", nil)

}

func handleRequest(conn net.Conn) {
  // Make a buffer to hold incoming data.
  buf := make([]byte, 1024)
  // Read the incoming connection into the buffer.
  reqLen, err := conn.Read(buf)
  if err != nil {
    log.Println("Error reading:", err.Error())
  }else{
    log.Printf("Len = %d",reqLen);
    log.Printf("Buff = %s",buf);
  }
  // Send a response back to person contacting us.
  conn.Write([]byte("Message received."))
  // Close the connection when you're done with it.
  // conn.Close()
}
