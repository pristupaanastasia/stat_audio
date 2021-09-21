package main

import (
	"fmt"
	"github.com/npat-efault/crc16"
	"github.com/tarm/serial"
	"log"
	"os"
	"time"
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
	var port string

	packet := &Packet{
		frameType: 0x7B,
		command:   [2]byte{0x00, 0x00},
	}

	if len(os.Args) < 2 {
		port = "COM3"
	} else {
		port = os.Args[0]
	}
	baud := 115200
	c := &serial.Config{Name: port, Baud: baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	packet.addresDevice = 5 //

	data := make([]byte, 8)
	conf := crc16.PPP
	crc := crc16.Checksum(conf, []byte(fmt.Sprintf("%v", packet)))
	data[0] = packet.frameType
	data[1] = packet.addresDevice
	data[2] = packet.command[0]
	data[3] = packet.command[1]
	data[4] = packet.lenData[0]
	data[5] = packet.lenData[1]
	data[7] = uint8(crc)
	crc = crc >> 8
	data[6] = uint8(crc)
	log.Println(data)
	s.Write(data)

	t := time.Tick(time.Second * 3)
	var i uint8
	for {
		select {
		case <-t:
			break
		default:
			_, err := s.Read(data)
			if err == nil {
				if data[2] == Ok && data[3] == 0 {
					log.Println("Conn success, address:", i)
					break
				}
			}
			log.Println(data)
		}
	}

}
