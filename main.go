package main

import (
	"fmt"
	"github.com/npat-efault/crc16"
	"go.bug.st/serial.v1"
	"time"

	//"go.bug.st/serial.v1"

	//"github.com/tarm/serial"
	"go.bug.st/serial.v1/enumerator"

	"log"
	//"os"
	//"time"
)

//Frame Type (1 byte).
//Addr Device (1 byte).
//Command (2 byte).
//Length of data (2 byte).
//Packet ID (2 byte)
//Reserved (2 byte)
//CRC (header) (2 byte).
//Data.
// CRC (data) (2 byte).

type Packet struct {
	frameType    byte    // 0x7B
	addresDevice byte    // пока неизвестен
	command      [2]byte //команда
	lenData      [2]byte
	packetId     [2]byte
	reserved     [2]byte
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

func SetByte(packet *Packet) []byte {
	data := make([]byte, 12)
	conf := crc16.PPP
	data[0] = packet.frameType
	data[1] = packet.addresDevice
	data[2] = packet.command[0]
	data[3] = packet.command[1]
	data[4] = packet.lenData[0]
	data[5] = packet.lenData[1]
	data[6] = packet.packetId[0]
	data[7] = packet.packetId[1]
	data[8] = packet.reserved[0]
	data[9] = packet.reserved[1]
	crc := crc16.Checksum(conf, data[:9])

	data[11] = uint8(crc)
	crc = crc >> 8
	data[10] = uint8(crc)
	return data
}

func ReadByte(s serial.Port) {
	buff := make([]byte, 100)
	for {
		i := 0
		n, err := s.Read(buff)
		if err != nil {
			log.Fatal(err)
			continue
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		for i < n {
			if buff[i] == 0x7B {
				i = i + 2
				if buff[i] == 0x80 {
					log.Println("Conn success")
					fmt.Printf("%v", string(buff[:n]))
				}
				log.Println(buff[i])
			}
			i++
		}
		//fmt.Printf("%v", string(buff[:n]))
	}
}
func main() {
	//var port string
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		return
	}

	packet := &Packet{
		frameType: 0x7B,
		command:   [2]byte{0x00, 0x00},
		packetId:  [2]byte{0x00, 0x01},
	}

	//if len(os.Args) < 2 {
	//	port = "/dev/tty"
	//} else {
	//	port = os.Args[0]
	//}
	//baud := 115200
	//c := &serial.Config{Name: port, Baud: baud, ReadTimeout: time.Second * 3}
	//s, err := serial.OpenPort(c)
	//if err != nil {
	//	log.Fatal(err)
	//}

	packet.addresDevice = 5 //

	data := SetByte(packet)
	log.Println(data)
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	for _, port := range ports {
		if port.IsUSB {
			fmt.Printf("Found port: %s\n", port.Name)
			fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB serial %s\n", port.SerialNumber)

			s, err := serial.Open(port.Name, mode)
			if err != nil {
				log.Fatal(err)
			}
			//err = s.SetDTR(false)
			//if err != nil {
			//	log.Println(err)
			//}

			//err = s.SetRTS(true)
			//if err != nil {
			//	log.Println(err)
			//}
			time.Sleep(time.Second / 2)
			n, err := s.Write(data)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(s.GetModemStatusBits())

			fmt.Printf("Sent %v bytes\n", n)
			buff := make([]byte, 100)
			for {
				i := 0
				n, err := s.Read(buff)
				if err != nil {
					log.Fatal(err)
					continue
				}
				if n == 0 {
					fmt.Println("\nEOF")
					break
				}
				for i < n {

					if buff[i] == 0x7B {
						i = i + 2
						if buff[i] == 0x80 {
							log.Println("Conn success")
							fmt.Printf("%v", string(buff[:n]))
						}
						log.Println(buff[i])
					}

					i++
				}
				//fmt.Printf("%v", string(buff[:n]))
			}
			packet.packetId[1] = 2
			packet.command[1] = 1
			data = SetByte(packet)

			n, err = s.Write(data)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second / 2)

		}

	}
	//s.Write(data)
	//time.Sleep(time.Second / 2)
	//buf := make([]byte, 128)
	//n, err := s.Read(buf)
	//if err == nil {
	//	if buf[2] == Ok && buf[3] == 0 {
	//		log.Println("Conn success, address:")
	//	}
	//}
	//log.Println(buf[:n])

}
