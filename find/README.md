# 子域名树状展示系统

一个基于 Go 的子域名扫描服务，前端使用树状结构展示从根域名发散的子域名信息。每个节点包含域名、IP 与地理位置。

## 功能

- 输入根域名进行子域名扫描
- DNS 解析获取 IP 列表
- 可选的地理位置查询（默认使用 ipapi.co）
- 前端树状展示，显示域名、IP、地理位置

## 运行

```bash
go run .
```

访问浏览器：`http://localhost:8080`

## 配置

- `PORT`：服务端口，默认 `8080`
- `GEO_PROVIDER`：地理位置提供方，默认 `ipapi`，设置为 `none` 可关闭
- `GEO_API_URL`：自定义地理位置 API 地址（可选）

## 接口

- `POST /scan`

```json
{
  "domain": "example.com",
  "maxDepth": 2,
  "workers": 20
}
```

- `GET /tree/{scanId}`

## 说明

- 子域名扫描基于 `wordlist.txt` 的字典枚举。
- 扫描深度与并发可在前端输入。
- 地理位置接口可能有频率限制，建议按需配置或关闭。
