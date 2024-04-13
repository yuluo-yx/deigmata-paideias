# Kubernetes Gateway API

## 诞生背景

在 kubernetes 中，流量的治理主要分为两个部分：

1. 南北向流量
2. 东西向流量

### 南北向流量（NORTH-SOUTH traffic）

> 在计算机网络中，南北向流量通常指数据流量从一个**内部网络（局域网）流向外部网络（公网）**的流量。
>
> 在 kubernetes 场景中，南北向流量为：主要指从集群外到集群内的流量。

#### kubernetes 南北向流量

客户端想要访问部署在 Kubernetes 集群中的服务时， 需要将集群内的服务暴露出来，最常用的方式就是通过 NodePort 或 LoadBalancer 类型的 Service。这两种方式相对来说都比较简单，但是如果用户有很多服务需要暴露到集群外，通过使用这些方式就会浪费很多的端口（NodePort）或者浪费很多 IP（Loadbalancer）， 而且由于 Service 的 API 在设计时定位相对清晰，并不包含一些 LB 或者 gateway 的能力，比如根据域名进行代理，认证和鉴权等能力，仅仅只通过 Service 来将集群内服务暴露出去，非常不合理。

基于如上分析，kubernetes 社区新增了一种资源配置：Ingress。资源定义非常之简单，仅仅只能指定 Host、Path 和应用的 Service、Port 和 使用的协议。对于许多场景，Ingress 只是达到了基本可用的状态，没有包含其他场景的功能。例如：

- 根据 Request Header/Method 进行匹配
- Path Rewrite
- ……

#### Ingress 的工作原理

用户在 Kubernetes 集群中创建 Ingress 资源之后，并没有像 Deployment、Pod 等资源那样，一并实现其对应的 Controller。所以需要一个 ingress controller 这样的组件存在，将 ingress 资源定义翻译成为数据面可以执行的配置，并在数据面生效。当前仅列出的 Ingress controller 就有 30 种：

https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/

同时，由于 Ingress API 只定义了有限的内容，它的表现力不够，为了满足不同场景下的需求，各个 Ingress controller 实现的时候，只好通过创建自己的 CRD（Custom Resource Definition）或者通过为 Ingress 资源增加 annotations 的方式来实现对应的需求。

这里的角色定位是：

> Ingress 资源 ==> 反向代理的规则，用于规定 HTTP 或 HTTPS 请求应该被转发到哪个 Service 上，例如可以根据请求中配置的不同 Host 和 Url 路径让请求负载到不同的 Service 上。
>
> Ingress Controlle ==> 反向代理程序，它负责解析 Ingress 的反向代理规则，如果 Ingress 有增删改的变动，所有的 Ingress Controller 都会及时更新自己相应的转发规则，当 Ingress Controller 收到请求后就会根据这些规则将请求转发到对应的 Service 。

### 东西向流量（EAST-WEST traffic）

> 东西向流量则是指数据在**内部网络**中的流动，从一个区域流向另一个区域的流量。
>
> 在 kubernetes 场景中，东西向流量指：集群内服务内彼此之间的流量，主要是服务之间的 RPC 流量。

通常是通过 service 来完成，集群内的应用，都通过 Service 将请求流量发送到目标 Pod 中。

但是通常情况下集群内的 Service 使用的是 Overlay 网络（【又叫叠加网络、覆盖网络】简单理解就是把一个逻辑网络建立在一个实体网络之上），存在一定的性能损耗。所以某些场景下，我们并不希望通过 Service 来访问其他的 Pod，而是希望能直接拿到对应的 Pod IP 进行直连。

这种场景下，Service 作用可以认为只是一个 DNS name， 或者自动化管理 endpoint/endpointslices 等资源/信息的组件。这时通常很容易想到 istio/Kuma 等 Service Mesh 项目，他们都通过自己的自定义 CRD 来完成东西向流量管理。当然，因为 kubernetes 未制定任何规范，所以各个实现之间的 CRDs 互不兼容。

### Gateway API 出现的背景

2018 年 Kubernetes 社区做的一项统计，调查用户的 Ingress 资源中包含多少为了实现特定需求而增加的 annotations。从结果上可以看到只有 8% 的用户没有额外添加 annotations。同时，将近 85% 的用户期望 Ingress 可以具备良好的移植性。这样用户在不同的环境、集群中就可以更容易的采用相同配置进行处理了。而且相比于更具备表达能力而言，用户更希望的还是可以具备移植性。

所以在 2019 年圣地亚哥 KubeCon 大会上，来自 Kubernetes Networking Special Interest Group (SIG Network) 的很多人共同商量讨论了 Ingress 的现状，用户需求等相关内容，并达成一致要启动 Gateway API 规范的制定。

### Gateway API 的现状

由于在设计 Gateway API 的时候，Ingress 资源已经存在了 4 年时间，并且收到了很多用户真实环境的反馈。基于这些经验，设计出的 Gateway API 有如下特点：

![img](https://static.goodrain.com/wechat/gateway-api-indepth/1.png)

#### 1. 面向角色

运维只需要关注 **Gateway** 和 **GatewayClass**，而开发只需要关注各种 **Route**。

- `GatewayClass`: 这是由基础设施供应商提供/配置的，可用于定义一些和基础设施相关的能力；
- `Gateway`: 这个是由集群运维来管理的，可以定义一些 Gateway 自身相关的能力，例如监听哪些端口，支持哪些协议等；
- `HTTPRoute`/`TCPRoute`/`*Route` : 路由规则可以由应用开发者进行管理和发布，因为应用开发知道应用需要暴露出哪些端口，访问路径，以及需要 Gateway 配合哪些功能等。

#### 2. 更丰富的表现力

Gateway API 相比于 Ingress 具有更强的表现力，也就是说通过它原生在 spec 中定义的内容，就可以达到预期的效果，而不需要像 Ingress 资源那样添加 annotations 进行辅助。

#### 3. 可扩展

GatewayAPI 在设计的时候，预留了一些可以进行扩展的点，例如以下的配置，可以在 `spec.parametersRef` 中关联到任意的自定义资源，从而实现更为复杂的需求。



