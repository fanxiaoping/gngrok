package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"
	vhost "github.com/inconshreveable/go-vhost"
)

const (
	connReadTimeout time.Duration = 10 * time.Second
	BadRequest = `HTTP/1.0 400 Bad Request
Content-Length: 12

Bad Request
`
	NotFound = `HTTP/1.0 404 Not Found
Content-Length: %d

Tunnel %s not found
`
)

type Message interface{}

type Envelope struct {
	Type    string
	Payload json.RawMessage
}

func main() {
	listener, err := net.Listen("tcp", ":4422")
	if err != nil {
		panic(err)
	}

	go func() {
		hl,err := net.Listen("tcp",":3321")
		if err != nil{
			return
		}
		for {
			rawConn,err := hl.Accept()
			if err != nil {
				log.Println("Failed to accept new TCP connection:", err)
			}
			go httpHandler(rawConn)
		}
	}()

	for {
		rawConn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept new TCP connection:", err)
		}
		go tcpHandler(rawConn)
	}
}

func httpHandler(conn net.Conn){
	data,err := ioutil.ReadAll(conn)
	if err != nil{
		log.Println(err)
		return
	}
	log.Println(string(data))
	vhostConn, err := vhost.HTTP(conn)
	if err != nil {
		log.Println("Failed to read valid:", err)
		conn.Write([]byte(BadRequest))
		return
	}
	host := strings.ToLower(vhostConn.Host())
	conn.Write([]byte(fmt.Sprintf(NotFound, len(host)+18, host)))
}

func tcpHandler(conn net.Conn) {
	//conn.SetReadDeadline(time.Now().Add(connReadTimeout))
	//buffer, err := readMsgShared(conn)
	//if err != nil {
	//	log.Println("Failed to read message:", err)
	//	conn.Close()
	//	return
	//}
	////conn.SetReadDeadline(time.Time{})
	//
	//log.Println(string(buffer))

	dconn,err := net.Dial("tcp","192.168.1.220:3389")
	if err != nil {
		log.Println("连接失败")
		return
	}
	//exitChan := make(chan bool,1)
	//发送数据
	go func() {
		_,err = io.Copy(dconn,conn)
		if err != nil{
			log.Println("send",err)
			return
		}
	}()
	//读取数据
	_,err = io.Copy(conn,dconn)
	if err != nil{
		log.Println("read",err)
		return
	}
	dconn.Close()
}

func readMsgShared(c net.Conn) (buffer []byte, err error) {
	log.Println("Waiting to read message")

	var sz int64
	err = binary.Read(c, binary.LittleEndian, &sz)
	if err != nil {
		return
	}
	log.Println("Reading message with length:", sz)

	buffer = make([]byte, sz)
	n, err := c.Read(buffer)
	log.Println("Read message ", buffer)

	if err != nil {
		return
	}

	if int64(n) != sz {
		err = errors.New(fmt.Sprintf("Expected to read %d bytes, but only read %d", sz, n))
		return
	}

	return
}
