package resources

import (
	"time"

	"golang-xds-impl/pkg/constant"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	ClusterName  = constant.CLUSTER_NAME
	RouteName    = constant.ROUTER_NAME
	ListerName   = constant.LISTENER_NAME
	ListenerPort = constant.LISTENER_PORT
	UpstreamHost = constant.UPSTREAM_SERVER_1
	UpstreamPort = constant.UPSTREAM_SERVER_PORT
)

// CDS 配置
func makeCluster(clusterName string) *clusterv3.Cluster {

	return &clusterv3.Cluster{
		Name:                 ClusterName,
		ConnectTimeout:       durationpb.New(3 * time.Second),
		// 集群类型
		ClusterDiscoveryType: &clusterv3.Cluster_Type{Type: clusterv3.Cluster_LOGICAL_DNS},
		// 负载均衡策略
		LbPolicy:             clusterv3.Cluster_ROUND_ROBIN,
		LoadAssignment:       makeEndpoint(clusterName),
		DnsLookupFamily:      clusterv3.Cluster_V4_ONLY,
	}
}

// makeEndpoint 设置 endpoint
func makeEndpoint(clusterName string) *endpointv3.ClusterLoadAssignment {
	return &endpointv3.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpointv3.LocalityLbEndpoints{
			{
				LbEndpoints: []*endpointv3.LbEndpoint{
					{
						HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
							Endpoint: &endpointv3.Endpoint{
								Address: &core.Address{
									Address: &core.Address_SocketAddress{
										SocketAddress: &core.SocketAddress{
											// 需要代理的服务信息
											Protocol: core.SocketAddress_TCP,
											// host 地址
											Address:  UpstreamHost,
											PortSpecifier: &core.SocketAddress_PortValue{
												// 需要代理的服务端口，nginx 应用运行在 nginx-app 容器的 80 端口
												PortValue: uint32(UpstreamPort),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func makeRoute(routeName string, clusterName string) *routev3.RouteConfiguration {
	return &routev3.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*routev3.VirtualHost{
			{
				Name:    "xds_test_hosts",
				Domains: []string{"*"},
				Routes: []*routev3.Route{
					{
						Match: &routev3.RouteMatch{
							PathSpecifier: &routev3.RouteMatch_Prefix{
								// 路由名字，不设置重写的时候，会全部转发给 nginx-app 容器
								Prefix: "/",
							},
						},
						Action: &routev3.Route_Route{
							Route: &routev3.RouteAction{
								ClusterSpecifier: &routev3.RouteAction_Cluster{
									Cluster: clusterName,
								},
								HostRewriteSpecifier: &routev3.RouteAction_HostRewriteLiteral{
									HostRewriteLiteral: UpstreamHost,
								},
							},
						},
					},
				},
			},
		},
	}
}

func makeHTTPListener(listenerName string, route string) *listenerv3.Listener {

	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    makeConfigSource(),
				RouteConfigName: route,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: wellknown.Router,
			// fix: gRPC config for type.googleapis.com/envoy.config.listener.v3.Listener rejected: Error adding/updating listener(s) listener_0: Didn't find a registered implementation for 'envoy.filters.http.router' with type URL: ''
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: &anypb.Any{
					TypeUrl: "type.googleapis.com/envoy.extensions.filters.http.router.v3.Router",
				},
			},
		}},
	}
	pbst, err := anypb.New(manager)
	if err != nil {
		panic(err)
	}

	return &listenerv3.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  constant.WILDCARD_HOST,
					PortSpecifier: &core.SocketAddress_PortValue{
						// envoy 监听的端口，可以通过 ip://ListenerPort 访问代理的服务
						PortValue: uint32(ListenerPort),
					},
				},
			},
		},
		// 过滤器链扩展
		FilterChains: []*listenerv3.FilterChain{{
			Filters: []*listenerv3.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listenerv3.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}
}

func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					// 这里的 ClusterName 需要和 上面的 ClusterName 区分开，这里是配置文件中配置的 xds server 的 ClusterName. 
					// 如果写错 envoy 不会报错，但是 envoy 代理服务访问失败。
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: constant.GRPC_STATIC_CLUSTER_NAME},
				},
			}},
		},
	}
	return source
}

// GenerateSnapshot 创建缓存快照
// 真正存放不同xds配置的文件, 通过snapshot构建出不同资源类型的服务发现资源(其实就是构造不同的配置结构体)
func GenerateSnapshot(version string) *cache.Snapshot {

	snap, _ := cache.NewSnapshot(version,
		map[resource.Type][]types.Resource{
			resource.ClusterType:  {makeCluster(ClusterName)},
			resource.RouteType:    {makeRoute(RouteName, ClusterName)},
			resource.ListenerType: {makeHTTPListener(ListerName, RouteName)},
		},
	)

	return snap
}
