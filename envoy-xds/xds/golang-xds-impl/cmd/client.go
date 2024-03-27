package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"golang-xds-impl/pkg/constant"

	envoyconfigclusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoyconfigcorev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// client 代码，请求获取配置

func main() {

	gOpts := []grpc.DialOption{

		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// XDS server addr
	addr := constant.LOCALHOST + ":" + strconv.Itoa(constant.SERVER_PORT)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, gOpts...)
	if err != nil {

		log.Fatalf("DialContext err is '%s'. \n", err.Error())
	}

	client := clusterservice.NewClusterDiscoveryServiceClient(conn)
	req := &discovery.DiscoveryRequest{
		Node: &envoyconfigcorev3.Node{
			Id: constant.NODE_ID,
		},
	}
	rsp, err := client.FetchClusters(context.Background(), req)
	if err != nil {

		log.Fatalf("FetchClusters err is '%s'. \n", err.Error())
	}

	getResource := rsp.GetResources()[0]
	cluster := &envoyconfigclusterv3.Cluster{}

	err = getResource.UnmarshalTo(cluster)
	if err != nil {

		log.Fatalf("UnmarshalTo err is: '%s' \n", err.Error())
	}

	fmt.Println("cluster: ", cluster)

}
