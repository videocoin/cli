package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/rtmp"
)

func getIpAddress() string {
	i, err := net.InterfaceByName("en0")
	if err != nil {
		fmt.Print(err)
		return ""
	}

	if addrs, err := i.Addrs(); err == nil {
		for _, addr := range addrs {
			switch ip := addr.(type) {
			case *net.IPNet:
				if ip.IP.DefaultMask() != nil {
					return ip.IP.String()
				}
			}
		}
	}

	return ""
}

func main() {
	port := ":1936"
	server := &rtmp.Server{Addr: "0.0.0.0" + port}
	l := &sync.RWMutex{}
	type Channel struct {
		que *pubsub.Queue
	}
	channels := map[string]*Channel{}

	server.HandlePlay = func(conn *rtmp.Conn) {
		defer func() {
			fmt.Println("client connection closed")
			conn.Close()
		}()
		l.RLock()
		ch := channels[conn.URL.Path]
		l.RUnlock()

		fmt.Println("new client connection", conn.URL.Path)
		if ch != nil {
			cursor := ch.que.Latest()
			_ = avutil.CopyFile(conn, cursor)
		} else {
			fmt.Printf("stream %s was not found\n", conn.URL.Path)
		}
	}

	server.HandlePublish = func(conn *rtmp.Conn) {
		defer func() {
			fmt.Printf("publish connection %s closed\n", conn.URL.Path)
			conn.Close()

			l.Lock()
			delete(channels, conn.URL.Path)
			l.Unlock()
		}()

		streams, err := conn.Streams()
		if err != nil {
			panic("failed to request incoming streams")
		}

		l.Lock()
		ch := channels[conn.URL.Path]
		if ch == nil {
			ch = &Channel{}
			ch.que = pubsub.NewQueue()
			err := ch.que.WriteHeader(streams)
			if err != nil {
				panic("failed to write headers")
			}
			channels[conn.URL.Path] = ch
		} else {
			ch = nil
		}
		l.Unlock()
		if ch == nil {
			return
		}

		fmt.Printf("new publish connection %s\n", conn.URL.Path)

		err = avutil.CopyPackets(ch.que, conn)
		if err != nil {
			panic("failed to copy packets")
		}
	}

	fmt.Printf("rtmp server is listening on %s%s\n", getIpAddress(), port)
	err := server.ListenAndServe()
	if err != nil {
		panic("failed to start rtmp server")
	}
}
