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

var Data_v = map[string] int{};

type Point struct{
  Ip string
  Port string
  Id byte
  Conn *net.Conn

}

type page struct {
  Title string
  Msg string
  Data map[string] int
}

type requestPage struct {
  Id string `json:"id"`
  Value string `json:"value"`
  Name string `json:"name"`
}

func index(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-type","text/html");
  t, err := template.ParseFiles("index.html")
  if err !=nil {log.Panic(err)}
  t.Execute(w, &page{Title:"Just Page",Msg: "Just Message", Data : Data_v});
}

func main() {
  go tcpServer_P1();
  go tcpServer_P2();
  go webServer();
  var input string
  fmt.Scanln(&input)
}


func webServer(){
  http.HandleFunc("/", index)
  http.HandleFunc("/P1", P1_action)
  http.HandleFunc("/P2", P2_action)
  http.ListenAndServe(":8080", nil)
}

func P1_action (w http.ResponseWriter, r *http.Request){
  if (Point_1.Conn == nil){
      log.Printf("P1 not connect");
      return
  }
  log.Printf("P1_action");
  decoder := json.NewDecoder(r.Body)
    var request requestPage
    err := decoder.Decode(&request)
    if err != nil {
      log.Fatal(err)
    }
    log.Printf("request = ",request);
    if (request.Id == "P1_RGB_G"){
      requestValue, err := strconv.Atoi(request.Value);
      if err != nil {
        log.Fatal(err)
      }
        arduino.AnalogWrite(*Point_1.Conn,5,byte(requestValue));
    }else if (request.Id == "P1_resetPoint"){
        log.Printf("ResetPoint");
        arduino.ResetPoint(*Point_1.Conn);
    }else if (request.Id == "P1_testBut"){
      log.Printf("testButton");
      arduino.DHTGetHumi(*Point_1.Conn,3);
    }else if (request.Id == "P1_testBut2"){
      log.Printf("testButton2");
      arduino.DHTGetTemp(*Point_1.Conn,3);
    }else if (request.Id == "P1_testBut3"){
      log.Printf("testButton3");
      arduino.DHTGetTemp(*Point_1.Conn,3);
      arduino.DHTGetHumi(*Point_1.Conn,3);
      arduino.AnalogRead(*Point_1.Conn,arduino.A0);
    }
}

func P2_action (w http.ResponseWriter, r *http.Request){
  if (Point_2.Conn == nil){
      log.Printf("P2 not connect");
      return
  }
  decoder := json.NewDecoder(r.Body)
    var request requestPage
    err := decoder.Decode(&request)
    if err != nil {
      log.Fatal(err)
    }
    if (request.Id == "P2_RGB_R"){
      requestValue, err := strconv.Atoi(request.Value);
      if err != nil {
        log.Fatal(err)
      }
      arduino.AnalogWrite(*Point_2.Conn,arduino.P2_PIN_RGB_R,byte(requestValue));
    }else if (request.Id == "P2_RGB_G"){
      requestValue, err := strconv.Atoi(request.Value);
      if err != nil {
        log.Fatal(err)
      }
      arduino.AnalogWrite(*Point_2.Conn,arduino.P2_PIN_RGB_G,byte(requestValue));
    }else if (request.Id == "P2_RGB_B"){
      requestValue, err := strconv.Atoi(request.Value);
      if err != nil {
        log.Fatal(err)
      }
      arduino.AnalogWrite(*Point_2.Conn,arduino.P2_PIN_RGB_B,byte(requestValue));
    }
    if (request.Id == "P2_resetPoint"){
        log.Printf("ResetPoint 2");
        arduino.ResetPoint(*Point_2.Conn);
    }
}

var Point_1 Point;
var Point_2 Point;

func tcpServer_P1(){
  Point_1.Ip = "10.10.2.5";
  Point_1.Port =  "3005";
  Point_1.Id = 5;

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

func tcpServer_P2(){
  Point_2.Ip = "10.10.2.6";
  Point_2.Port =  "3006";
  Point_2.Id = 6;

  p2, err := net.Listen("tcp", ":"+Point_2.Port)
  if err != nil {
      log.Println("Error listening:", err.Error())
      os.Exit(1)
  }
  // Close the listener when the application closes.
  defer p2.Close()
  log.Println("Listening P_2 on "  + Point_2.Port + " port")
  for {
      // Listen for an incoming connection.
      conn, err := p2.Accept()
      Point_2.Conn = &conn;
      if err != nil {
          log.Println("P_2 -> Error accepting: ", err.Error())
          os.Exit(1)
      }else{
        log.Println("P_2 -> Connected")
      }
      // Handle connections in a new goroutine.
      go point2Server(conn)
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
        Data_v["Temp"] = ((int(buf[ARG_2])>>4) + (int(buf[ARG_3])<<4));
        log.Printf("get temp on pin %d -> %d",buf[ARG_1],((int(buf[ARG_2])>>4) + (int(buf[ARG_3])<<4)));
      }
      if (buf[COMMAND] == arduino.COMMAND_DHT_GET_HUMI){
        Data_v["Humi"] = ((int(buf[ARG_2])>>4) + (int(buf[ARG_3])<<4));
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

func point2Server(conn net.Conn) {
  for{
      buf := make([]byte, 5)
      reqLen, err := conn.Read(buf)
      if err != nil {
        // log.Println("Error reading:", err.Error())
      }else{
        log.Printf("P_2 -> Incoming data! Len = %d Data = %d",reqLen,buf);
      }

      if (buf[COMMAND] == arduino.COMMAND_SETUP){
        arduino.SetupPoint_2(conn);
      }
      if (buf[COMMAND] == arduino.COMMAND_DHT_GET_TEMP){
        Data_v["Temp"] = ((int(buf[ARG_2])>>4) + (int(buf[ARG_3])<<4));
        log.Printf("get temp on pin %d -> %d",buf[ARG_1],((int(buf[ARG_2])>>4) + (int(buf[ARG_3])<<4)));
      }
      if (buf[COMMAND] == arduino.COMMAND_DHT_GET_HUMI){
        Data_v["Humi"] = ((int(buf[ARG_2])>>4) + (int(buf[ARG_3])<<4));
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
