package main

import (
	"fmt"
	"github.com/coreos/go-iptables/iptables"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"time"
)
const (
	defaultSnapLen = 262144
)

type SyncPacket struct{
	timeRecieved time.Time
	nextSeq uint32
	port uint16
}



func main(){

	syncPacketsList := make(map[string]SyncPacket)

	iptable,err := iptables.NewWithProtocol(iptables.ProtocolIPv4)

	if err!=nil{
		panic(err)
	}

	handle, err := pcap.OpenLive("enp0s8", defaultSnapLen, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter("port 80"); err != nil {
		panic(err)
	}

	// below is to remove ip tables after 5 minutes
	packets := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()
	ticker := time.NewTicker(15 * time.Second)
	quit := make(chan struct{})
	defer close(quit)
	go func() {
		for {
			select {
			case <- ticker.C:
				fmt.Print(".")
				for add,val:= range syncPacketsList{
					if time.Now().Sub(val.timeRecieved) >= 5 * time.Minute{
						iptable.Delete("filter","INPUT",
							"-s",
							fmt.Sprintf("%s",add),
							"-p",
							"tcp",
							"--destination-port",
							fmt.Sprintf("%d",val.port),
							"--jump","DROP")
						fmt.Print("Removing Rule")
						delete(syncPacketsList,add)
					}
				}
			case <- quit:
				err := iptable.DeleteAll()
				if err !=nil{
					panic(err)
				}
				ticker.Stop()
				return
			}
		}
	}()
	for pkt := range packets {
		ipPacket := pkt.NetworkLayer().(*layers.IPv4)
		srcIpString := ipPacket.SrcIP.String()
		tcpPacket := pkt.TransportLayer().(*layers.TCP)
		destPort := tcpPacket.DstPort
		if tcpPacket.SYN{
			if lastpacket,ok:=syncPacketsList[srcIpString]; ok == true{
				if lastpacket.nextSeq != tcpPacket.Seq {
					iptablErr := iptable.AppendUnique("filter","INPUT",
						"-s",
						fmt.Sprintf("%s",srcIpString),
						"-p",
						"tcp",
						"--destination-port",
						fmt.Sprintf("%d",destPort),
						"--jump","DROP","--wait")
					if iptablErr!=nil{
						panic(iptablErr)
					}
					lastpacket.timeRecieved = time.Now()

				}
			}else{
				syncPacketsList[srcIpString] = SyncPacket{
					nextSeq: tcpPacket.Seq + uint32(len(tcpPacket.Payload)) + 1,
					timeRecieved: time.Now(),
					port: uint16(tcpPacket.DstPort),
				}
			}
		}

	}

}
