# KOOK2QQ-bot
KOOK与QQ消息互通机器人

## Usage

初次使用 `go test` 生成一个设备描述文件 `device.json`

### 0. 配置 application.yaml

编辑 `application.yaml` 如下：

```yaml
bot:
  # 账号
  account: 123456
  # 密码
  password: your_password
  loginmethod: common/qrcode
```

设置你使用的qq号和密码，登录方式默认为`common`，使用账号密码登录，一般无法直接登录，会自动尝试使用`qrcode`方式登录，将会在终端打印出二维码，使用手机QQ扫描即可登录。

成功登陆后将会在本地保存`session`，下次登录时会优先使用`session`登录。

所以当首次由于QQ风控无法在服务器上登录时，可以先在本地成功登陆后，将`session`与`device`复制到服务器端使用。

### 1. 配置 config.json

```json
{
 "kook2qq": {
  "1111111112375343":"222222738",
  "1111111111521062":"222222174"
 },
 "kookinvite": {
  "1111111112375343": "https://kook.top/123456",
  "1111111111521062" : "https://kook.top/654321"
 },
  "masterid": "30000000",
  "stdoutchannel": "2000000",
  "token": "your bots login token"
}
```

- `kook2qq` 为 `KOOK->QQ` 的转发数组。可以实现多组 `KOOK` 频道与 `QQ` 群之间的映射。其中key为KOOK**频道**的ID，value为QQ群号。
- `kookinvite` 为 `KOOK` 的对应频道的邀请链接。当无法成功转发消息至 `QQ` 时，将会建议至 `KOOK` 查看消息。
- `stdoutchannel` 为 `KOOK` 上机器人输出调试信息的频道，不使用可以留空。
- `masterid` 为你的 `KOOK` ID
- `token` 为你的机器人的登录 `token`

