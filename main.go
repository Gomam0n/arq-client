package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

var serverIP = "127.0.0.1"
var serverPort = 8888
var clientPort = 9999

func main() {

	serverAddr := serverIP + ":" + strconv.Itoa(serverPort)
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	clientAddress := "0.0.0.0:" + strconv.Itoa(clientPort)
	udpAddr, err := net.ResolveUDPAddr("udp", clientAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	readConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer readConn.Close()

	_, err = conn.Write([]byte(strconv.Itoa(clientPort)))
	if err != nil {
		fmt.Println(err)
		return
	}

	var content []byte
	for {
		data := make([]byte, 16)
		len, addr, err := readConn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if addr.IP.String() != serverIP {
			continue
		}
		content = append(content, data[:len]...)
		if len < 16 {
			break
		}
	}
	file, err := os.Create("receive.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(file, string(content))

}
