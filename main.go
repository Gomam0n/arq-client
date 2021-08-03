package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

// IP local IP address
var IP = "0.0.0.0"

// change this IP according to the server
var serverIP = "10.10.80.2"
var serverPort = 8888
var clientPort = 9999

func main() {
	// dial the server
	serverAddr := serverIP + ":" + strconv.Itoa(serverPort)
	conn, err := net.Dial("udp", serverAddr)
	checkError(err)
	defer conn.Close()

	udpAddr, err := getUDPAddr(clientPort)
	checkError(err)

	// listen to @clientPort
	readConn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	defer readConn.Close()

	// tell server the @clientPort
	_, err = conn.Write([]byte(strconv.Itoa(clientPort)))
	checkError(err)
	var newConn net.Conn
	// wait for server to respond
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
		// dial server according to the port
		newConn, err = DialServer(addr, string(data[:len]))
		fmt.Println("server address: ", newConn.RemoteAddr())
		break
	}

	sequenceNumber := 0
	sendData := "ACK" + strconv.Itoa(sequenceNumber)
	// the first ack
	newConn.Write([]byte(sendData))
	fmt.Println(sendData)
	sequenceNumber = 1
	var content []byte
	for {
		data := make([]byte, 17)
		len, addr, err := readConn.ReadFromUDP(data)

		// the data should be from server IP and with correct sequence number
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
		fmt.Println("receive data: ", string(data[1:len]))

		// send ACK
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

// WriteFile write content to file according to the filename
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

func DialServer(addr *net.UDPAddr, data string) (conn net.Conn, err error) {
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
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
