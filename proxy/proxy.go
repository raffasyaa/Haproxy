package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

const (
	buflen        = 4096 * 4
	timeout       = 60
	defaultHost   = "127.0.0.1:22"
	response      = "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: upgrade\r\nSec-WebSocket-Accept: foo\r\n\r\n"
)

var (
	listeningAddr = "127.0.0.1"
	listeningPort = 700
	password      = ""
)

type connectionHandler struct {
	client net.Conn
	server net.Conn
}

func main() {
	fmt.Println(":-------Goproxy-------:")
	fmt.Printf("Listening addr: %s\n", listeningAddr)
	fmt.Printf("Listening port: %d\n", listeningPort)
	fmt.Println(":---------------------:")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listeningAddr, listeningPort))
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(clientConn)
	}
}

func handleClient(clientConn net.Conn) {
	defer clientConn.Close()

	clientBuffer := make([]byte, buflen)
	_, err := clientConn.Read(clientBuffer)
	if err != nil {
		fmt.Println("Error reading client request:", err)
		return
	}

	hostPort := findHeader(clientBuffer, "X-Real-Host")

	if hostPort == "" {
		hostPort = defaultHost
	}

	split := findHeader(clientBuffer, "X-Split")

	if split != "" {
		clientConn.Read(clientBuffer) // Read and discard the extra data
	}

	if hostPort != "" {
		passwd := findHeader(clientBuffer, "X-Pass")

		if len(password) != 0 && passwd == password {
			methodConnect(clientConn, hostPort)
		} else if len(password) != 0 && passwd != password {
			clientConn.Write([]byte("HTTP/1.1 400 Wrong Password!\r\n\r\n"))
		} else if strings.HasPrefix(hostPort, "127.0.0.1") || strings.HasPrefix(hostPort, "localhost") {
			methodConnect(clientConn, hostPort)
		} else {
			clientConn.Write([]byte("HTTP/1.1 403 Forbidden!\r\n\r\n"))
		}
	} else {
		fmt.Println("- No X-Real-Host!")
		clientConn.Write([]byte("HTTP/1.1 400 No X-Real-Host!\r\n\r\n"))
	}
}

func findHeader(header []byte, name string) string {
	headerStr := string(header)
	index := strings.Index(headerStr, name+": ")

	if index == -1 {
		return ""
	}

	headerStr = headerStr[index+len(name)+2:]
	index = strings.Index(headerStr, "\r\n")

	if index == -1 {
		return ""
	}

	return headerStr[:index]
}

func connectTarget(host string) (net.Conn, error) {
	i := strings.Index(host, ":")
	if i != -1 {
		port := host[i+1:]
		host = host[:i]
	} else {
		if strings.HasPrefix(host, "127.0.0.1") || strings.HasPrefix(host, "localhost") {
			port = "22"
		} else {
			port = "80"
		}
	}

	server, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error connecting to target:", err)
		return nil, err
	}

	return server, nil
}

func methodConnect(clientConn net.Conn, path string) {
	fmt.Printf(" - Connect %s\n", path)

	server, err := connectTarget(path)
	if err != nil {
		clientConn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
		return
	}
	defer server.Close()

	clientConn.Write([]byte(response))
	clientBuffer := make([]byte, buflen)

	for {
		_, err := clientConn.Read(clientBuffer)
		if err != nil {
			fmt.Println("Error reading client data:", err)
			return
		}

		_, err = server.Write(clientBuffer)
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}

		_, err = server.Read(clientBuffer)
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}

		_, err = clientConn.Write(clientBuffer)
		if err != nil {
			fmt.Println("Error writing to client:", err)
			return
		}
	}
}