## socket client for socket gateway


### 整体说明
基础部分
- client 基础连接的维持，数据读写的处理
- gateway 基础网关通信功能
- security 数据安全机制的处理
- 其他服务
  - 每个服务都提供 whenStarted whenPaused 接口处理