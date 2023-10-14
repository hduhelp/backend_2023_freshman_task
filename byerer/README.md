# 2.0
增加鉴权功能

# 2.1
> gorm 默认设置 ID (todoid) 为主键，不同的用户产生相同的id

solution: 
- todoID = userID + time

> 关于越权，由于所有todo放在同一个库中，可能产生使用自己的 token ,删除别人的 todo 
> 的问题

solution：
- 解析 userid 与 todoid 是否适配
- 多条件查询数据库

似乎后者会降低查询效率？   
实现前者

# 3.0
添加前端页面 目前仅实现部分接口


# 等待升级

- 配置yaml
- 日志
- 邮件提醒
- 数据库查询效率
- docker

# 感想
docker太妙了  
js好复杂