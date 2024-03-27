package indi.yuluo.envoy.xds.constant;

/**
 * @author yuluo
 * @author 1481556636@qq.com
 */

public interface Constants {

    String LOCALHOST = "127.0.0.1";

    int SERVER_PORT = 9090;

    String WILDCARD_HOST = "0.0.0.0";

    int DYNAMIC_TEST_SERVER_PORT = 19090;

    String UPSTREAM_HOST = "127.0.0.1";

    String NETWORK_MODEL = "tcp";

    String NODE_ID = "test_1";

    String TEST_URI = "/test";

    String CLUSTER_NAME = "xds_nginx_cluster";

    String LISTENER_NAME = "listener_0";

    int LISTENER_PORT = 10000;

    String UPSTREAM_SERVER_1 = "nginx-app-1";

    String UPSTREAM_SERVER_2 = "nginx-app-2";

    int UPSTREAM_SERVER_PORT = 80;

    String ROUTER_NAME = "local_router";

    String FILTER_ENVOY_ROUTER = "envoy.filters.http.router";

    String FILTER_HTTP_CONNECTION_MANAGER = "envoy.http_connection_manager";

    String GROUP = "key";

    String GRPC_STATIC_CLUSTER_NAME = "grpc_xds_cluster";
}
