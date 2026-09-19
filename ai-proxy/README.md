# AI Proxy 本机部署

使用一份 `compose.yaml` 启动 CLIProxyAPI 网关、CPA Manager Plus 完整管理面板和 New API。适用于当前 macOS / OrbStack 环境，所有发布端口仅监听本机。

## 服务入口

| 服务 | 地址 | 登录与用途 |
| --- | --- | --- |
| CPA Manager Plus | <http://localhost:18317/management.html> | 使用 `.env` 中的 `CPA_MANAGER_ADMIN_KEY` 登录 |
| New API | <http://localhost:3000> | 首次访问按向导创建管理员账号 |
| CPA API | <http://localhost:8317/v1> | 使用 `.env` 中的 `CPA_API_KEY` 调用 |

管理面板已通过 `http://cli-proxy-api:8317` 连接 CPA，无需再次填写上游地址。CPA 管理密钥 `CPA_MANAGEMENT_KEY` 与面板登录密钥、模型调用密钥分别生成。CPA 自带的管理页面已关闭，请使用 18317 端口的完整面板。

## 启动与停止

从仓库下载后，先进入 `ai-proxy` 目录，使用 Python 3 初始化本机配置：

```bash
python3 init-config.py
```

脚本使用标准库生成随机密钥，并从模板创建 `.env` 和 `cpa/config.yaml`，权限为 `600`。如果任一配置文件已存在，脚本会退出，不覆盖原有配置。已部署的实例跳过初始化。

在本目录执行：

```bash
docker compose config --quiet
docker compose up -d --wait --wait-timeout 180
docker compose ps
docker compose logs --tail 100
```

正常运行时三个服务均显示 `healthy`。停止并移除本套容器与网络：

```bash
docker compose down
```

数据保留在当前目录，重新执行启动命令即可恢复。服务配置了 `restart: unless-stopped`。

## 本机代理

`cpa/config.yaml` 中的 `proxy-url` 已设为 `http://host.docker.internal:7897`。CPA 和管理面板同时配置 HTTP / HTTPS 代理环境变量，容器之间的连接通过 `NO_PROXY` 直连。

请保持本机 7897 代理运行。容器中的 `127.0.0.1` 指向容器自身，因此代理使用 OrbStack 提供的 `host.docker.internal`。这套配置已针对当前主机验证；迁移到 Linux 时，需要重新配置宿主机地址并验证代理的监听地址。

New API 未设置全局出站代理。如果将 CPA 作为渠道，请使用下文的容器内地址，其上游请求由 CPA 经本机代理发送。

## 接入模型

1. 打开 CPA 管理面板，登录后添加服务商密钥或导入授权账号。
2. 打开 New API，完成管理员初始化。
3. 新建 OpenAI 兼容渠道，代理地址填 `http://cli-proxy-api:8317`，渠道密钥填 `.env` 中的 `CPA_API_KEY`，选择 CPA 实际提供的模型。
4. 在 New API 配置模型价格、用户额度和访问令牌，并执行渠道测试。

New API 渠道地址不要填写 `localhost` 或管理面板的 18317 端口。OpenAI 兼容渠道的基础地址无需追加 `/v1`；调用 New API 的客户端 Base URL 则使用 `http://localhost:3000/v1`，密钥使用 New API 签发的令牌。

CPA 的 8085、1455、54545、51121、11451 端口按上游 Compose 保留，用于不同服务商的 OAuth 回调，仅绑定本机。尚未导入上游账号时，CPA 模型列表为空属于正常情况。

## 文件与持久化

| 路径 | 内容 |
| --- | --- |
| `compose.yaml` | 三个服务、健康检查、端口和数据挂载 |
| `init-config.py` | 首次部署时生成私有配置和随机密钥 |
| `.env.example`、`cpa/config.example.yaml` | 可提交 Git 的无密钥配置模板 |
| `.env` | 镜像摘要、随机登录密钥、CPA 调用密钥及 New API 会话密钥 |
| `cpa/config.yaml` | CPA 配置、代理地址和访问密钥 |
| `data/cpa/` | 授权文件、日志与插件 |
| `data/cpa-manager-plus/` | 管理面板 SQLite 数据库及 `data.key` |
| `data/new-api/` | New API SQLite 数据库 |
| `data/new-api-logs/` | New API 日志 |
| `backups/` | 修改配置前的备份 |

`.env` 与 CPA 配置创建时权限为 `600`，并与运行数据一起加入 `.gitignore`。请在本机编辑器中查看密钥，避免将完整配置或日志公开。CPA 启动后会将管理密钥转换为哈希；`.env` 保留供管理面板连接使用的原始值。

Git 仓库只保存部署文件、配置模板和说明；实际密钥、数据库、运行日志、备份和本机维护记录保留在本机。

New API 使用官方支持的 SQLite 单实例模式，无需另启 PostgreSQL 和 Redis。镜像通过 SHA-256 摘要锁定：CPA `v7.3.8`、CPA Manager Plus `v1.13.0`、New API `v1.0.0-rc.38`；New API 所选版本是上游当时 latest release 指向的候选版本。

## 备份与回滚

首次部署无需迁移。备份时先停止服务，保证 SQLite 文件一致，再将配置、密钥和数据一起归档：

```bash
docker compose stop
backup_dir="backups/$(date +%Y%m%d-%H%M%S)"
mkdir -p "$backup_dir"
tar -czf "$backup_dir/deployment.tar.gz" compose.yaml .env cpa data
docker compose start
```

升级前保存完整备份，然后修改 `.env` 中的镜像摘要并重新启动。需要回滚时先停止服务，保留故障状态副本，再从同一份备份恢复配置和整个数据目录，最后执行启动命令。管理面板数据库必须与 `data.key` 一起恢复。

## 来源与验证范围

资料核对日期：2026-09-19。检索范围为用户指定的两个 GitHub 项目、CPA 官方配置及正式发布信息，重点核对镜像、端口、持久化、代理和管理连接变量：

- [CPA Manager Plus 部署说明](https://github.com/seakee/CPA-Manager-Plus)：确认完整面板需要独立 CPA 网关。
- [Manager Server 配置说明](https://seakee.github.io/CPA-Manager-Plus/docs/en/operations/manager-server.html)：确认管理密钥、上游地址及 HTTP 采集配置。
- [CLIProxyAPI 配置示例](https://github.com/router-for-me/CLIProxyAPI/blob/main/config.example.yaml)及 [Compose](https://github.com/router-for-me/CLIProxyAPI/blob/main/docker-compose.yml)：确认代理、挂载与 OAuth 回调端口。
- [New API 官方说明](https://github.com/QuantumNous/new-api/blob/main/README.zh_CN.md)：确认 SQLite 本地部署及首次初始化流程。

验收以 Compose 静态校验、容器健康状态、HTTP 入口、管理认证与容器代理连通性为准。本任务未修改应用源码，未运行上游单元测试、覆盖率、性能、压力、容量、混沌测试或镜像 CVE 扫描。没有上游账号，未验证真实模型推理、OAuth 登录和实际用量采集；服务启动成功不代表已经配置可调用模型。具体实测结果见 `verification.md`。
