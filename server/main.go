package main

import (
	"crypto/md5"
	"encoding/binary"
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
		check(err)

		dataLength := packet[0:2]
		len := binary.BigEndian.Uint16(dataLength) + 18
		var checksumHeader [16]byte
		copy(checksumHeader[:], packet[2:18])
		data := packet[18:len]
		checksum := md5.Sum(data)

		fmt.Printf("Checksum Header: %x\n", string(checksumHeader[:]))
		fmt.Printf("Checksum Validate: %x\n", string(checksum[:]))
		fmt.Printf("Data: %s\n", string(data))

		if checksum != checksumHeader {
			fmt.Println("DATA CORRUPTED INVALID CHECKSUM. GET OUTTA HERE WITH THAT GARBAGE.")
		} else {
			fmt.Println("Sending ACK")
			ack := []byte("ACK")
			err = syscall.Sendto(socket, ack, 0, &syscall.SockaddrInet4{Port: 5050})
			check(err)
		}
		fmt.Println()
	}

}

func check(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}
