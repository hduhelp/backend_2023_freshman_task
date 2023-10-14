# Todo List API

这个API允许你管理你的待办事项列表

## API接口

| Endpoint                  | Method | Description                |
| ------------------------- | ------ | -------------------------- |
| `/api/users/register`     | POST   | 注册用户                   |
| `/api/users/login`        | POST   | 登陆用户获取jwt token      |
| `/api/users/logout`       | DELETE | 登出用户                   |
| `/api/todos/:id`          | DELETE | 删除任务                   |
| `/api/todos/add`          | POST   | 添加任务                   |
| `/api/todos/:id`          | PUT    | 更新任务                   |
| `/api/todos/:id`          | GET    | 通过任务ID查看任务         |
| `/api/todos/before/:date` | GET    | 通过日期查看日期之前的任务 |
| `/api/todos/all`          | GET    | 查看所有任务               |

## 配置文件

配置文件是一个 TOML 格式文件，通常命名为 `config.toml`。

以下是一个示例配置文件：

```
# 服务器配置
ServerIP = "localhost"
ServerPort = 8080

# 数据库配置
DBPath = "database.db"
```

你可以根据你的项目需求自定义配置文件，确保包括必要的参数和值。

> jwt的秘钥`SECRET_KEY`需要存储在环境变量 (* 必要)

## CLI 使用

程序还支持命令行界面（CLI）以下是如何使用 CLI 参数：

```
shellCopy code
$ ./your-app -dbPath /path/to/database.db -host localhost -port 8081
```

### 可用 CLI 参数

- `-dbPath`: 指定数据库文件的路径。
- `-host`: 指定服务器主机名。
- `-port`: 指定服务器端口号。

如果你提供了这些 CLI 参数，它们将覆盖配置文件中的对应值。



## 登陆认证

要验证并获得JWT令牌，请使用有效凭据向' /api/auth/login '发出POST请求。然后，您可以将JWT令牌包含在其他请求的“Authorization”头中。

## 请求参数

### 注册用户 (POST `/api/users/register`)

- **username** (string): 用户名
- **password** (string): 密码
- **info**(string): 简介

### 登陆用户 (POST `/api/users/login`)

- **username** (string): 用户名
- **password** (string): 密码

### 登出用户 (DELETE `/api/users/logout`)

> 需要Authorization

### 删除任务 (DELETE `/api/todos/:id`)

> 需要Authorization

### 添加任务 (POST `/api/todos/add`)

> 需要Authorization

- **title** (string): 任务名称
- **date** (string): 任务日期（`ISO 8601`）

### 更新任务 (PUT `/api/todos/:id`)

> 需要Authorization

- **todo_id** (string): 任务id

- **title** (string): 任务名称
- **date** (string): 任务日期（`ISO 8601`）
- **completed** (bool): 任务是否完成

### 查看之前任务 (GET `/api/todos/before/:date`)

> 需要Authorization

### 错误回复体

- Status Code: 404

- Response Body:

  ```json
  {
    "error": "错误内容"
  }
  ```

## 回复体示例

`/api/todo/all`获取所有任务

```json
{
  "data": [
    {
      "ID": 1,
      "UserID": 2,
      "Title": "Test yeah",
      "Completed": false,
      "DueDate": "2023-10-13T00:00:00Z",
      "CreatedAt": "2023-10-06T16:54:51Z"
    },
    {
      "ID": 2,
      "UserID": 2,
      "Title": "Test yeah2",
      "Completed": false,
      "DueDate": "2023-10-13T00:00:00Z",
      "CreatedAt": "2023-10-06T16:54:59Z"
    }
  ],
  "message": "获取成功"
}

```

`/api/todo/:id`获取单个任务


```json
{
    "data": {
        "ID": 5,
        "UserID": 3,
        "Title": "TT day",
        "Completed": false,
        "DueDate": "2023-10-13T00:00:00Z",
        "CreatedAt": "2023-10-12T16:09:42Z"
    },
    "message": "获取成功"
}
```

