# ts-derp-verifier

基于 Tailscale 设备列表的 DERP 客户端鉴权服务。它会定时拉取 Tailscale 设备信息，缓存授权的 node key，并提供 HTTP 接口给 DERP 回调鉴权。

## 功能
- 固定间隔从 Tailscale API 拉取设备列表。
- 依据已授权设备的 `NodePublic` 进行鉴权。
- 使用 JSON 配置文件，提供单一 HTTP 接口。

## 配置
创建 `config.json`：
```json
{
  "tailnet": "your-tailnet", //这个值是你的tailscale的邮箱
  "api_key": "tskey-...",    //需要自己申请
  "listen": ":8080", 
  "refresh_interval": 600,
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
- `refresh_interval` 单位为秒。
- `listen` 为空时默认 `:8080`。

## 运行
```bash
go run ./cmd -config ./config.json
```

## API
POST `/derp/verify`

请求/响应为 JSON，字段名首字母大写：
```json
{"NodePublic":"nodekey:1ce861bc476e736324998259467c0504febb26dbc7bec6d3dc809ca017bf937b","Source":"11.22.33.44"}
```

响应：
```json
{
  "Allow": true
}
```