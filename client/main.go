package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"syscall"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter port: ")
	var port int
	_, err := fmt.Scan(&port)
	check(err)

	socket := socketBuild()
	fmt.Println("Reliable UDP")
	fmt.Println("************")
	for {
		fmt.Print("$ ")
		if s.Scan() {
			msg := s.Bytes()
			transmit(msg, socket, port)
		}
	}
}

func socketBuild() int {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	check(err)

	err = syscall.Bind(socket, &syscall.SockaddrInet4{})
	check(err)

	return socket

}
func transmit(msg []byte, socket int, port int) {
	err := syscall.Sendto(socket, msg, 0, &syscall.SockaddrInet4{Port: port})
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}
