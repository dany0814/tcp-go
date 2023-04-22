package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/dany0814/tcp-go/pkg/utils"
)

const BUFFER_SIZE = 1024

func main() {
	serverAddr := "localhost:11112"
	// connect to server
	conn, err := net.Dial("tcp", serverAddr)
	utils.CheckError(err)
	defer conn.Close()

	//read from
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter 'send <filename>' to transfer files to the server\n\n")
	inputFromUser, _ := reader.ReadString('\n')
	arrayOfCommands := strings.Split(inputFromUser, " ")

	if arrayOfCommands[0] == "send" {
		sendFileToServer(arrayOfCommands[1], conn)
	} else {
		fmt.Println("Bad Command")
	}
}

func sendFileToServer(fileName string, conn net.Conn) {

	var currentByte int64 = 0
	fmt.Println("send to client")
	fileBuffer := make([]byte, BUFFER_SIZE)

	var err error

	//file to read
	file, err := os.Open(strings.TrimSpace(fileName)) // For read access.
	if err != nil {
		conn.Write([]byte("-1"))
		log.Fatal(err)
	}
	conn.Write([]byte("send " + fileName))
	//read file until there is an error
	for err == nil || err != io.EOF {

		_, err = file.ReadAt(fileBuffer, currentByte)
		currentByte += BUFFER_SIZE
		fmt.Println(fileBuffer)
		conn.Write(fileBuffer)
	}

	file.Close()
	conn.Close()

}
