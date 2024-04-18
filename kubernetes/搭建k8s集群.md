## 一：搭建基础环境

k8s官方网站：https://kubernetes.io/zh/ ，可自行查看相关文档说明

k8s基础概念：https://zhuanlan.zhihu.com/p/621593509

k8s-master：Ubuntu–192.168.152.100
k8s-node01：Ubuntu–192.168.152.101
k8s-node02：Ubuntu–192.168.152.102

## 二：前置条件设置

全部已安装docker，未安装可根据官方文档安装：https://docs.docker.com/get-docker/

在线安装docker环境

```bash
curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun
systemctl enable docker
#启动Docker服务：
systemctl start docker
#检查是否安装成功：
docker version
```

### 1，禁止swap分区

K8s的要求，确保禁止掉swap分区，不禁止，初始化会报错。
在每个宿主机上执行：

```shell
sudo swapoff -a
#修改/etc/fstab，注释掉swap那行，持久化生效
sudo vi /etc/fstab
```

### 2，确保时区和时间正确

时区设置

```shell
sudo timedatectl set-timezone Asia/Shanghai
#同时使系统日志时间戳也立即生效
sudo systemctl restart rsyslog
```

### 3，关闭防火墙和selinux

ubuntu 查看防火墙命令，`ufw status`可查看状态，ubuntu20.04默认全部关闭，无需设置。

### 4，主机名和hosts设置（可选）

非必须，但是为了直观方便管理，建议设置。

在宿主机分别设置主机名：k8s-master，k8s-node01，k8s-node02

```shell
sudo vim /etc/hostname
```

hosts设置

```shell
sudo vim /etc/hosts
#添加内容
192.168.152.100 k8s-master
192.168.152.101 k8s-node01
192.168.152.102 k8s-node02
```

## 三：安装组件

### 1，更改docker默认驱动为systemd

为防止初始化出现一系列的错误，请检查docker和kubectl驱动是否一致，否则kubectl没法启动造成报错。版本不一样，docker有些为cgroupfs，而kubectl默认驱动为systemd，所以需要更改docker驱动。

可查看自己docker驱动命令：

```shell
sudo docker info|grep Driver
```

更改docker驱动，编辑 /etc/docker/daemon.json (没有就新建一个），添加如下启动项参数即可：

```shell
#编辑创建文件
sudo vim /etc/docker/daemon.json
#添加内容
{
    "registry-mirrors": ["https://registry.docker-cn.com",
"https://mirror.ccs.tencentyun.com",
"http://docker.mirrors.ustc.edu.cn",
"http://hub-mirror.c.163.com",
"https://dpxn2pal.mirror.aliyuncs.com"
],
   "exec-opts": ["native.cgroupdriver=systemd"]
}
```

> 上面这个配置，必须在每台电脑上都操作

重启docker

```shell
sudo systemctl restart docker.service
```

### 2，更新 apt 包索引并安装使用 Kubernetes apt 仓库所需要的包

安装软件包以允许apt通过HTTPS使用存储库，已安装软件的可以忽略

```shell
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl
```

### 3，下载公开签名秘钥、并添加k8s库

**国外** ：下载 Google Cloud 公开签名秘钥：

```shell
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add
sudo apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"
```

**国内**：可以用阿里源即可：

```shell
curl -s https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo apt-key add -
sudo tee /etc/apt/sources.list.d/kubernetes.list <<EOF
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF
```

请注意，在命令中，使用的是Ubuntu 16.04 Xenial 版本， 是可用的最新 Kubernetes 存储库。所以而非20.04 的focal。

### 4，更新 apt 包索引，安装 kubelet、kubeadm 和 kubectl，并锁定其版本

锁定版本，防止出现不兼容情况，例如，1.7.0 版本的 kubelet 可以完全兼容 1.8.0 版本的 API 服务器，反之则不可以。

- kubeadm：用来初始化集群的指令。
- kubelet：在集群中的每个节点上用来启动 Pod 和容器等。
- kubectl：用来与集群通信的命令行工具。

```shell
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
```

## 四：初始化主节点master

只需要在master上操作即可。

```shell
kubeadm init --apiserver-advertise-address=192.168.152.100 --pod-network-cidr=172.16.0.0/16
```

- apiserver-advertise-address：定义API服务器将公布其监听内容的IP地址。如果未设置，将使用默认网络接口。

- pod-network-cidr：指定Pod网络的IP地址范围。如果设置，控制平面将自动为每个节点分配CIDR。 如果您的网络插件的首选Pod网络与某些主机网络之间发生冲突，则可以使用它来定义首选的网络范围

> 这里会遇见一个坑，因为要下载K8s所需的Docker镜像文件，众所周知K8s是谷歌研发出来的，registry服务器在谷歌，国内连不上，所以这里我们要用阿里云镜像
>
> 最后添加`--image-repository=registry.aliyuncs.com/google_containers` 即可

### 1，初始化错误解决（没有报错的可以跳过这条）

#### **错误提示1：**

```shell
[kubelet-check] It seems like the kubelet isn't running or healthy.
[kubelet-check] The HTTP call equal to 'curl -sSL http://localhost:10248/healthz' failed with error: Get "http://localhost:10248/healthz": dial tcp 127.0.0.1:10248: connect: connection refused.
```

原因：kubectl没法启动，journalctl -xe查看启动错误信息。

```shell
journalctl -xe
#信息显示docker和kubectel驱动不一致
kubelet cgroup driver: \"systemd\" is different from docker cgroup driver: \"cgroupfs\""
```

解决方案：k8s建议systemd驱动，所以更改docker驱动即可，编辑 /etc/docker/daemon.json (没有就新建一个），添加如下启动项参数即可：

```shell
#编辑创建文件
sudo vim  /etc/docker/daemon.json
#添加内容
{
  "exec-opts": ["native.cgroupdriver=systemd"]
}
```

重启docker和kubectel

```shell
#重启docker
sudo systemctl restart docker.service
#重载kubectl
sudo systemctl daemon-reload  
#重启kubectl
sudo systemctl restart kubelet.service 
#查看kubectl服务状态恢复正常
sudo systemctl status kubelet.service
```

#### **错误提示2：**

```shell
error execution phase preflight: [preflight] Some fatal errors occurred:
    [ERROR FileAvailable--etc-kubernetes-manifests-kube-apiserver.yaml]: /etc/kubernetes/manifests/kube-apiserver.yaml already exists
    [ERROR FileAvailable--etc-kubernetes-manifests-kube-controller-manager.yaml]: /etc/kubernetes/manifests/kube-controller-manager.yaml already exists
    [ERROR FileAvailable--etc-kubernetes-manifests-kube-scheduler.yaml]: /etc/kubernetes/manifests/kube-scheduler.yaml already exists
    [ERROR FileAvailable--etc-kubernetes-manifests-etcd.yaml]: /etc/kubernetes/manifests/etcd.yaml already exists
[preflight] If you know what you are doing, you can make a check non-fatal with `--ignore-preflight-errors=...`
```

原因：初始化生产的文件，重新初始化，需要删除即可

```shell
rm -fr /etc/kubernetes/manifests/*
```

#### **错误提示3：**

```shell
 error execution phase preflight: [preflight] Some fatal errors occurred:
    [ERROR Port-10250]: Port 10250 is in use
```

解决方法：重置配置

```shell
sudo kubeadm reset
```

#### 错误提示4

```shell
[control-plane] Creating static Pod manifest for "kube-controller-manager"
[control-plane] Creating static Pod manifest for "kube-scheduler"
[etcd] Creating static Pod manifest for local etcd in "/etc/kubernetes/manifests"
[wait-control-plane] Waiting for the kubelet to boot up the control plane as static Pods from directory "/etc/kubernetes/manifests". This can take up to 4m0s
[kubelet-check] Initial timeout of 40s passed.

Unfortunately, an error has occurred:
        timed out waiting for the condition

This error is likely caused by:
        - The kubelet is not running
        - The kubelet is unhealthy due to a misconfiguration of the node in some way (required cgroups disabled)

If you are on a systemd-powered system, you can try to troubleshoot the error with the following commands:
        - 'systemctl status kubelet'
        - 'journalctl -xeu kubelet'

Additionally, a control plane component may have crashed or exited when started by the container runtime.
To troubleshoot, list all containers using your preferred container runtimes CLI.
Here is one example how you may list all running Kubernetes containers by using crictl:
        - 'crictl --runtime-endpoint unix:///var/run/containerd/containerd.sock ps -a | grep kube | grep -v pause'
        Once you have found the failing container, you can inspect its logs with:
        - 'crictl --runtime-endpoint unix:///var/run/containerd/containerd.sock logs CONTAINERID'
error execution phase wait-control-plane: couldn't initialize a Kubernetes cluster
To see the stack trace of this error execute with --v=5 or higher
```

解决方法：

修改/etc/containerd/config.toml文件，复制以下配置覆盖

```shell
disabled_plugins = []
imports = []
oom_score = 0
plugin_dir = ""
required_plugins = []
root = "/var/lib/containerd"
state = "/run/containerd"
temp = ""
version = 2

[cgroup]
  path = ""

[debug]
  address = ""
  format = ""
  gid = 0
  level = ""
  uid = 0

[grpc]
  address = "/run/containerd/containerd.sock"
  gid = 0
  max_recv_message_size = 16777216
  max_send_message_size = 16777216
  tcp_address = ""
  tcp_tls_ca = ""
  tcp_tls_cert = ""
  tcp_tls_key = ""
  uid = 0

[metrics]
  address = ""
  grpc_histogram = false

[plugins]

  [plugins."io.containerd.gc.v1.scheduler"]
    deletion_threshold = 0
    mutation_threshold = 100
    pause_threshold = 0.02
    schedule_delay = "0s"
    startup_delay = "100ms"

  [plugins."io.containerd.grpc.v1.cri"]
    device_ownership_from_security_context = false
    disable_apparmor = false
    disable_cgroup = false
    disable_hugetlb_controller = true
    disable_proc_mount = false
    disable_tcp_service = true
    enable_selinux = false
    enable_tls_streaming = false
    enable_unprivileged_icmp = false
    enable_unprivileged_ports = false
    ignore_image_defined_volumes = false
    max_concurrent_downloads = 3
    max_container_log_line_size = 16384
    netns_mounts_under_state_dir = false
    restrict_oom_score_adj = false
    sandbox_image = "registry.aliyuncs.com/google_containers/pause:3.9"
    selinux_category_range = 1024
    stats_collect_period = 10
    stream_idle_timeout = "4h0m0s"
    stream_server_address = "127.0.0.1"
    stream_server_port = "0"
    systemd_cgroup = false
    tolerate_missing_hugetlb_controller = true
    unset_seccomp_profile = ""

    [plugins."io.containerd.grpc.v1.cri".cni]
      bin_dir = "/opt/cni/bin"
      conf_dir = "/etc/cni/net.d"
      conf_template = ""
      ip_pref = ""
      max_conf_num = 1

    [plugins."io.containerd.grpc.v1.cri".containerd]
      default_runtime_name = "runc"
      disable_snapshot_annotations = true
      discard_unpacked_layers = false
      ignore_rdt_not_enabled_errors = false
      no_pivot = false
      snapshotter = "overlayfs"

      [plugins."io.containerd.grpc.v1.cri".containerd.default_runtime]
        base_runtime_spec = ""
        cni_conf_dir = ""
        cni_max_conf_num = 0
        container_annotations = []
        pod_annotations = []
        privileged_without_host_devices = false
        runtime_engine = ""
        runtime_path = ""
        runtime_root = ""
        runtime_type = ""

        [plugins."io.containerd.grpc.v1.cri".containerd.default_runtime.options]

      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes]

        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
          base_runtime_spec = ""
          cni_conf_dir = ""
          cni_max_conf_num = 0
          container_annotations = []
          pod_annotations = []
          privileged_without_host_devices = false
          runtime_engine = ""
          runtime_path = ""
          runtime_root = ""
          runtime_type = "io.containerd.runc.v2"

          [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
            BinaryName = ""
            CriuImagePath = ""
            CriuPath = ""
            CriuWorkPath = ""
            IoGid = 0
            IoUid = 0
            NoNewKeyring = false
            NoPivotRoot = false
            Root = ""
            ShimCgroup = ""
            SystemdCgroup = true

      [plugins."io.containerd.grpc.v1.cri".containerd.untrusted_workload_runtime]
        base_runtime_spec = ""
        cni_conf_dir = ""
        cni_max_conf_num = 0
        container_annotations = []
        pod_annotations = []
        privileged_without_host_devices = false
        runtime_engine = ""
        runtime_path = ""
        runtime_root = ""
        runtime_type = ""

        [plugins."io.containerd.grpc.v1.cri".containerd.untrusted_workload_runtime.options]

    [plugins."io.containerd.grpc.v1.cri".image_decryption]
      key_model = "node"

    [plugins."io.containerd.grpc.v1.cri".registry]
      config_path = ""

      [plugins."io.containerd.grpc.v1.cri".registry.auths]

      [plugins."io.containerd.grpc.v1.cri".registry.configs]

      [plugins."io.containerd.grpc.v1.cri".registry.headers]

      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]

    [plugins."io.containerd.grpc.v1.cri".x509_key_pair_streaming]
      tls_cert_file = ""
      tls_key_file = ""

  [plugins."io.containerd.internal.v1.opt"]
    path = "/opt/containerd"

  [plugins."io.containerd.internal.v1.restart"]
    interval = "10s"

  [plugins."io.containerd.internal.v1.tracing"]
    sampling_ratio = 1.0
    service_name = "containerd"

  [plugins."io.containerd.metadata.v1.bolt"]
    content_sharing_policy = "shared"

  [plugins."io.containerd.monitor.v1.cgroups"]
    no_prometheus = false

  [plugins."io.containerd.runtime.v1.linux"]
    no_shim = false
    runtime = "runc"
    runtime_root = ""
    shim = "containerd-shim"
    shim_debug = false

  [plugins."io.containerd.runtime.v2.task"]
    platforms = ["linux/amd64"]
    sched_core = false

  [plugins."io.containerd.service.v1.diff-service"]
    default = ["walking"]

  [plugins."io.containerd.service.v1.tasks-service"]
    rdt_config_file = ""

  [plugins."io.containerd.snapshotter.v1.aufs"]
    root_path = ""

  [plugins."io.containerd.snapshotter.v1.btrfs"]
    root_path = ""

  [plugins."io.containerd.snapshotter.v1.devmapper"]
    async_remove = false
    base_image_size = ""
    discard_blocks = false
    fs_options = ""
    fs_type = ""
    pool_name = ""
    root_path = ""

  [plugins."io.containerd.snapshotter.v1.native"]
    root_path = ""

  [plugins."io.containerd.snapshotter.v1.overlayfs"]
    root_path = ""
    upperdir_label = false

  [plugins."io.containerd.snapshotter.v1.zfs"]
    root_path = ""

  [plugins."io.containerd.tracing.processor.v1.otlp"]
    endpoint = ""
    insecure = false
    protocol = ""

[proxy_plugins]

[stream_processors]

  [stream_processors."io.containerd.ocicrypt.decoder.v1.tar"]
    accepts = ["application/vnd.oci.image.layer.v1.tar+encrypted"]
    args = ["--decryption-keys-path", "/etc/containerd/ocicrypt/keys"]
    env = ["OCICRYPT_KEYPROVIDER_CONFIG=/etc/containerd/ocicrypt/ocicrypt_keyprovider.conf"]
    path = "ctd-decoder"
    returns = "application/vnd.oci.image.layer.v1.tar"

  [stream_processors."io.containerd.ocicrypt.decoder.v1.tar.gzip"]
    accepts = ["application/vnd.oci.image.layer.v1.tar+gzip+encrypted"]
    args = ["--decryption-keys-path", "/etc/containerd/ocicrypt/keys"]
    env = ["OCICRYPT_KEYPROVIDER_CONFIG=/etc/containerd/ocicrypt/ocicrypt_keyprovider.conf"]
    path = "ctd-decoder"
    returns = "application/vnd.oci.image.layer.v1.tar+gzip"

[timeouts]
  "io.containerd.timeout.bolt.open" = "0s"
  "io.containerd.timeout.shim.cleanup" = "5s"
  "io.containerd.timeout.shim.load" = "5s"
  "io.containerd.timeout.shim.shutdown" = "3s"
  "io.containerd.timeout.task.state" = "2s"

[ttrpc]
  address = ""
  gid = 0
  uid = 0
```

- 让配置生效

  ```bash
  systemctl daemon-reload && systemctl restart containerd
  ```

  重新执行kubeadm init

#### 错误提示5

```shell
[init] Using Kubernetes version: v1.26.3
[preflight] Running pre-flight checks
error execution phase preflight: [preflight] Some fatal errors occurred:
        [ERROR CRI]: 
    container runtime is not running: output: time="2023-03-24T19:16:15+08:00" 
    level=fatal msg="validate service connection: 
    CRI v1 runtime API is not implemented for endpoint 
    \"unix:///var/run/containerd/containerd.sock\"
    : rpc error: code = Unimplemented desc = 
    unknown service runtime.v1.RuntimeService"
```

错误原因：

- **CRI（Container Runtime Interface）的远程调用接口**，这个接口**定义了容器运行时的各项核心操作**，比如：启动一个容器需要的所有参数。没有容器运行时就创建不了容器

解决方法：

查看 / 目录的 /etc/containerd/config.toml文件，这个是容器运行时的配置文件

```shell
vim /etc/containerd/config.toml
```

如果看到了这行：

```shell
disabled_plugins : ["cri"]
```

将这行用#注释或者将"cri"删除

```shell
#disabled_plugins : ["cri"]
 
disabled_plugins : []
```

 重启容器运行时

```shell
systemctl restart containerd
```

#### 错误提示6

在安装[k8s](https://so.csdn.net/so/search?q=k8s&spm=1001.2101.3001.7020)集群时，需要把docker的cgroup改成systemd，修改完后报错

```shell
Nov 15 17:25:08 node1 systemd[1]: docker.service: Main process exited, code=exited, status=1/FAILURE
Nov 15 17:25:08 node1 systemd[1]: docker.service: Failed with result 'exit-code'.
Nov 15 17:25:08 node1 systemd[1]: Failed to start Docker Application Container Engine.
Nov 15 17:25:10 node1 systemd[1]: docker.service: Scheduled restart job, restart counter is at 3.
Nov 15 17:25:10 node1 systemd[1]: Stopped Docker Application Container Engine.
Nov 15 17:25:10 node1 systemd[1]: docker.service: Start request repeated too quickly.
Nov 15 17:25:10 node1 systemd[1]: docker.service: Failed with result 'exit-code'.
Nov 15 17:25:10 node1 systemd[1]: Failed to start Docker Application Container Engine.
```

解决方法：

```shell
# 解决报错
# vi /lib/systemd/system/docker.service
# 找到并删除下面这句话，保存退出，即可解决
# --exec-opt native.cgroupdriver=cgroupfs \

[root@master ~]# systemctl daemon-reload
[root@master ~]# systemctl restart docker
```



### 2，初始化完成

无报错，最后出现以下，表示初始化完成，根据提示还需要操作。

```shell
[addons] Applied essential addon: CoreDNS
[addons] Applied essential addon: kube-proxy
 
Your Kubernetes control-plane has initialized successfully!
 
To start using your cluster, you need to run the following as a regular user:
 
  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config
 
Alternatively, if you are the root user, you can run:
 
  export KUBECONFIG=/etc/kubernetes/admin.conf
 
You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/
 
Then you can join any number of worker nodes by running the following on each as root:
 
kubeadm join 192.168.152.100:6443 --token tojkw9.v4nageqhftd7v2vc \
    --discovery-token-ca-cert-hash sha256:6a5b372144d6cc2a12f8e41853554549f1290e665381f48fdd92bbf92de7b884
```

根据提示用户是root或者普通用户操作，由于大多环境不会是root用户，我也是普通用户，所以选择普通用户操作命令：

```shell
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

> 这个过程大概率会执行成功，但是也有个小坑，就是无法在后面直接使用`kubectl get node`这个命令，而是需要使用`kubectl --kubeconfig ~/.kube/config get node`才能成功执行
>
> 改变配置方法：
>
> `nano ~/.bashrc`,在最下面添加一行export KUBECONFIG=~/.kube/config
>
> 然后执行`source ~/.bashrc`，这样就可以正常使用了 

如果是root用户，执行以下命令：

```shell
export KUBECONFIG=/etc/kubernetes/admin.conf
```

初始化完成，用最后的提示命令`kubeadm join....` 在node机器上加入集群即可。

### 3，主节点pod网络设置

主节点支持网络插件：https://kubernetes.io/zh/docs/concepts/cluster-administration/addons/
这里安装Calico网络插件：https://docs.projectcalico.org/getting-started/kubernetes/self-managed-onprem/onpremises
Calico官网提供三种安装方式，1）低于50个节点，2）高于50个节点，3）etcd datastore（官方不建议此方法）。

这里选择第一种：

```shell
curl https://docs.projectcalico.org/manifests/calico.yaml -O
kubectl apply -f calico.yaml
```

> 这里可能遇到的坑是下载https://docs.projectcalico.org/manifests/calico.yaml这个文件，文件下载到本地了，但是文件内容不对，如果不对的话，需要在浏览器打开文件，然后将文件内容粘贴在这个文件中

安装完成后，`kubectl get node` 可查看节点状态，由NotReady变成Ready则正常，需要等几分钟完成。

```shell
#未安装网络插件
ubuntu@k8s-master:~$ kubectl get node
NAME         STATUS     ROLES                  AGE   VERSION
k8s-master   NotReady   control-plane,master   80m   v1.22.3
#已安装网络插件
ubuntu@k8s-master:~$ kubectl get node
NAME         STATUS   ROLES                  AGE   VERSION
k8s-master   Ready    control-plane,master   83m   v1.22.3
```

## 五：加入node节点

### 1，node加入master节点

在所有node节点机器操作，统一已安装完成 kubelet、kubeadm 和 kubectl，用master初始化完成后最后提示命令加入，切记要用root用户。

```shell
kubeadm join 192.168.152.100:6443 --token tojkw9.v4nageqhftd7v2vc \
    --discovery-token-ca-cert-hash sha256:6a5b372144d6cc2a12f8e41853554549f1290e665381f48fdd92bbf92de7b884
```

加入成功后，提示如下：

```shell
[preflight] Running pre-flight checks
[preflight] Reading configuration from the cluster...
[preflight] FYI: You can look at this config file with 'kubectl -n kube-system get cm kubeadm-config -o yaml'
[kubelet-start] Writing kubelet configuration to file "/var/lib/kubelet/config.yaml"
[kubelet-start] Writing kubelet environment file with flags to file "/var/lib/kubelet/kubeadm-flags.env"
[kubelet-start] Starting the kubelet
[kubelet-start] Waiting for the kubelet to perform the TLS Bootstrap...
 
This node has joined the cluster:
* Certificate signing request was sent to apiserver and a response was received.
* The Kubelet was informed of the new secure connection details.
 
Run 'kubectl get nodes' on the control-plane to see this node join the cluster.
```

再次查看kubelet服务已正常启动。

### 2，需注意的坑

1：加入主节点，需要`root`用户执行词条命令，才可以加入master主节点。
node在没有加入主节点master之前，kubelet服务是没法启动的，是正常情况，会报错如下：

```shell
"Failed to load kubelet config file" err="failed to load Kubelet config file /var/lib/kubelet/config.yaml, 
error failed to read kubelet config file \"/var/lib/kubelet/config.yaml,
error: open /var/lib/kubelet/config.yaml: no such file or directory" path="/var/lib/kubelet/config.yaml"
```

原因是缺失文件，主节点master初始化 `kubeadm init`生成。
node节点是不需要初始化的，所以只需要用root用户`kubeadm join`加入master即可生成。

2：如果加入提示某些文件已存在，如：

```shell
[ERROR FileAvailable--etc-kubernetes-pki-ca.crt]: /etc/kubernetes/pki/ca.crt already exists
```

原因是加入过主节点，即使没成功加入，文件也会创建，所以需要重置节点，重新加入即可，重置命令：

```shell
kubeadm reset
```

### 3，在master查看节点

加入完成后，在master节点`kubectl get node`可查看已加入的所有节点：

```shell
ubuntu@k8s-master:~$ kubectl get node
NAME         STATUS   ROLES                  AGE     VERSION
k8s-master   Ready    control-plane,master   16h     v1.22.3
k8s-node01   Ready    <none>                 24m     v1.22.3
k8s-node02   Ready    <none>                 6m54s   v1.22.3
```

这里k8s集群创建完成。