package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"time"
)

func main() {
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("Covert Tcp Server using sequence numbers")
	fmt.Println("\t\t\t Althaf Mohamed")
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++")


	sourceIpAddress := flag.String("srcIp", "192.168.50.2", "Source ip address")

	portToFilter := flag.Uint("p", 80, "-p 80, Port to filter messages on")

	fmt.Println("Listening on : ", *sourceIpAddress)

	conn, err := net.ListenIP("ip4:tcp", &net.IPAddr{ IP: net.ParseIP(*sourceIpAddress) })

	if err!=nil{
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte,1500)
	for{

		len,err :=conn.Read(buffer)

		if err==nil{
			data := buffer[:len]
			packet := gopacket.NewPacket(data,layers.LayerTypeIPv4,gopacket.Default)
			tcpPacket := packet.TransportLayer().(*layers.TCP)
			if tcpPacket.DstPort == layers.TCPPort(uint16(*portToFilter)){
				seqRecvd := tcpPacket.Seq
				messageRecvd := fmt.Sprintf("%c",seqRecvd)
				fmt.Println(time.Now().Unix(),":\t","Seq: ",seqRecvd," Recieving : ", messageRecvd)
			}
		}else{
			fmt.Println(err)
		}
	}

}