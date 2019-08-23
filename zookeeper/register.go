package zookeeper

import (
	"log"
	"sync"

	"github.com/samuel/go-zookeeper/zk"
)

const schema = "go-modian"

var swg sync.WaitGroup
var register = make(chan struct{})

func Register(target, server, value string) error {

	zkc, err := InitConn(target)
	if err != nil {

		return err
	}

	path := "/" + schema + "/" + server + "/" + value

	_, err = zkc.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
	if err != nil {

		return err
	}

	log.Println("register =>", path)

	swg.Add(1)

	go func() {

		for range register {
		}

		err := zkc.Delete(path, -1)
		if err == nil {

			log.Println("unregister =>", path)
		}

		swg.Done()
	}()

	return nil
}

func UnRegister() {

	close(register)
	swg.Wait()
}
