## 项目架构

经典的 客户端 - controller - service - dao - 数据库

客户端和controller层采用RESTFul API的形式进行交互

功能分成几个模块

- user
- video
- favorite
- comment
- relation
- message

controller层、service层、dao层分别对应这几个模块

controller层使用JWT鉴权、FTP服务器和FFMpeg实现截图

service层使用Redis保存热键，使用RabbitMQ

数据库使用MySQL数据库

## 数据库表设计

按模块设计，应该有6个数据表

每个数据表都有一个主键id

整体设计思路参照抖音api

### 用户 user

- id
- name
- password

### 视频 video

- id
- user_id (-> user.id)
- publish_time
- play_url (ftp server)
- cover_url (ftp server)
- title

### 点赞 favorite

- id
- from_user_id (-> user.id)
- to_user_id (-> user.id)

### 评论 comment

- id
- from_user_id (-> user.id)
- to_video_id (-> video.id)
- content
- create_date

### 关系 relation

- id 
- from_user_id (-> user.id)
- to_user_id (-> user.id)
 
### 消息 message

- id
- from_user_id (-> user.id)
- to_user_id (-> user.id)
- content
- create_time

### 表关系



## 进度

- 写路由 ☑️
- 规范一下命名
- 重新设计数据库，比如索引等
- 实现token
- - 生成token ☑️
- - 验证token ☑️
- 实现user模块 register login info
- - 实现register 
- - - 实现各层逻辑 ☑️
- - - 测试各层代码 ☑️
- - 实现login
- - - 实现各层逻辑 ☑️
- - - 测试各层代码 ☑️
- - 实现info
- - - 实现各层逻辑 ☑️
- - - 测试各层逻辑 （待完成）
- 实现video模块，feed、publish_action、publish_list
- - publish_action
- - - 实现逻辑 ing
- - - 测试逻辑
- 学ftp和nginx ☑️
- 学ffmpeg ☑️
- 学长连接
- 学通道和协程 （并发此处不实现）
- 学ssh调用
- 学习gin-gorm-mysql ☑️
- 设计数据表 ☑️
- 连接数据库 initDao() ☑️
- 根据数据表，利用gorm和mysql实现基本框架(从下层到上层)
- - 在mysql创建数据表 ☑️
- - 在dao层中定义结构体和写相关函数接口

## User模块

## 命名

增删改查

add delete update find count

## bug解决

Q: jwt生成token时，出现错误key is of invalid type

A: 把SignedString函数参数从string类型换成byte数组


---

Q: 测试ftp文件上传和截图时出现如下错误

发现有空指针引用问题，一个是var c *ftp.ServerConn，在后面这个变量不能用短变量声明

一个是os.Open()参数中相对路径错误

还有一个时assert.True(t, err == nil, err.Error())中err如果为空，就不能调用err.Error(),
所以在msgAndArgs参数中不能写err.Error()

用ssh截图时出现Screenshot()返回错误，错误为Process exited with status 127

A:一个session只能run一条命令，所以每次都要newSession()，发送心跳也是如此

session最好是局部变量，用完及销

最最傻逼的问题就是openssh可以用ffmpeg命令，但是golang的ssh只能用/usr/local/ffmpeg/bin/ffmpeg代替ffmpeg命令
否则找不到命令，就算设置了环境变量



## 问题

数据库中唯一索引和唯一约束

数据库应不应该用外键

https://blog.csdn.net/wch020928/article/details/126714294

是否要使用gorm自动建表

如何正确使用go的错误处理，什么场合，什么处理方式？

omitempty

当用户的请求有问题时，应该返回404，还是返回200，然后在response中说用户请求有问题

## 项目问题

上传的视频中title不能有空格

## 参考资料

gorm标签

https://blog.csdn.net/abc54250/article/details/129233456

gorm crud

https://blog.csdn.net/weixin_45604257/article/details/105139862

gin postform和query的区别

https://gin-gonic.com/zh-cn/docs/examples/query-and-post-form/

go与ftp

https://github.com/jlaffaye/ftp










