package indi.yuluo.envoy.xds;

import indi.yuluo.envoy.xds.constant.Constants;
import indi.yuluo.envoy.xds.resource.Resources;
import io.envoyproxy.controlplane.cache.v3.SimpleCache;
import io.envoyproxy.controlplane.cache.v3.Snapshot;
import io.envoyproxy.controlplane.server.V3DiscoveryServer;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.netty.NettyServerBuilder;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.ServerSocket;
import java.net.Socket;

/**
 * @author yuluo
 * @author 1481556636@qq.com
 */

public class Main {

    public static void main(String[] args) {

        SimpleCache<String> cache = new SimpleCache<>(node -> Constants.GROUP);

        Resources resources = new Resources();

        cache.setSnapshot(Constants.GROUP, resources.generateSnapshot("1", Constants.UPSTREAM_SERVER_1));

        V3DiscoveryServer v3DiscoveryServer = new V3DiscoveryServer(cache);
        ServerBuilder<NettyServerBuilder> builder =
                NettyServerBuilder.forPort(Constants.SERVER_PORT)
                        .addService(v3DiscoveryServer.getAggregatedDiscoveryServiceImpl())
                        .addService(v3DiscoveryServer.getClusterDiscoveryServiceImpl())
                        .addService(v3DiscoveryServer.getEndpointDiscoveryServiceImpl())
                        .addService(v3DiscoveryServer.getListenerDiscoveryServiceImpl())
                        .addService(v3DiscoveryServer.getRouteDiscoveryServiceImpl());
        Server server = builder.build();

        try {
            server.start();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }

        System.out.println("Server has started on port " + server.getPort());

        Runtime.getRuntime().addShutdownHook(new Thread(server::shutdown));

        try {
            Thread.sleep(10000);
        } catch (InterruptedException e) {
            throw new RuntimeException(e);
        }

        try {
            ServerSocket serverSocket = new ServerSocket(Constants.DYNAMIC_TEST_SERVER_PORT);
            System.out.println("更新 xds 资源服务启动，监听端口：" + serverSocket.getLocalPort());

            while (true) {
                Socket clientSocket = serverSocket.accept();
                BufferedReader in = new BufferedReader(new InputStreamReader(clientSocket.getInputStream()));
                String request = in.readLine().substring(4, 9);

                if ("/test".equals(request)) {
                  Snapshot snapshot = resources.generateSnapshot("2", Constants.UPSTREAM_SERVER_2);
                  cache.setSnapshot(Constants.GROUP, snapshot);
                  System.out.println("更新 v2 版本 xds 资源成功！");
                } else {
                    if ("/show".equals(request)) {
                        System.out.println(cache.getSnapshot(Constants.GROUP));
                    } else {
                        System.out.println("无效请求！");
                    }
                }

                in.close();
                clientSocket.close();
            }
        } catch (IOException e) {
            e.printStackTrace();
        }

        try {
            server.awaitTermination();
        } catch (InterruptedException e) {
            throw new RuntimeException(e);
        }
    }

}
