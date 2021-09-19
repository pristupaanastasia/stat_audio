package main

import (
	"bufio"
	"fmt"
	"github.com/npat-efault/crc16"
	"log"
	"net"
	"os"
	"strconv"
)

//Frame Type (1 byte).
//Addr Device (1 byte).
//Command (2 byte).
//Length of data (2 byte).
//CRC (header) (2 byte).
//Data.
// CRC (data) (2 byte).

type Packet struct {
	frameType    byte    // 0x7B
	addresDevice byte    // пока неизвестен
	command      [2]byte //команда
	lenData      [2]byte
	crcHeader    [2]byte
}

const Ok byte = 0x80

//Command
// Пинг (0x0000).
// Блокировка клавиатуры прибора (0x0001).
// Чтение версий прибора (0x0002).
// Установка скорости обмена по UART (0x0003).
// Установка нового адреса прибора (0x0004).
// Перезагрузка (RESET) прибора (0x0005).
// Выключение питания прибора (0x0006).
// Читаем номер прибора (0x0009).
// Установка времени (0x0100).
// Установка даты (0x0101).
// Настройка шаблонов для чтения данных (0x0102).
func main() {
	ln, _ := net.Listen("tcp", ":8081")

	packet := &Packet{
		frameType: 0x7B,
		command:   [2]byte{0x0000, 0x0000},
		crcHeader: [2]byte{0x0000, 0x0000},
	}
	address, err := strconv.Atoi(os.Args[0])
	if err != nil {
		log.Println(err)
	}
	packet.addresDevice = byte(address)
	conn, _ := ln.Accept()

	data := make([]byte, 8)
	conf := crc16.PPP
	buf := crc16.Checksum(conf, []byte(fmt.Sprintf("%v", packet)))
	packet.crcHeader[0] = byte(buf)
	buf = buf >> 8
	packet.crcHeader[1] = byte(buf)
	conn.Write([]byte(fmt.Sprintf("%v", packet)))
	for {
		_, err := bufio.NewReader(conn).Read(data)
		if err != nil {
			continue
		}
		if data[2] == Ok && data[3] == 0 {
			log.Println("Conn success")
			break
		}
	}

	for {
		// Будем прослушивать все сообщения разделенные \n
		n, err := bufio.NewReader(conn).Read(data)
		if err != nil {
			continue
		}
		// Распечатываем полученое сообщение
		fmt.Print("Message Received:", data[:n])
		// Процесс выборки для полученной строки
		// Отправить новую строку обратно клиенту

	}

}
