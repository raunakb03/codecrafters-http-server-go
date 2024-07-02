package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	fmt.Println(line)
	params := strings.Split(strings.Split(line, " ")[1], "/")[2]
	res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(params), params)
    fmt.Println(res)
	_, err = conn.Write([]byte(res))
	if err != nil {
		panic(err)
	}
	// if reqUrl == "" {
	// 	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// } else {
	// 	_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
}
