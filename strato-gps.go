package main

import (
	"flag"
	"fmt"
	"github.com/stratoberry/gps"
	"log"
	"net"
	"os"
	"time"
)

var outFile = flag.String("output", "/data/gps.csv", "output filename")
var udpPort = flag.Int64("udp", 4773, "UDP port for sending data")

var updateFreq = flag.Int("freq", 10, "update frequency in seconds")

// Feb 1st 2014
const EPOCH = 1391212800

func main() {
	var dev *gps.Device
	var err error
	if dev, err = gps.Open("/dev/ttyAMA0"); err != nil {
		log.Panicln("Failed to open GPS device", err)
	}

	dev.Watch()

	f, err := os.OpenFile(*outFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Panicln("Failed to open output file", err)
	}
	defer f.Close()

	conn, _ := net.Dial("udp", fmt.Sprintf("255.255.255.255:%d", *udpPort))
	defer conn.Close()

	var lastUpdate int64 = 0
	for fix := range dev.Fixes {
		now := time.Now().Unix()
		if now-lastUpdate > *updateFreq {
			lastUpdate = now
			stanza := fmt.Sprintf("%d;%.8f;%.8f;%.2f;%.2f\n", now-EPOCH, fix.Lat, fix.Lon, fix.Alt, fix.TrackAngle)

			f.WriteString(stanza)
			conn.Write([]byte(stanza))
		}
	}
}
