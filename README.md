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

### 1. 配置 config.json

```json
{
  "kookchannel": "1000000",
  "masterid": "30000000",
  "qqgroup": "4000000",
  "stdoutchannel": "2000000",
  "token": "your bots login token"
}
```

`kookchannel` 为 `KOOK` 上收发 `QQ` 消息的频道。

`stdoutchannel` 为 `KOOK` 上机器人输出调试信息的频道，不使用可以留空。

`qqgroup` 为 `QQ` 上收发 `KOOK` 消息的群号。

`masterid` 为你的 `KOOK` ID

`token` 为你的机器人的登录 `token`

