# indoor_positioning

## 项目架构
项目架构大致参考[The Clean Code Blog](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)([中文译](https://learnku.com/go/t/43569))所提及的架构设计

![架构图](https://cdn.learnku.com/uploads/images/202004/21/54739/cHZbDZpxWt.png!large)


![架构模式](https://picx.zhimg.com/80/08415618172ea7a3d2b916ab0c555346_720w.webp?source=1940ef5c)

## config
1. 配置文件
2. 配置配置文件程序

## router
1. 路由设置

## handler
handler <=> Delivery
仅负责请求参数解析，具体业务逻辑由service实现

## service
service <=> Usecase
负责完成具体业务逻辑

## dao
dao <=> repository
持久层，负责数据库的CURD

## pkg
> 工具包

- error：错误码工具