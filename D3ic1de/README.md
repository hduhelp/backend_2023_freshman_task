# Go-TodoList

------

**实现的功能**

- 基本Todo的增删改查
- 用户鉴权

## API

------



### 登录

**POST**: `/login`

| 参数名   | 参数类型 | 备注   |
| -------- | -------- | ------ |
| username | string   | 用户名 |
| password | string   | 密码   |

### 注册

**POST**: `/register`

| 参数名   | 参数类型 | 备注   |
| -------- | -------- | ------ |
| username | string   | 用户名 |
| password | string   | 密码   |

------

### 添加Todo

**POST**: `/todo`

```json
{
    "content": string,
    "done": bool
}
```

- `content`: 内容
- `done`: 是否完成

### 获取Todo

```
GET /todo

Cookie: token=登录后携带的cookie
```

返回TodoList

### 查询Todo

```
GET /todo/:index

Cookie: token=登录后携带的cookie
```

返回查询第index个的todo

### 修改Todo

```
PUT /todo/:index

Cookie: token=登录后携带的cookie

{
    "content": string,
    "done": bool
}
```

**响应**

```json
{
	"code": int,
    "message": string
}
```

- `code`: 响应状态码
- `message`: 响应成功/错误

### 删除Todo

```
DELETE /todo/:index

Cookie: token=登录后携带的cookie

{
    "content": string,
    "done": bool
}
```

**响应**同上

