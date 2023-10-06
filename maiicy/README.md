# Todo List API

这个API允许你管理你的待办事项列表

## API接口

| Endpoint             | Method | Description                |
| -------------------- | ------ | -------------------------- |
| `/api/user/register` | POST   | 注册用户                   |
| `/api/user/login`    | POST   | 登陆用户获取jwt token      |
| `/api/user/logout`   | POST   | 登出用户                   |
| `/api/todo/delete`   | POST   | 删除任务                   |
| `/api/todo/add`      | POST   | 添加任务                   |
| `/api/todo/update`   | POST   | 更新任务                   |
| `/api/todo/:id"`     | GET    | 通过任务ID查看任务         |
| `/api/todo/:date`    | GET    | 通过日期查看日期之前的任务 |
| `/api/todo/all`      | GET    | 查看所有任务               |

## 登陆认证

要验证并获得JWT令牌，请使用有效凭据向' /api/auth/login '发出POST请求。然后，您可以将JWT令牌包含在其他请求的“Authorization”头中。

## 请求参数

### 注册用户 (POST `/api/user/register`)

- **username** (string): 用户名
- **password** (string): 密码

### 登陆用户 (POST `/api/user/login`)

- **username** (string): 用户名
- **password** (string): 密码

### 登出用户 (POST `/api/user/logout`)

> 需要Authorization

### 添加任务 (POST `/api/todo/add`)

> 需要Authorization

- **title** (string): 任务名称
- **date** (string): 任务日期（`ISO 8601`）

### 删除任务 (POST `/api/todo/delete`)

> 需要Authorization

- **todo_id** (string): 任务id

### 更新任务 (POST `/api/todo/update`)

> 需要Authorization

- **todo_id** (string): 任务id

- **title** (string): 任务名称
- **date** (string): 任务日期（`ISO 8601`）
- **completed** (bool): 任务是否完成

### 错误回复体

- Status Code: 404 Not Found

- Response Body:

  ```json
  {
    "error": "Todo not found"
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
    }
  ],
  "message": "获取成功"
}

```

