# APISIX 体验指南

> 所有的 sh 脚本通过 git bash 执行。
>
> 出现错误仔细核对文档。
>
> 导入 postman 测试脚本。
>
> github 地址：https://github.com/yuluo-yx/apisix-deploy

## 使用 docker 安装 apisix

确保本地安装 Docker 和 Docker-compose 如未安装参开以下文档安装：
Docker：https://docs.docker.com/engine/install/centos/
Docker-Compose：https://docs.docker.com/compose/install/

Clone 并修改配置连接到服务器（区别本地），新建 apisix-3.7 目录

```shell
mkdir apisix-3.7 && cd apisix-3.7
git init
git clone https://github.com/yuluo-yx/apisix-deploy
```

### 配置修改

这里的配置必须要修改，确保服务安全

#### Dashboard

```shell
cd /example/dashboard/
vim conf.yaml

# 修改如下配置

 35   expire_time: 3600     # jwt token expire time, in second
 36   users:                # yamllint enable rule:comments-indentation
 37     - username: admin   # username and password for login `manager api`
 38       password: deamgodeamgo@666

 **用户名和密码需要修改，密码使用强类型密码**
```

#### Apisix-server

```shell
# 修改 apisix admin api 的验证 key

 32     admin_key:
 33       - name: "admin"
 34         key: 054f7cf07e344346cd3f287985e76a21
 35         role: admin                 # admin: manage all configuration data
 36 

 **必须修改，用于 curl 下发配置时使用**
```

#### 启动

启动之前，先确定 linux 系统架构，执行以下命令，选择不同的版本执行！ `dpkg --print-architecture`

```shell
# 运行启动脚本（以 amd 架构为例启动）
./start.sh

# 关闭
./stop.sh

# 测试服务是否启动成功, test.sh 中的 key 需要修改为 apisix-server 中配置的 key！
./test.sh

# 出现以下结果部署成功
{"value":{"pass_host":"pass","nodes":{"httpbin.org:80":1},"update_time":1701241286,"scheme":"http","create_time":1701241286,"hash_on":"vars","id":"1","type":"roundrobin"},"key":"/apisix/upstreams/1"}
```

#### 访问测试
- 服务器开放 9000 端口，本地机器访问 http://ip:9000 输入在 dashboard 中配置的用户名和密码访问控制面板！

- 如果是本地，浏览器访问：http://127.0.0.1:9000 即可查看 dashboard 页面。

## 运行 api 服务

打开 IDEA，运行 spring boot 项目之后，进入 postman 中打开 `api 原生 ` 文件夹访问接口进行测试！

## apisix 发布接口

```shell
$ cat publish-router.sh
curl "http://127.0.0.1:9180/apisix/admin/routes/1" -H "X-API-KEY: 054f7cf07e344346cd3f287985e76a21" -X PUT -d '{"methods": ["GET"],"uri": "/api/test","upstream": {"type": "roundrobin","nodes": {"127.0.0.1:8080": 1}}}'

# 发布成功如下所示
$ sh publish-router.sh
{"key":"/apisix/routes/1","value":{"methods":["GET"],"upstream":{"type":"roundrobin","pass_host":"pass","nodes":{"127.0.0.1:8080":1},"hash_on":"vars","scheme":"http"},"id":"1","uri":"/api/test","update_time":1701359555,"priority":0,"status":1,"create_time":1701359555}}
```

> **noded 的 ip 必须是 192 开头，127 会报 502**
>
> **etcd 的地址必须配置为 ip 地址，127 会报 502**

### 测试

postman 请求 `apisix-test` `apisix-user-getone` 接口进行测试！

## apisix 保护接口

> 给 test 接口加上 限流插件 防护
>
> 规则：时间窗口为 60 s，在 60 s 之内只能被访问两次，超过直接返回 503 错误码。

```shell
$ cat protected-test.sh
curl -i http://127.0.0.1:9180/apisix/admin/routes/1 -H 'X-API-KEY: 054f7cf07e344346cd3f287985e76a21' -X PUT -d '{"uri": "/api/test","plugins": {"limit-count": {"count": 2, "time_window": 60,"rejected_code": 503,"key_type": "var","key": "remote_addr"}},"upstream": {"type": "roundrobin","nodes": {"192.168.2.27:8080": 1}}}'

$ sh protected-test.sh
{"key":"/apisix/routes/1","value":{"plugins":{"limit-count":{"key_type":"var","count":2,"allow_degradation":false,"show_limit_quota_header":true,"time_window":60,"rejected_code":503,"key":"remote_addr","policy":"local"}},"update_time":1701395836,"uri":"/api/test","create_time":1701359555,"upstream":{"pass_host":"pass","nodes":{"192.168.2.27:8080":1},"hash_on":"vars","scheme":"http","type":"roundrobin"},"priority":0,"id":"1","status":1}}
```

访问 postman `保护 test` 进行测试。

## 运行前端服务

### 本地运行

```shell
cd /app/front/apisix-test-front
yarn
yarn dev --host
```

### nginx 部署前端

```shell
cd /app/front/docker
docker-compose up
```

## apisix 发布前端路由

### 本地运行路由发布

```shell
curl "http://127.0.0.1:9180/apisix/admin/routes/1" -H "X-API-KEY: 054f7cf07e344346cd3f287985e76a21" -X PUT -d '{"methods": ["GET"],"uri": "/*","upstream": {"type": "roundrobin","nodes": {"192.168.2.27:5173": 1}}}'

$ sh publish-front-router.sh
{"key":"/apisix/routes/1","value":{"priority":0,"update_time":1701397737,"uri":"/*","status":1,"create_time":1701359555,"upstream":{"pass_host":"pass","nodes":{"192.168.2.27:5173":1},"hash_on":"vars","scheme":"http","type":"roundrobin"},"id":"1","methods":["GET"]}}
```

### nginx 路由发布

```shell
curl "http://127.0.0.1:9180/apisix/admin/routes/1" -H "X-API-KEY: 054f7cf07e344346cd3f287985e76a21" -X PUT -d '{"methods": ["GET"],"uri": "/*","upstream": {"type": "roundrobin","nodes": {"192.168.2.27": 1}}}'

$ sh publish-front-router.sh
{"key":"/apisix/routes/1","value":{"priority":0,"update_time":1701397737,"uri":"/*","status":1,"create_time":1701359555,"upstream":{"pass_host":"pass","nodes":{"192.168.2.27":1},"hash_on":"vars","scheme":"http","type":"roundrobin"},"id":"1","methods":["GET"]}}
```

## 请求测试

- 本地：浏览器访问 http://127.0.0.1:9080/ 查看
