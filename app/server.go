package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var directoryPath string

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
    fmt.Println("handle user agent")
	header := req.Header.Get("User-Agent")
	res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(header), header)
	sendResponse(conn, res)
}

func handleFileRequest(conn net.Conn, reqUrl string) {
    fmt.Println("handle file req is exe")
	fileName := strings.Split(reqUrl, "/")[1]

	// checking if file exists or not
	if _, err := os.Stat(directoryPath + fileName); os.IsNotExist(err) {
		sendResponse(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
		return
	} else if err != nil {
		handleError(err, "Failed to check if file exists")
	}

	// opening the file
	file, err := os.Open(directoryPath + fileName)
	if err != nil {
		handleError(err, "Failed while opening file")
	}

	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		handleError(err, "Failed to get file info")
	}

	f, err := os.ReadFile(directoryPath + fileName)
	if err != nil {
		handleError(err, "Failed to read file")
	}
	fileContent := string(f)
	res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", fileInfo.Size(), fileContent)
	sendResponse(conn, res)
}

func handlePOSTRequest(conn net.Conn, reqUrl string, body string) {
    fmt.Println("handle post req is exe ")
	fileName := strings.Split(reqUrl, "/")[1]

	// create a file with this name
	_, err := os.Create(directoryPath + fileName)
	handleError(err, "Failed to create file")
	err = os.WriteFile(directoryPath+fileName, []byte(body), 0644)
	handleError(err, "Failed to write to file")
	sendResponse(conn, "HTTP/1.1 201 Created\r\n\r\n")
}

func handleEncoding(conn net.Conn, header []string) {
    fmt.Println("handle encoding")
	var res string
	var foundGzip bool

	for _, h := range header {
        fmt.Println(h)
		if h == "gzip" {
			foundGzip = true
			break
		}
	}
	if foundGzip {
		res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: %s\r\n\r\n", "gzip")
	} else {
		res = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n"
	}
	sendResponse(conn, res)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	handleError(err, "Failed to read request")

	reqUrl := req.URL.Path[1:]
	splitLine := strings.Split(reqUrl, "/")

	filereq, err := regexp.MatchString("^files/.+$", reqUrl)
	handleError(err, "Failed to match regex")

	// get the body from the request
	body, err := io.ReadAll(req.Body)
	handleError(err, "Failed to read body")

	// get the headers
	header := strings.Split(req.Header.Get("Accept-Encoding"), ",")
    var finalHeader []string
    for _, h := range header {
        temp := strings.TrimSpace(h)
        if temp != "" {
            finalHeader = append(finalHeader, temp)
        }
    }

    fmt.Println("final header len is '", len(finalHeader))

	if len(finalHeader) > 0 {
		handleEncoding(conn, header)
	} else if reqUrl == "user-agent" {
		handleUserAgent(conn, req)
	} else if filereq && string(body) != "" {
		handlePOSTRequest(conn, reqUrl, string(body))
	} else if filereq {
		handleFileRequest(conn, reqUrl)
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

	dirFlag := flag.String("directory", "", "the directory to serve files from")
	flag.Parse()
	directoryPath = *dirFlag
	if directoryPath != "" {
		if _, err = os.Stat(*dirFlag); os.IsNotExist(err) {
			err = os.MkdirAll(*dirFlag, 0755)
			handleError(err, "Failed to create directory")
		} else if err != nil {
			handleError(err, "Failed to check if directory exists")
		}
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		handleError(err, "Failed to accept connection")

		go handleConnection(conn)
	}
}
