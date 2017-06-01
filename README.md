# magpie
HiMagpie 网关，负责 TCP 长连接建立、心跳、消费消息队列并推送等。

### 协议
使用 pb 作为与客户端数据交互的协议。

### Session 和心跳
1. 定时心跳
2. 过期关闭连接
3. 首次建立连接进行校验

### 消息推送
1. 消息有第三放 Web Server 通过 HTTP API 发送到 Housekeeper，由 Housekeeper 校验 Web Server 和客户端（ClientId）等信息。  
2. Housekeeper 将待推送消息放到 ClientId Session 所在及机器的消费队列。  
3. mappie 进行消息消费和 ACK 确认、失败重传、过期判断等。
