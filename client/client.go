package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	addr := "localhost:1234"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("failed to connect to server at %s: %v", addr, err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Anonymous"
	}

	_, err = fmt.Fprintf(conn, "%s\n", name)
	if err != nil {
		log.Fatalf("failed to send name: %v", err)
	}

	fmt.Printf("Welcome %s! Type messages and press Enter to send.\n", name)
	fmt.Println("Type 'exit' to disconnect.")

	go func() {
		serverReader := bufio.NewReader(conn)
		for {
			line, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("Disconnected from server.")
				os.Exit(0)
			}
			line = strings.TrimRight(line, "\r\n")
			if line != "" {
				fmt.Println(line)
			}
		}
	}()

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if strings.EqualFold(text, "exit") {
			fmt.Println("Bye!")
			fmt.Fprintf(conn, "exit\n")
			return
		}
		if text == "" {
			continue
		}

		_, err := fmt.Fprintf(conn, "%s\n", text)
		if err != nil {
			fmt.Printf("send error: %v\n", err)
			return
		}
	}
}
