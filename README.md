# zinx
源码在https://github.com/aceld/zinx  对zinx 1.0做了基本实现
- server模块：服务器基本方法与配置
- connection模块：配置TCP连接基本属性与方法
- request模块：将连接与消息绑定
- message模块：封装基本消息体
- datapack模块：工具模块，用于解决TCP数据包的粘包
- router模块：根据消息类型处理不同业务
- msghandle模块：实习工作池和任务队列
- connmanager模块：用于服务器的连接管理
# MMO_game
设定proto协议，用于client与服务器之间的通信
- AOI兴趣点
- proto协议
- 玩家上线
- 世界聊天
- 位置广播
- 玩家下线
  具体流程图可参考zinx.docx
