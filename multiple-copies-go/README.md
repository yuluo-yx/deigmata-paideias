# 高可用多副本实现 Go

用 Go 基于 file 实现一个简单的分布式多副本 Demo。可以将文件换为 mysql redis etcd 等等

## 流程

1. 每个副本启动后，通过 `flock(LOCK_EX|LOCK_NB)` 非阻塞尝试获取文件排他锁
2. 获取成功 → 成为 Leader，写入 holderID + 时间戳到锁文件，定时 Heartbeat 续期
3. 获取失败 → 成为 Follower，读取锁文件获取当前 Leader 和上次 Heartbeat 时间
4. Follower 检查 lease 是否过期（超过 `--lease-duration` 无 Heartbeat）
   - 未过期：继续等待，显示剩余时间
   - 已过期：尝试抢锁，成为新 Leader
5. Leader 正常退出时释放锁；crash 时内核自动释放 flock，Follower 通过 lease 超时检测接管

## 运行

```bash
go build -o election .
./election -id=replica-1 -lease-duration=5s
./election -id=replica-2 -lease-duration=5s
./election -id=replica-3 -lease-duration=5s
```

## 参考

- https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/client-go/tools/leaderelection/resourcelock/leaselock.go
- https://github.com/kubernetes/client-go/blob/master/tools/leaderelection/resourcelock/interface.go
