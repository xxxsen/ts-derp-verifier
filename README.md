# ts-derp-verifier

基于 Tailscale 设备列表的 DERP 客户端鉴权服务。定时拉取 Tailscale 设备信息，缓存授权的 node key，并提供 HTTP 接口给 DERP 回调鉴权。

## 功能
- 固定间隔从 Tailscale API 拉取设备列表。
- 依据已授权设备的 `NodePublic` 进行鉴权。
- 使用 JSON 配置文件，提供单一 HTTP 接口。

## 配置
创建 `config.json`：
```jsonc
{
  "tailnet": "your-tailnet", //默认是你的邮箱地址
  "client_id": "your-client-id", //oauth clientid
  "client_secret": "your-client-secret", //oauth secret
  "listen": ":8080", //监听地址
  "refresh_interval": 600, //单位为秒
  "log": {
    "file": "/tmp/derp.log",
    "level": "info",
    "file_count": 5,
    "file_size": 10240000,
    "keep_days": 7,
    "console": true
  }
}
```

说明：
- `client_id`/`client_secret` 为 OAuth Client 凭据, 凭据需要至少包含 `devices:core:read` 权限。
- `refresh_interval` 单位为秒。
- `listen` 为空时默认 `:8080`。

**申请OAuth凭据的路径如下:**

tailscale->Settings->Trust credentials->+Credential->Devices(Core:Read)

## 运行

**使用docker**

```yaml 
services:
  ts-derp-verifier:
    image: xxxsen/ts-derp-verifier:latest
    container_name: ts-derp-verifier
    restart: always
    expose:
      - 8080
    volumes:
      - ./config:/config # 挂在配置目录
    command: --config=/config/config.json # 指定配置文件启动
```

**直接二进制启动**

```shell
# go 安装, 自己准备go环境, 需要go 1.24及以上
go install github.com/xxxsen/ts-derp-verifier/cmd/ts-derp-verifier@latest

# 本地运行
ts-derp-verifier --config=./config.json
```

## API

默认路径为 /derp/verify, 自己看情况在nginx侧进行调整吧。

服务部署成功后, 后续derp侧使用, 可以使用下面的命令: 

```shell 
# --verify-client-url 指定部署当前工具server的host+path
# --verify-client-url-fail-open 当校验地址请求失败的时候的行为, 默认为放行, 这里false为拒绝
derper --verify-client-url="https://yourhost/derp/verify" --verify-client-url-fail-open=false 
```
