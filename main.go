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
    ARG_3  = 4
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
  t, err := template.ParseFiles("index.html")
  if err !=nil {log.Panic(err)}
  t.Execute(w, &page{Title:"Just Page",Msg: "Just Message"});
}

func main() {
  go tcpServer();
  go webServer();
  var input string
  fmt.Scanln(&input)
}


func webServer(){
  http.HandleFunc("/", index)
  http.HandleFunc("/api/v1/values", changeRGB)
  http.HandleFunc("/api/v1/resetPoint", resetPoint)
  http.HandleFunc("/api/v1/tBut", tBut)
  http.HandleFunc("/api/v1/tBut2", tBut2)
  http.HandleFunc("/api/v1/tBut3", tBut3)
  http.ListenAndServe(":8080", nil)
}
func resetPoint (w http.ResponseWriter, r *http.Request){
  log.Printf("ResetPoint");
  arduino.ResetPoint(*Point_1.Conn);
}

func tBut (w http.ResponseWriter, r *http.Request){
  log.Printf("testButton");
  arduino.DHTGetHumi(*Point_1.Conn,3);
}
func tBut2 (w http.ResponseWriter, r *http.Request){
  log.Printf("testButton2");
  arduino.DHTGetTemp(*Point_1.Conn,3);
}
func tBut3 (w http.ResponseWriter, r *http.Request){
  log.Printf("testButton3");
  arduino.DHTGetTemp(*Point_1.Conn,3);
  arduino.DHTGetHumi(*Point_1.Conn,3);
  arduino.AnalogRead(*Point_1.Conn,arduino.A0);
}

func changeRGB (w http.ResponseWriter, r *http.Request){
  decoder := json.NewDecoder(r.Body)
    var color rgbResponse
    err := decoder.Decode(&color)
    if err != nil {
      log.Fatal(err)
    }
    colorValue, err := strconv.Atoi(color.Value);
    if (color.Id == "Red"){
      arduino.AnalogWrite(*Point_1.Conn,5,byte(colorValue));
    }
}
var Point_1 Point;

func tcpServer(){
  Point_1.Ip = "10.10.2.9";
  Point_1.Port =  "3333";
  Point_1.Id = 9;

  p1, err := net.Listen("tcp", ":"+Point_1.Port)
  if err != nil {
      log.Println("Error listening:", err.Error())
      os.Exit(1)
  }
  // Close the listener when the application closes.
  defer p1.Close()
  log.Println("Listening P_1 on "  + Point_1.Port + " port")
  for {
      // Listen for an incoming connection.
      conn, err := p1.Accept()
      Point_1.Conn = &conn;
      if err != nil {
          log.Println("P_1 -> Error accepting: ", err.Error())
          os.Exit(1)
      }else{
        log.Println("P_1 -> Connected")
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
        log.Printf("P_1 -> Incoming data! Len = %d Data = %d",reqLen,buf);
      }

      if (buf[COMMAND] == arduino.COMMAND_SETUP){
        arduino.SetupPoint_1(conn);
      }
      if (buf[COMMAND] == arduino.COMMAND_DHT_GET_TEMP){
        log.Printf("get temp on pin %d -> %d",buf[ARG_1],((int(buf[ARG_2])>>4) + (int(buf[ARG_3])<<4)));
      }
      if (buf[COMMAND] == arduino.COMMAND_DHT_GET_HUMI){
        log.Printf("get humi on pin %d -> %d",buf[ARG_1],((int(buf[ARG_2])>>4) + (int(buf[ARG_3])<<4)));
      }
      if (buf[COMMAND] == arduino.COMMAND_DIGITAL_READ){
        if (buf[ARG_1] == 2){
            arduino.DigitalWrite(*Point_1.Conn,6,buf[ARG_2])
        }
      }
      if (buf[COMMAND] == arduino.COMMAND_ANALOG_READ){
        log.Printf("get value on pin %d -> %d",buf[ARG_1],((int(buf[ARG_2])>>4) + (int(buf[ARG_3])<<4)));
      }
    }
}
