package indi.yuluo.envoy.xds.resource;

import indi.yuluo.envoy.xds.constant.Constants;

import com.google.common.collect.ImmutableList;
import com.google.protobuf.Any;
import com.google.protobuf.Duration;
import io.envoyproxy.controlplane.cache.v3.Snapshot;
import io.envoyproxy.envoy.config.cluster.v3.Cluster;
import io.envoyproxy.envoy.config.core.v3.*;
import io.envoyproxy.envoy.config.endpoint.v3.ClusterLoadAssignment;
import io.envoyproxy.envoy.config.endpoint.v3.Endpoint;
import io.envoyproxy.envoy.config.endpoint.v3.LbEndpoint;
import io.envoyproxy.envoy.config.endpoint.v3.LocalityLbEndpoints;
import io.envoyproxy.envoy.config.listener.v3.Filter;
import io.envoyproxy.envoy.config.listener.v3.FilterChain;
import io.envoyproxy.envoy.config.listener.v3.Listener;
import io.envoyproxy.envoy.config.route.v3.*;
import io.envoyproxy.envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager;
import io.envoyproxy.envoy.extensions.filters.http.router.v3.Router;

/**
 * @author yuluo
 * @author 1481556636@qq.com
 */

public class Resources {

    public Cluster makeCluster(String clusterName, String upstreamHost) {

        return Cluster.newBuilder()
                .setName(clusterName)
                .setType(Cluster.DiscoveryType.LOGICAL_DNS)
                .setConnectTimeout(Duration.newBuilder().setSeconds(3).build())
                .setLbPolicy(Cluster.LbPolicy.ROUND_ROBIN)
                .setLoadAssignment(makeEndpoint(clusterName, upstreamHost))
                .setDnsLookupFamily(Cluster.DnsLookupFamily.V4_ONLY)
                .build();
    }

    public ClusterLoadAssignment makeEndpoint(String clusterName, String upstreamName) {
        return ClusterLoadAssignment.newBuilder()
                .setClusterName(clusterName)
                .addEndpoints(LocalityLbEndpoints.newBuilder()
                        .addLbEndpoints(
                                LbEndpoint.newBuilder()
                                        .setEndpoint(Endpoint.newBuilder()
                                                .setAddress(Address.newBuilder()
                                                        .setSocketAddress(SocketAddress.newBuilder()
                                                                .setAddress(upstreamName)
                                                                .setPortValue(Constants.UPSTREAM_SERVER_PORT)
                                                                .setProtocol(SocketAddress.Protocol.TCP)
                                                                .build())
                                                        .build())
                                                .build())
                                        .build()
                        )
                        .build())
                .build();
    }

    public RouteConfiguration makeRoute(String routeName, String clusterName, String upstreamHost) {

        return RouteConfiguration.newBuilder()
                .setName(routeName)
                .addVirtualHosts(
                        VirtualHost.newBuilder()
                                .setName("xds_test_java_hosts")
                                .addDomains("*")
                                .addRoutes(
                                        Route.newBuilder()
                                                .setMatch(RouteMatch.newBuilder()
                                                        .setPrefix("/")
                                                        .build()
                                                )
                                                .setRoute(RouteAction.newBuilder()
                                                        .setCluster(clusterName)
                                                        .setHostRewriteLiteral(upstreamHost)
                                                        .build()
                                                )
                                                .build()
                                )
                )
                .build();

    }

    public ConfigSource makeConfigSource() {

        return ConfigSource.newBuilder()
                .setResourceApiVersion(ApiVersion.V3)
                .setApiConfigSource(ApiConfigSource.newBuilder()
                        .setTransportApiVersion(ApiVersion.V3)
                        .setApiType(ApiConfigSource.ApiType.GRPC)
                        .setSetNodeOnFirstMessageOnly(Boolean.TRUE)
                        .addGrpcServices(
                                GrpcService.newBuilder()
                                        .setEnvoyGrpc(GrpcService.EnvoyGrpc.newBuilder()
                                                .setClusterName(Constants.GRPC_STATIC_CLUSTER_NAME)
                                                .build())
                                        .build()
                        )
                        .build())
                .build();
    }

    public Listener makeHttpListener(String listenerName, String routeName) {

        HttpConnectionManager manager = HttpConnectionManager.newBuilder()
                .setCodecType(HttpConnectionManager.CodecType.AUTO)
                .setStatPrefix("http")
                .setRds(
                        io.envoyproxy.envoy.extensions.filters.network.http_connection_manager.v3.Rds
                                .newBuilder()
                                .setConfigSource(makeConfigSource())
                                .setRouteConfigName(routeName))
                .addHttpFilters(
                        io.envoyproxy.envoy.extensions.filters.network.http_connection_manager.v3.HttpFilter
                                .newBuilder()
                                .setName(Constants.FILTER_ENVOY_ROUTER)
                                .setTypedConfig(Any.pack(Router.newBuilder().build())))
                .build();


        return Listener.newBuilder()
                .setName(listenerName)
                .setAddress(Address.newBuilder()
                        .setSocketAddress(SocketAddress.newBuilder()
                                .setProtocol(SocketAddress.Protocol.TCP)
                                .setAddress(Constants.WILDCARD_HOST)
                                .setPortValue(Constants.LISTENER_PORT)
                                .build())
                        .build())
                .addFilterChains(FilterChain.newBuilder()
                        .addFilters(
                                Filter.newBuilder()
                                        .setName(Constants.FILTER_HTTP_CONNECTION_MANAGER)
                                        .setTypedConfig(Any.pack(manager))
                                        .build()
                        )
                        .build())
                .build();
    }

    public Snapshot generateSnapshot(String version, String upstreamName) {

        return Snapshot.create(
                ImmutableList.of(makeCluster(Constants.CLUSTER_NAME, upstreamName)),
                ImmutableList.of(makeEndpoint(Constants.CLUSTER_NAME, upstreamName)),
                ImmutableList.of(makeHttpListener(Constants.LISTENER_NAME, Constants.ROUTER_NAME)),
                ImmutableList.of(makeRoute(Constants.ROUTER_NAME, Constants.CLUSTER_NAME, upstreamName)),
                ImmutableList.of(),
                version
        );

    }

}
