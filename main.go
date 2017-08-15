package main
//GOOS=linux GOARCH=arm GOARM=6 go build
//scp ~/Work/goHomeServer/* pi@192.168.1.50:GoServer

import (
    "fmt"
    "net/http"
    "net"
    "html/template"
    "log"
    "encoding/json"
    // "time"
    "os"
    "./arduino"
    "strconv"
)

const (
    ID  = 0
    COMMAND = 1
    ARG_1 = 2
    ARG_2  = 3
)

type Point struct{
  Ip string
  Port string
  Id byte
  Conn *net.Conn

}

type page struct {
  Title string
  Msg string
}

type rgbResponse struct {
  Value string `json:"value"`
  Id string `json:"id"`
}

func index(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-type","text/html");
  t, err := template.ParseFiles("./index.html")
  if err !=nil {log.Panic(err)}
  t.Execute(w, &page{Title:"Just Page",Msg: "Just Message"});
}

func main() {
  go webServer();
  go tcpServer();

  var input string
  fmt.Scanln(&input)
}

func webServer(){
  for {
    http.HandleFunc("/", index)
    http.HandleFunc("/api/v1/values", changeRGB)
    http.ListenAndServe(":8080", nil)
  }
}

func changeRGB (w http.ResponseWriter, r *http.Request){
  fmt.Println("cnhage")
  decoder := json.NewDecoder(r.Body)
    var color rgbResponse
    err := decoder.Decode(&color)
    if err != nil {
      log.Fatal(err)
    }
    colorValue, err := strconv.Atoi(color.Value);
    if (color.Id == "Red"){
      // arduino.DigitalWrite(Point_1.Conn,1,1)
      fmt.Println("Red = ",colorValue)
    }
}

func tcpServer(){
  Point_1 := Point{
    Ip : "10.10.2.9",
    Port : "3333",
    Id : 9};
  p1, err := net.Listen("tcp", ":"+Point_1.Port)
  if err != nil {
      fmt.Println("Error listening:", err.Error())
      os.Exit(1)
  }
  // Close the listener when the application closes.
  defer p1.Close()
  fmt.Println("Listening on "  + Point_1.Port + " port")
  for {
      // Listen for an incoming connection.
      conn, err := p1.Accept()
      Point_1.Conn = &conn;
      if err != nil {
          fmt.Println("Error accepting: ", err.Error())
          os.Exit(1)
      }else{
        fmt.Println("Connection accept")
      }
      // Handle connections in a new goroutine.
      go point1Server(conn)
     //  l.Close();
  }
}

func point1Server(conn net.Conn) {
  for{
      buf := make([]byte, 5)
      reqLen, err := conn.Read(buf)
      if err != nil {
        // log.Println("Error reading:", err.Error())
      }else{
        log.Printf("Len = %d",reqLen);
        log.Printf("Buff = %d",buf);
      }

      if (buf[COMMAND] == arduino.COMMAND_SETUP){
        arduino.SetupPoint_1(conn);
      }
    }
}

// func pinMode(conn net.Conn, pin byte, mode byte){
//   Send_message := ;
//
// }
