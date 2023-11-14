# mirai-api-http的Go SDK

![](https://img.shields.io/github/languages/top/CuteReimu/mirai-sdk-http "语言")
[![](https://img.shields.io/github/actions/workflow/status/CuteReimu/mirai-sdk-http/golangci-lint.yml?branch=master)](https://github.com/CuteReimu/mirai-sdk-http/actions/workflows/golangci-lint.yml "代码分析")
[![](https://img.shields.io/github/contributors/CuteReimu/mirai-sdk-http)](https://github.com/CuteReimu/mirai-sdk-http/graphs/contributors "贡献者")
[![](https://img.shields.io/github/license/CuteReimu/mirai-sdk-http)](https://github.com/CuteReimu/mirai-sdk-http/blob/master/LICENSE "许可协议")

这是针对[mirai-api-http](https://github.com/project-mirai/mirai-api-http)编写的Go SDK。

相较于直接使用mirai-core和mirai-console而言，mirai-api-http的好处是，在你更新代码时，你无需进行重新登录。

## 开始

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

    ## 就填-1
    reservedSyncId: -1
```

引入项目：

```bash
go get -u github.com/CuteReimu/mirai-sdk-http
```

关于如何使用，可以参考`examples`文件夹下的例子

## 注意事项

所有`ListenXXXXXX`函数之间都不支持并发，你可以在启动机器人的情况下调用这些函数，但是不要同时在多个协程调用这些函数。~~（不过在一般情况下确实不会有这种奇怪的需求）~~

## 进度

目前已支持的功能有：

- 消息链
  - [x] 所有消息类型
  - [x] 所有消息解析
  - [x] 所有其它客户端同步消息解析
- 事件
  - [ ] Bot自身事件
  - [ ] 好友事件
  - [ ] 群事件
  - [ ] 申请事件
  - [ ] 其它客户端事件
  - [ ] 命令事件
- 请求
  - [x] 获取插件信息
  - [x] 缓存操作
  - [ ] 获取账号信息
  - [x] 消息发送与撤回
  - [ ] 文件操作
  - [ ] 多媒体内容上传
  - [ ] 账号管理
  - [x] 群管理
  - [ ] 群公告
  - [ ] 事件处理
  - [ ] Console命令
- 其它
  - [x] 连接与认证
  - [ ] 断线重连
  - [ ] MiraiCode解析

