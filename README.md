# mirai-api-http的Go SDK

![](https://img.shields.io/github/languages/top/CuteReimu/mirai-sdk-http "语言")
[![](https://img.shields.io/github/actions/workflow/status/CuteReimu/mirai-sdk-http/golangci-lint.yml?branch=master)](https://github.com/CuteReimu/mirai-sdk-http/actions/workflows/golangci-lint.yml "代码分析")
[![](https://img.shields.io/github/contributors/CuteReimu/mirai-sdk-http)](https://github.com/CuteReimu/mirai-sdk-http/graphs/contributors "贡献者")
[![](https://img.shields.io/github/license/CuteReimu/mirai-sdk-http)](https://github.com/CuteReimu/mirai-sdk-http/blob/master/LICENSE "许可协议")

这是针对[mirai-api-http](https://github.com/project-mirai/mirai-api-http)编写的Go SDK。

在使用本项目之前，你应该知道如何使用[mirai](https://github.com/mamoe/mirai)进行登录，并安装[mirai-api-http](https://github.com/project-mirai/mirai-api-http)插件。

请多参阅mirai-api-http的[文档](https://docs.mirai.mamoe.net/mirai-api-http/api/API.html)

本项目使用ws接口，因此你需要修改mirai的配置文件`config/net.mamoe.mirai-api-http/setting.yml`，开启ws监听。

```yaml
adapters:
  - ws
verifyKey: ABCDEFGHIJK
adapterSettings:
  ws:
    ## websocket server 监听的本地地址
    ## 一般为 localhost 即可, 如果多网卡等情况，自定设置
    host: localhost

    ## websocket server 监听的端口
    ## 与 http server 可以重复, 由于协议与路径不同, 不会产生冲突
    port: 8080

    ## 填-1
    reservedSyncId: -1
```

引入项目：

```bash
go get -u github.com/CuteReimu/mirai-sdk-http
```

举例：

```go
package main

import (
    "github.com/CuteReimu/mirai-sdk-http"
    "github.com/CuteReimu/mirai-sdk-http/message"
)

func main() {
    b, _ := miraihttp.Connect("localhost", 8080, miraihttp.WsChannelAll, "ABCDEFGHIJK", 123456789)
    b.WriteMessage(&message.SendGroupMessage{
        Target: 987654321,
        MessageChain: message.Chain{
            (&message.Plain{Text: "Mirai牛逼"}).toMap(),
        },
    })
}
```
