# indoor_positioning

## 项目架构
项目架构大致参考[The Clean Code Blog](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)([中文译](https://learnku.com/go/t/43569))所提及的架构设计

![架构图](https://cdn.learnku.com/uploads/images/202004/21/54739/cHZbDZpxWt.png!large)


![架构模式](https://picx.zhimg.com/80/08415618172ea7a3d2b916ab0c555346_720w.webp?source=1940ef5c)

## config
1. 配置文件
2. 配置配置文件程序

> 关于database初始化代码的放置位置问题，起初我一直觉得应当放在config中才合理，直到需要写user类型的数据库操作接口时，发现需要引用数据库初始化代码中的DB对象，此时如果config.DB显然不太合适，因此还是放置到model中了

## router
1. 路由设置

## handler
handler <=> Delivery
仅负责请求参数解析，具体业务逻辑由service实现

## service
service <=> Usecase
负责完成具体业务逻辑

## model
model <=> repository
持久层，负责数据库的CURD,同时定义各种类型

## pkg
> 工具包

- auth: 身份认证工具
- error：错误码工具

## 说明
handler/user/user.go 和 model/user.go 区别在于，前者是和user-api相关的请求响应结构，后者是和user相关的数据库操作