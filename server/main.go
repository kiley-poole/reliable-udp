package main

import (
	"fmt"
	"log"
	"syscall"
)

func main() {

	address := syscall.SockaddrInet4{Port: 9500}

	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	check(err)

	defer syscall.Close(socket)

	err = syscall.SetsockoptInt(socket, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	check(err)

	err = syscall.Bind(socket, &address)
	check(err)

	for {
		packet := make([]byte, 1460)
		_, _, err = syscall.Recvfrom(socket, packet, 0)
		fmt.Printf("Message: %s\n", string(packet))
		check(err)

		fmt.Println("Sending ACK")
		test := "ACK"
		err = syscall.Sendto(socket, []byte(test), 0, &syscall.SockaddrInet4{Port: 5050})
		check(err)

	}

}

func check(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}
