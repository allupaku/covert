package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"time"
)

func main() {

	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("Covert Tcp Client using sequence numbers")
	fmt.Println("\t\t\t Althaf Mohamed")
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

	sourceIpAddress := flag.String("srcIp", "192.168.50.3", "Source ip address")

	destIpAddress := flag.String("destIp", "192.168.50.2", "Source ip address")

	portToSend := flag.Uint("p", 80, "-p 80, Port to send messages to")

	net.ParseIP(*sourceIpAddress)

	fmt.Println("Sending from source IP ",*sourceIpAddress )

	fmt.Println("Sending to destination IP ",*destIpAddress )

	fmt.Println("Sending to port ", *portToSend)

	srcIp := net.ParseIP(*sourceIpAddress)

	destIp := net.ParseIP(*destIpAddress)



	destPort := layers.TCPPort(uint16(*portToSend))

	dataToSend, err := ioutil.ReadFile("message.txt")

	fmt.Println("Message to send :" ,string(dataToSend))

	if err!=nil{
		panic(err)
	}

	packetBuffer := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		 FixLengths: true,
		 ComputeChecksums: true,
	}

	ip := layers.IPv4{
		Version:    4,
		Protocol:   layers.IPProtocolIPv4,
		SrcIP:      srcIp,
		DstIP:      destIp,
	}
	conn, err := net.Dial("ip4:tcp", destIp.String())
	if err!=nil{
		panic(err)
	}
	defer conn.Close()
	for _,c := range dataToSend{
		rand.Seed(time.Now().UnixNano())

		sourcePort := layers.TCPPort(rand.Intn( 65535-1025) + 1025)
		fmt.Println("Sending : ", c, "(" ,fmt.Sprintf("%c",c),")")

		tcp := layers.TCP{
			SrcPort:    sourcePort,
			DstPort:    destPort,
			Seq:        uint32(c),
			SYN:        true,
		}
		gopacket.SerializeLayers(packetBuffer, opts,
			&ip,
			&tcp,
			gopacket.Payload([]byte("XX")))
		packetData := packetBuffer.Bytes()


		if err != nil {
			log.Fatalf("Dial: %s\n", err)
		}
		conn.Write(packetData)
		time.Sleep(time.Millisecond * 100)
	}

}