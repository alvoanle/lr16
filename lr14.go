package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	HOST                  = "localhost"
	PORT                  = "8083"
	TYPE                  = "tcp"
	CLIENT_CHANNEL_BUFFER = 10
)

var (
	messages = make(chan string)
	joining  = make(chan chan<- string)
	leaving  = make(chan chan<- string)
)

func manager() {
	clients := make(map[chan<- string]bool)

	for {
		select {
		case msg := <-messages:
			for clientChan := range clients {
				select {
				case clientChan <- msg:
				default:
					leaving <- clientChan
					delete(clients, clientChan)
					fmt.Println("Клієнт від'єднаний через переповнення буфера.") //коментарій
				}
			}

		case clientChan := <-joining:
			clients[clientChan] = true
			fmt.Printf("Chat: Новий клієнт приєднався. Активних: %d\n", len(clients))

		case clientChan := <-leaving:
			delete(clients, clientChan)
			close(clientChan)
			fmt.Printf("Chat: Клієнт від'єднався. Активних: %d\n", len(clients))
		}
	}
}

func clientWriter(conn net.Conn, out <-chan string) {
	defer conn.Close()
	for msg := range out {
		_, err := conn.Write([]byte(msg))
		if err != nil {
			break
		}
	}
}

func handleConnection(conn net.Conn) {
	outgoing := make(chan string, CLIENT_CHANNEL_BUFFER)
	go clientWriter(conn, outgoing)
	joining <- outgoing

	clientAddr := conn.RemoteAddr().String()
	input := bufio.NewScanner(conn)

	for input.Scan() {
		messages <- fmt.Sprintf("[%s]: %s\n", clientAddr, input.Text())
	}
	leaving <- outgoing
}

func main() {
	go manager()

	listener, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Помилка запуску сервера: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("TCP Chat Сервер запущено на %s:%s. Очікування клієнтів...\n", HOST, PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Помилка Accept: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}
