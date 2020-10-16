package main

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
)

func main() {
	ctlConn, err := net.Dial("tcp", ":4422")
	if err != nil {
		panic(err)
	}
	defer ctlConn.Close()

	buffer, _ := json.Marshal(struct {
		Name string
		Age  int32
	}{
		Name: "张三",
		Age:  12,
	})

	err = binary.Write(ctlConn,binary.LittleEndian,int64(len(buffer)))
	if err != nil{
		panic(err)
	}

	if _,err = ctlConn.Write(buffer);err != nil{
		panic(err)
	}
	log.Println("请求结束")
}
