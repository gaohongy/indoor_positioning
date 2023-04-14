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

## service
> 目前实现不考虑性能，先完成实现，故service目前是删除状态

一般在 handler 中主要做解析参数、返回数据操作，简单的逻辑也可以在 handler 中做，像新增用户、删除用户、更新用户，代码量不大，所以也可以放在 handler 中。
有些代码量很大的逻辑就不适合放在 handler 中，因为这样会导致 handler 逻辑不是很清晰，这时候实际处理的部分通常放在 service 包中。比如 LisReferencepoint() 函数

## 说明
handler/user/user.go 和 model/user.go 区别在于，前者是和user-api相关的请求响应结构，后者是和user相关的数据库操作

需要明白的是，api接口所接收的数据格式是不同于数据库读写的数据格式的。api接口包含哪些数据是由业务逻辑决定的，数据库读写所需数据结构是由数据库表的定义决定的

// TODO 管理员用户和普通用户均可切换场所，但管理员切换场所后仅具有普通用户权限，这样写的话，user需要加字段。
// 先写为仅普通用户可以修改place_id

## 路径点筛选程序逻辑
> 以下时间格式转换的流程未必正确，大致是这样

1. 前端获取到的是Date()生成的日期，通过toISOString()转为ISO8601格式(go语言的time.Time似乎就是这个格式)，通过encodeURIComponent()解决一些特殊字符无法在url中传输的问题
2. 后端拿到数据后，通过QueryUnescape将上述转换后的时间字符串还原为原字符串
3. 通过time.parse()将字符串解析为time.Time格式
4. 在查询语句中使用time.Time时间格式进行查询
难点在于时间筛选条件和用户筛选条件并非总是全部都有，需要考虑请求参数为空时的程序运行逻辑
![](https://img2023.cnblogs.com/blog/1898659/202304/1898659-20230411223201651-424563567.png)