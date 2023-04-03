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
- 实现user模块 register login info
- - 学习gin-gorm-mysql ☑️
- - 设计数据表 ☑️
- - 根据数据表，利用gorm和mysql实现基本框架

## User模块

## 命名

增删改查

add delete update find










