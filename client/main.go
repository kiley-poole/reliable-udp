package main

import (
	"bufio"
	"crypto/md5"
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
	defer syscall.Close(socket)
	fmt.Println("Reliable UDP")
	fmt.Println("************")
	for {
		fmt.Print("$ ")
		if s.Scan() {
			data := s.Bytes()
			for {
				msg := buildMsg(data)
				transmit(msg, socket, port)
				exit := receiveValidation(socket)
				if exit {
					break
				}
			}

		}
	}
}

func socketBuild() int {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	check(err)

	err = syscall.Bind(socket, &syscall.SockaddrInet4{Port: 5050})
	check(err)

	return socket

}

func buildMsg(data []byte) []byte {

	checksum := md5.Sum(data)
	msg := append(checksum[:], data...)
	return msg
}

func transmit(msg []byte, socket int, port int) {
	err := syscall.Sendto(socket, msg, 0, &syscall.SockaddrInet4{Port: port})
	check(err)
}

//refactor this
func receiveValidation(socket int) bool {
	var fds syscall.FdSet
	fds.Bits[0] = 1 << uint(socket)

	timeout, err := syscall.Select(socket+1, &fds, nil, nil, &syscall.Timeval{Sec: 0, Usec: 500000})
	check(err)

	if timeout == 0 {
		return false
	}

	res := make([]byte, 1460)
	_, _, err = syscall.Recvfrom(socket, res, 0)
	check(err)
	fmt.Printf("%s\n", res)

	return true
}

func check(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}
