package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const BUFFER_SIZE = 1024

func main() {
	// Listen for incoming connections.
	addr := "0.0.0.0:11112"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	host, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening TCP Server on host: %s, port: %s\n", host, port)

	for {
		// Listen for an incoming connection
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		// Handle connections in a new goroutine
		go func(conn net.Conn) {
			buf := make([]byte, 1024)
			len, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("Error reading: %#v\n", err)
				return
			}
			fmt.Printf("Message received: %s\n", string(buf[:len]))

			cleanedBuffer := bytes.Trim(buf, "\x00")
			cleanedInputCommandString := strings.TrimSpace(string(cleanedBuffer))
			arrayOfCommands := strings.Split(cleanedInputCommandString, " ")

			fmt.Println(arrayOfCommands[0])
			if arrayOfCommands[0] == "send" {
				fmt.Println("getting a file")

				getFileFromClient(arrayOfCommands[1], conn)

			} else {
				_, err = conn.Write([]byte("bad command"))
			}
			conn.Write([]byte("Message received.\n"))
			conn.Close()
		}(conn)
	}
}

func getFileFromClient(fileName string, connection net.Conn) {

	var currentByte int64 = 0

	fileBuffer := make([]byte, BUFFER_SIZE)

	var err error
	file, err := os.Create("./test/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	connection.Write([]byte("get " + fileName))

	for err == nil || err != io.EOF {

		connection.Read(fileBuffer)

		cleanedFileBuffer := bytes.Trim(fileBuffer, "\x00")

		_, err = file.WriteAt(cleanedFileBuffer, currentByte)
		if len(string(fileBuffer)) != len(string(cleanedFileBuffer)) {
			break
		}
		currentByte += BUFFER_SIZE

		// basepath := path.Dir(file.FileName)
		// fileName := path.Base(file.FileName)
		// err = os.MkdirAll(basepath, 0777)
		// checkError(err)
		// filePath := path.Join(basepath, fileName)

	}

	connection.Close()
	file.Close()
	return

}
