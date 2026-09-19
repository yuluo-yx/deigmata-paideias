# 部署验收记录

最后验证时间：2026-09-19 11:19，Asia/Shanghai。

## 结果

| 检查 | 命令或方法 | 结果 |
| --- | --- | --- |
| Compose 静态校验 | `docker compose config --quiet` | 通过，退出码 0 |
| 启动与健康检查 | `docker compose up -d --wait --wait-timeout 180` | 通过，退出码 0，三个服务均 healthy |
| 容器状态 | `docker compose ps --format json` | 三个服务 running，健康状态 healthy |
| New API 页面 | `GET http://127.0.0.1:3000/` | HTTP 200 |
| New API 状态 | `GET /api/status` | HTTP 200 |
| New API 初始化状态 | `GET /api/setup` | HTTP 200，SQLite，等待首次创建管理员 |
| CPA 管理页面 | `GET http://127.0.0.1:18317/management.html` | HTTP 200 |
| 管理面板健康检查 | `GET :18317/health` | HTTP 200 |
| 管理面板采集状态 | 携带面板密钥读取 `GET :18317/status` | HTTP 200，collector running，transport http，上游 cli-proxy-api:8317 |
| 面板到 CPA 的管理连接 | 携带面板密钥读取 `GET :18317/v0/management/proxy-url` | HTTP 200，返回 host.docker.internal:7897 代理地址 |
| CPA 代理配置 | 携带 CPA 管理密钥读取 `GET :8317/v0/management/proxy-url` | HTTP 200，返回配置的 HTTP 代理 |
| CPA 模型接口 | 携带 CPA 调用密钥读取 `GET :8317/v1/models` | HTTP 200，模型列表为空，尚未添加上游账号 |
| 未认证访问 | 不携带密钥读取 CPA 模型接口及两侧管理代理接口 | 三个请求均被 401 或 403 拒绝 |
| 主机代理 | `curl -x http://127.0.0.1:7897 https://www.gstatic.com/generate_204` | HTTP 204 |
| 容器代理 | 在 CPA 容器使用 `openssl s_client -proxy host.docker.internal:7897 -connect www.gstatic.com:443 -servername www.gstatic.com -verify_return_error -quiet`，发送 generate_204 HTTP 请求 | TLS 验证成功，HTTP 204 |

## 部署调整

CPA 当前镜像基于 Debian，未包含 wget。第一次启动的健康检查失败后，已将 CPA 健康检查调整为镜像自带的 Bash TCP / HTTP 检查，要求响应状态为 HTTP 200；再次启动通过。修改前配置保存在 `backups/20260919-111811/compose.yaml`，该备份仅用于记录初始配置，不应作为可用部署版本恢复。

Docker socket 超出沙箱默认访问范围，已通过工具权限提升完成本机部署。未修改其他已有容器，未触发远端 CI/CD。

## 交付与边界

交付文件：`compose.yaml`、`.env`、`cpa/config.yaml`、`.gitignore`、`README.md`、本记录及运行数据目录。配置和数据均保留在当前工作区，镜像摘要锁定在 `.env`。

这是新部署，无需迁移。停止回滚命令为 `docker compose down`，持久化数据仍保留；完整备份与升级回滚步骤见 `README.md`。

真实模型推理、OAuth 登录和请求用量采集需要用户添加上游账号后验证；New API 管理员账号由用户首次访问时创建。未运行源码覆盖率、性能、压力、容量、混沌测试或镜像漏洞扫描，因此不对这些项目作通过声明。

## Git 发布验证

2026-09-19 增加 `.env.example`、`cpa/config.example.yaml` 和 `init-config.py`，供克隆仓库后生成独立的本机密钥。原实例的 `.env`、CPA 配置及数据保持不变。

在临时目录复制公开文件后，已验证初始化成功、CPA 两端密钥一致、配置权限为 `600`、重复初始化拒绝覆盖，以及 `docker compose config --quiet` 通过。另检查 Compose 渲染结果，确认 New API 的 `/data` 仍挂载到部署目录的 `data/new-api`。

公开文件在提交前逐一检查，未包含当前实例密钥、数据库、日志、备份或维护记录。目标仓库使用 push 触发 CI，因此提交消息包含 `[skip ci]`，遵循本任务不触发远端 CI/CD 的要求；机制依据见 [GitHub 跳过工作流说明](https://docs.github.com/en/actions/how-tos/manage-workflow-runs/skip-workflow-runs)。
