package zookeeper

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func InitConn(target string) (*zk.Conn, error) {

	addrs := strings.Split(target, ",")
	zkc, events, err := zk.Connect(addrs, time.Second*5)
	if err != nil {

		return nil, err
	}

	for {

		isConnected := false
		select {

		case connEvent := <-events:

			if connEvent.State == zk.StateConnected {

				isConnected = true
				log.Println("connect to zookeeper server success!")
			}

		case _ = <-time.After(time.Second * 5):

			log.Println("connect to zookeeper server timeout!")

			return nil, errors.New("connect to zookeeper server timeout!")
		}

		if isConnected {

			break
		}
	}

	return zkc, nil
}
