package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

var IP = "127.0.0.1"
var serverIP = "127.0.0.1"
var serverPort = 8888
var clientPort = 9999

func main() {

	serverAddr := serverIP + ":" + strconv.Itoa(serverPort)
	conn, err := net.Dial("udp", serverAddr)
	checkError(err)

	defer conn.Close()

	udpAddr, err := getUDPAddr(clientPort)
	checkError(err)

	readConn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	defer readConn.Close()

	_, err = conn.Write([]byte(strconv.Itoa(clientPort)))
	checkError(err)

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
	WriteFile("receive.txt", content)

}
func WriteFile(name string, content []byte) {
	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(file, string(content))
}
func getUDPAddr(port int) (udpAddr *net.UDPAddr, err error) {
	serverAddress := IP + ":" + strconv.Itoa(port)
	udpAddr, err = net.ResolveUDPAddr("udp", serverAddress)
	return
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
