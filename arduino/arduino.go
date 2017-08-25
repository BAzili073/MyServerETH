package arduino

import (
  "net"
  "log"
)
const(
  PIN_INPUT = 0
  PIN_OUTPUT = 1

  COMMAND_SETUP = 1
  COMMAND_PINMODE = 2
  COMMAND_DIGITAL_WRITE = 3
  COMMAND_ANALOG_WRITE = 4
  COMMAND_PIN_NOTIFY_SETTING = 5
  COMMAND_PIN_SET_SETTING = 6
  COMMAND_PIN_RESET_SETTING = 7
  COMMAND_DIGITAL_READ = 8
  COMMAND_ANALOG_READ = 9
  COMMAND_DHT_GET_TEMP = 10
  COMMAND_DHT_TEMP_RESP = 11
  COMMAND_DHT_GET_HUMI = 12
  COMMAND_DHT_HUMI_RESP = 13
  COMMAND_DHT_ADD = 14
  COMMAND_FULL_RESET = 15
  COMMAND_CHANGE_POINT_ID = 16
  COMMAND_POINT_ON = 17
  COMMAND_PING = 18
  COMMAND_PONG = 19
  COMMAND_PING_FOR_ALL = 20
  COMMAND_CHANGE_MASTER_ID = 21

  DHT11 = 11
  DHT21 = 21
  DHT22 = 22

  A0 = 14
  A1 = 15
  A2 = 16
  A3 = 17
  A4 = 18
  A5 = 19
  A6 = 20
  A7 = 21


  P2_PIN_RGB_R = 5
  P2_PIN_RGB_G = 3
  P2_PIN_RGB_B = 6
)

func sendMessage(conn net.Conn,command byte,arg_1 byte,arg_2 byte,arg_3 byte){
    conn.Write([]byte{1,command,arg_1,arg_2,arg_3});
}

func PinMode(conn net.Conn, pin byte, mode byte){
  sendMessage(conn,COMMAND_PINMODE,pin,mode,0);
}

func DigitalWrite(conn net.Conn, pin byte, state byte){
  sendMessage(conn,COMMAND_DIGITAL_WRITE,pin,state,0);
}

func DigitalRead(conn net.Conn, pin byte){
  sendMessage(conn,COMMAND_DIGITAL_READ,pin,0,0);
}

func AnalogWrite(conn net.Conn, pin byte, state byte){
  sendMessage(conn,COMMAND_ANALOG_WRITE,pin,state,0);
}

func AnalogRead(conn net.Conn, pin byte){
  sendMessage(conn,COMMAND_ANALOG_READ,pin,0,0);
}

func ResetPoint(conn net.Conn){
  sendMessage(conn,COMMAND_FULL_RESET,0,0,0);
}

func PinNotify(conn net.Conn, pin byte, state byte){
  sendMessage(conn,COMMAND_PIN_NOTIFY_SETTING,pin,state,0);
}

func DHTAdd(conn net.Conn, pin byte, DHT_type byte){
  sendMessage(conn,COMMAND_DHT_ADD,pin,DHT_type,0);
}

func DHTGetTemp(conn net.Conn, pin byte){
  sendMessage(conn,COMMAND_DHT_GET_TEMP,pin,0,0);
}

func DHTGetHumi(conn net.Conn, pin byte){
  sendMessage(conn,COMMAND_DHT_GET_HUMI,pin,0,0);
}

func SetupPoint_1(conn net.Conn){
  log.Printf("Setup point 1");
  PinMode(conn,A0,PIN_INPUT);
  PinMode(conn,7,PIN_INPUT);
  PinNotify(conn,7,1);
  PinMode(conn,2,PIN_INPUT);
  PinNotify(conn,2,1);
  PinMode(conn,6,PIN_OUTPUT);
  PinMode(conn,5,PIN_OUTPUT);
  PinMode(conn,4,PIN_OUTPUT);
  DHTAdd(conn,3,DHT11);
  sendMessage(conn,COMMAND_SETUP,1,1,0);
}

func SetupPoint_2(conn net.Conn){
  log.Printf("Setup point 2");
  PinMode(conn,P2_PIN_RGB_R,PIN_OUTPUT);
  PinMode(conn,P2_PIN_RGB_G,PIN_OUTPUT);
  PinMode(conn,P2_PIN_RGB_B,PIN_OUTPUT);
  sendMessage(conn,COMMAND_SETUP,1,1,0);
}
