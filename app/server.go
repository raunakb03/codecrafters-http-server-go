package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

func handleError(err error, message string) {
	if err != nil {
		fmt.Println(message)
		fmt.Println(err)
		os.Exit(1)
	}
}

func sendResponse(conn net.Conn, res string) {
	_, err := conn.Write([]byte(res))
	handleError(err, "Failed to send response")
}

func handleUserAgent(conn net.Conn, req *http.Request) {
	header := req.Header.Get("User-Agent")
	res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(header), header)
	sendResponse(conn, res)
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	handleError(err, "Failed to read request")

	reqUrl := req.URL.Path[1:]
	splitLine := strings.Split(reqUrl, "/")

	if reqUrl == "user-agent" {
		handleUserAgent(conn, req)
	} else if len(splitLine) > 1 && splitLine[0] == "echo" {
		res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(splitLine[1]), splitLine[1])
		sendResponse(conn, res)
	} else {
		if reqUrl == "" {
			sendResponse(conn, "HTTP/1.1 200 OK\r\n\r\n")
		} else {
			sendResponse(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
		}
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	handleError(err, "Failed to bind to port 4221")
    defer l.Close()

	for {
		conn, err := l.Accept()
		handleError(err, "Failed to accept connection")

		go handleConnection(conn)
	}
}
