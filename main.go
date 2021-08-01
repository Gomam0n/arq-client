package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

var IP = "0.0.0.0"

// change this IP
var serverIP = "10.10.80.2"
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
	var newConn net.Conn
	for {
		data := make([]byte, 17)
		len, addr, err := readConn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if addr.IP.String() != serverIP {
			fmt.Println("IP is not server IP")
			continue
		}
		newConn, err = DialClient(addr, string(data[:len]))
		fmt.Println("server address: ", newConn.RemoteAddr())
		break
	}

	sequenceNumber := 0
	sendData := "ACK" + strconv.Itoa(sequenceNumber)
	newConn.Write([]byte(sendData))
	fmt.Println(sendData)
	fmt.Println(readConn.LocalAddr())
	sequenceNumber = 1
	var content []byte
	for {
		data := make([]byte, 17)
		len, addr, err := readConn.ReadFromUDP(data)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if addr.IP.String() != serverIP {
			continue
		}
		if data[0] != byte(sequenceNumber+'0') {
			fmt.Println("sequence number not correct, ignore the packet and resend ACK")
			newConn.Write([]byte(sendData))
			continue
		}
		fmt.Println("receive data ", string(data[1:len]))
		sendData = "ACK" + strconv.Itoa(sequenceNumber)
		newConn.Write([]byte(sendData))
		fmt.Println(sendData)
		sequenceNumber = 1 - sequenceNumber

		content = append(content, data[1:len]...)
		if len < 17 {
			break
		}
	}
	WriteFile("receive.txt", content)
	fmt.Println("file stored")
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
func DialClient(addr *net.UDPAddr, data string) (conn net.Conn, err error) {
	port, err := strconv.ParseUint(data, 10, 17)
	if err != nil {
		fmt.Println("Parse unsigned int error: " + err.Error())
		return
	}
	if port < 1024 {
		fmt.Println(strconv.Itoa(addr.Port) + ": the port should be 1024~65535")
		return
	}
	clientAddr := addr.IP.String() + ":" + strconv.Itoa(int(port))
	// error?
	conn, err = net.Dial("udp", clientAddr)
	if err != nil {
		fmt.Println("The connection with " + addr.IP.String() + " has error: " + err.Error())
	}
	return
}
