# ToDoList for hduhelp freshman task

## Detail

本项目使用 Go 语言实现了一个基础的 TodoList 后端。

### Features

- 基本的 Todo 增删改查；
- 基础的用户功能；
- TodoList 的临期提醒（通过邮箱，目前禁支持 24h 时提醒）。

## Config

配置文件应当放置于执行目录下的 `config.toml` 中。

- `database_file` 数据库文件地址
- `email.address` 邮箱地址
- `email.smtp_user` SMTP 用户名
- `email.smtp_address` SMTP 服务器
- `email.smtp_port` SMTP 端口
- `email.password` 邮箱密码

## Api

API 的请求和返回均应使用 json 格式。

### Utils

#### Auth

```json
{
    "username": string
    "password": string
}
```

- username: 用户名；
- password：用户经 bcrypt hash 后的密码。

#### User 

```json
{
    "username": string
    "email": string
    "password": string
}
```

- username: 用户名；
- email: 用户邮箱；
- password：用户密码。

### Todo

```json
{
    "title": string
    "location": string
    "duedate": int
    "userid": int
    "done": bool
}
```

- title：标题；
- location：地址；
- duedate：截止时间的时间戳；
- userid：所属用户 ID；
- done： 是否完成。

### User & Token

POST /user

```
{
    "user": User
}
```

创建一个新用户，返回该用户的 ID 等信息。


PUT /user/:id

```
{
    "token": string
    "user": User
}
```

认证 Token 属于 ID 为 `id` 的用户 ，并更新该用户的用户信息。

POST /Token

```
{
    "auth": Auth
}
```

认证，并返回一个新的属于认证用户的 Token。

Token 会在最后一次使用后 24h 后被删除。

### Todo

> 你能够更新/删除/创建/访问一个 Todo，当前仅当 Token 所指的用户拥有该 Todo。

GET /todo

```
{
    "token": string
}
```

返回属于被认证用户的 todo 列表

POST /todo

```
{
    "token": string
    "todo": Todo
}
```

创建一个新的属于被创建用户的 Todo。

PUT /todo/:id

```
{
    "token": string
    "todo": Todo
}
```

更新 ID 为 `id` 的 Todo。

GET /todo/:id

```
{
    "token": string
}
```

返回 ID 为 `id` 的 Todo。

DELETE /todo/:id

```
{
    "token": string
}
```

删除 ID 为 `id` 的 Todo。
