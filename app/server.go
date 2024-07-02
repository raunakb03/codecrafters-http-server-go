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
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	// split the line by spaces
	reqUrl := strings.Split(line, " ")[1][1:]
	splitLine := strings.Split(line, " ")[1]
	bodySplit := strings.Split(splitLine, "/")
	if len(bodySplit) > 2 && bodySplit[1] == "echo" {
        fmt.Println("exe first")
		res := 
            fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(bodySplit[2]), bodySplit[2])
        conn.Write([]byte(res))
        return
	}
	fmt.Println(reqUrl)
	if reqUrl == "" {
        fmt.Println("exe second")
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			panic(err)
		}
	} else {
        fmt.Println("exe third")
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		if err != nil {
			panic(err)
		}
	}
}
