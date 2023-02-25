# Route2QQ-bot
QQ消息互通机器人

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
 "routes": [
  {
   "type": "kook2qq",
   "host": "1111111112375343",
   "qqgroup": "222222738",
   "hostinvite":"https://kook.top/123456"
  },
  {
   "type": "kook2qq",
   "host": "1111111111521062",
   "qqgroup": "222222174",
   "hostinvite":"https://kook.top/654321"
  },
  {
   "type": "vc2qq",
   "vcurl": "https://vocechat.test/api/bot/send_to_group/1",
   "secret": "556c5957b22756964223a332c226e6f6e6365223a227a67446",
   "qqgroup": "5543054283"
  }
 ],
  "vcport": "25535",
  "masterid": "30000000",
  "stdoutchannel": "2000000",
  "token": "your bots login token"
}
```

- `routes` 为转发路由信息。可以实现多组其他IM与 `QQ` 群之间的映射。其中`type`为映射类型。目前支持`kook2qq`(kook)和`vc2qq`(vocechat)
  - 在所有类型中
    - `qqgroup` 为 `QQ` 群号
    - `hostinvite` 为 其他IM对应的邀请链接或者访问链接，当无法成功转发消息至 `QQ` 时，将会建议至此链接查看消息。
  - 在 `kook2qq` 类型中，
    - `host` 为 `KOOK` 频道的ID
  - 在 `vc2qq` 类型中，
    - `vcurl` 为 `vocechat` 推送群消息的接口
    - `secret` 为 `vocechat` 推送群消息的 `api-key`
- `vcport` 为 接收 `vocechat webhook` 消息推送的端口，将会监听此端口
- `stdoutchannel` 为 `KOOK` 上机器人输出调试信息的频道，不使用可以留空。
- `masterid` 为你的 `KOOK` ID
- `token` 为你的机器人的登录 `token`
