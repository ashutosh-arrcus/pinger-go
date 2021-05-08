package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"os"
	"os/signal"
	"time"
)

func getStats(p *ping.Pinger, done <-chan struct{}) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
			fmt.Println("--------------------------------")
			fmt.Println("Current time: ", t)
			fmt.Println("Stats: ", p.Statistics())
			fmt.Println("Pkt Sent: ", p.Statistics().PacketsSent)
			fmt.Println("Pkt Recv: ", p.Statistics().PacketsRecv)
			return
		}
	}

}

func main() {
	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
		panic(err)
	}
	//pinger.Interval = 2
	pinger.Count = 1
	//pinger.Timeout = 10
	//pinger.Debug = true
	pinger.RecordRtts = true

	doneCh := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			pinger.Stop()
		}
		close(doneCh)
	}()

	go getStats(pinger, doneCh)
	err = pinger.Run()
	if err != nil {
		panic(err)
	}
	//pinger.Stop()
	//fmt.Println(pinger.Statistics())
	<-doneCh
}
