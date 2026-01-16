package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8083"
	SERVER_TYPE = "tcp"
)

func main() {
	fmt.Print("Введіть свій нікнейм: ") //commit
	userReader := bufio.NewReader(os.Stdin)
	nickname, _ := userReader.ReadString('\n')

	conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Помилка підключення:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Підключено до TCP Chat-сервера на", SERVER_HOST+":"+SERVER_PORT)

	go func() {
		serverReader := bufio.NewReader(conn)
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil {
				break
			}
			fmt.Print(msg)
		}
	}()

	for {
		text, _ := userReader.ReadString('\n')
		if len(text) > 1 {
			message := fmt.Sprintf("[%s]: %s", nickname[:len(nickname)-1], text)
			conn.Write([]byte(message))
		}
	}
}
