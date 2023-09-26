# Go-Todo

## API

### 登录
endpoint: `/login`

method: `POST`

|   字段名    |  数据类型  | 说明  |
|:--------:|:------:|:---:|
| username | string | 用户名 |
| password | string | 密码  |

### 注册
endpoint: `/register`

method: `POST`

|   字段名    |  数据类型  | 说明  |
|:--------:|:------:|:---:|
| username | string | 用户名 |
| password | string | 密码  |


### 修改密码
endpoint: `/user/:username/reset`

method: `POST`

| 字段名 |  数据类型  | 说明  |
|:---:|:------:|:---:|
| old | string | 旧密码 |
| new | string | 新密码 |

### 登出
endpoint: `/user/:username/logout`

method: `GET`

|   字段名    |  数据类型  | 说明  |
|:--------:|:------:|:---:|

### 添加Todo
endpoint: `/user/:username/add`

method: `POST`

|  字段名   |  数据类型  |       说明       |
|:------:|:------:|:--------------:|
|  name  | string |      事项名       |
| detail | string |      事项描述      |
|  time  | int64  | 一定时间后截至(单位: 秒) |

### 列出Todo
endpoint: `/user/:username/list`

method: `GET`

| 字段名  | 数据类型 |  说明  |
|:----:|:----:|:----:|
| from | int  | 起始位置 |
| num  | int  | 总取数量 |

### 获取单项Todo
endpoint: `/user/:username/get`

method: `GET`

|  字段名   |  数据类型  |       说明       |
|:------:|:------:|:--------------:|
|  name  | string |      事项名       |

### 添加Todo
endpoint: `/user/:username/delete`

method: `POST`

|  字段名   |  数据类型  |       说明       |
|:------:|:------:|:--------------:|
|  name  | string |      事项名       |

### 修改Todo截至时间
endpoint: `/user/:username/update`

method: `POST`

|  字段名   |  数据类型  |         说明          |
|:------:|:------:|:-------------------:|
|  name  | string |         事项名         |
|  time  | int64  | 从当前时间一定时间后截至(单位: 秒) |

### 设置接收邮箱
endpoint: `/user/:username/email`

method: `POST`

|  字段名  |  数据类型  |  说明  |
|:-----:|:------:|:----:|
| email | string | 邮箱地址 |

### 获取全局Todo(需要特殊权限)
endpoint: `/admin/list`

method: `GET`

| 字段名  | 数据类型 |  说明  |
|:----:|:----:|:----:|
| from | int  | 起始位置 |
| num  | int  | 总取数量 |

### 为某位用户添加Todo(需要特殊权限)
endpoint: `/admin/add`

method: `POST`

|   字段名    |  数据类型  |       说明       |
|:--------:|:------:|:--------------:|
|   name   | string |      事项名       |
|  detail  | string |      事项描述      |
|   time   | int64  | 一定时间后截至(单位: 秒) |
| username | string |      用户名       |

### 删除某位用户Todo(需要特殊权限)
endpoint: `/admin/delete`

method: `POST`

|   字段名    |  数据类型  | 说明  |
|:--------:|:------:|:---:|
|   name   | string | 事项名 |
| username | string | 用户名 |

### Response

统一为

|  字段名   |  数据类型  |                   说明                   |
|:------:|:------:|:--------------------------------------:|
| status | string |                  状态描述                  |
|  code  |  uint  |                  状态码                   |
|  time  | string | 时间 (2023-09-26T12:08:08.0012063+08:00) |
|  data  |  map   |                 用户数据相关                 |

`data` 数据主要包括

|   字段名    |  数据类型  |                   说明                   |
|:--------:|:------:|:--------------------------------------:|
| username | string |                  状态描述                  |
|  token   | string |                  状态码                   |
|   time   | string | 时间 (2023-09-26T12:08:08.0012063+08:00) |
|   ...    |  ...   |                  数据相关                  |

当请求错误或未经行身份验证，`data` 字段将为 `null`

## StatusCode

| Code |       Status       |
|:----:|:------------------:|
|  1   |   DatabaseError    |
|  2   |  UserExistsError   |
|  3   | UserNotExistsError |
|  4   |  ChangePWDFailed   |
|  5   |   ListTodoFailed   |
|  6   |    LoginFailed     |
|  7   |  ItemExistsError   |
|  8   | ItemNotExistsError |
|  9   |  ChangeTodoFailed  |
|  10  |  EmailFormatError  |
|  11  |   SetEmailFailed   |
|  0   |  ChangePWDSuccess  |
|  0   |   AddTodoSuccess   |
|  0   |  ListTodoSuccess   |
|  0   |   DelTodoSuccess   |
|  0   |    LoginSuccess    |
|  0   |  RegisterSuccess   |
|  0   | ChangeTodoSuccess  |
|  0   |  SetEmailSuccess   |
|  0   |   GetUserSuccess   |

## Config

default.yml
```yaml
server:
  port: 8000
db: data.db
salt: -S@ltPaThlJbnRz
email:
  sender:
  server:
  port: 25
  secret:
```

|     字段名      |        说明        |
|:------------:|:----------------:|
| server.port  |        端口        |
|      db      |       数据库        |
|     salt     | hash(pwd + salt) |
| email.sender |        账号        |
| email.sender |     SMTP 服务器     |
|  email.port  |     SMTP 端口      |
|    secret    |    账号密码(授权码)     |