package zookeeper

import (
	"log"

	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc/resolver"
)

const (
	grpcScheme  = "grpczk"
	GrpcDialUrl = "grpczk:///zookeeper.grpc.io"
)

func NewGrpcResolver(target, sever string) resolver.Builder {

	return &GrpcResolver{Target: target, Server: sever}
}

type GrpcResolver struct {
	Target string
	Server string
	cc     resolver.ClientConn
	zkc    *zk.Conn
}

func (r *GrpcResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {

	r.cc = cc
	err := r.initZkc()
	if err != nil {

		return nil, err
	}

	go r.watch()

	return r, nil
}

func (r *GrpcResolver) initZkc() error {

	zkc, err := InitConn(r.Target)
	if err != nil {

		return err
	}

	r.zkc = zkc

	return nil
}

func (r *GrpcResolver) updateState(addrs []string) {

	updateAddrs := make([]resolver.Address, len(addrs))
	for i, s := range addrs {

		updateAddrs[i] = resolver.Address{Addr: s}
	}

	r.cc.UpdateState(resolver.State{Addresses: updateAddrs})

	log.Println(r.Server, "addrs", addrs)
}

func (r *GrpcResolver) watch() {

	for {

		addrs, _, wch, err := r.zkc.ChildrenW("/" + schema + "/" + r.Server)

		if err == nil {

			r.updateState(addrs)
		} else {

			log.Println(grpcScheme, "watch", err)
		}

		for ev := range wch {

			if ev.Type == zk.EventNodeChildrenChanged {

				continue
			}
		}
	}
}

func (r *GrpcResolver) Scheme() string {

	return grpcScheme
}

func (r *GrpcResolver) ResolveNow(rn resolver.ResolveNowOption) {}

func (r *GrpcResolver) Close() {}
