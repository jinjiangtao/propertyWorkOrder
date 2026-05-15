# 物业报修管理系统

## 项目简介

物业报修管理系统是一个基于 Go Gin + Vue3 + Element Plus 的全栈应用，用于管理用户报修请求和管理员处理流程。

## 技术栈

### 后端
- Go 1.21+
- Gin 框架
- SQLite 数据库

### 前端
- Vue 3
- Vue Router
- Element Plus
- Vite

## 功能特性

### 用户功能
- 用户注册/登录
- 提交报修申请（报修类型、描述、图片）
- 查看个人报修记录

### 管理员功能
- 管理员登录
- 查看所有报修记录
- 更新报修状态（未处理 -> 处理中 -> 已完成）

## 项目结构

```
propertyWorkOrder/
├── backend/                 # 后端代码
│   ├── controller/          # 控制器层
│   ├── entity/              # 实体模型
│   ├── middleware/          # 中间件
│   ├── model/               # 数据模型层
│   ├── router/              # 路由配置
│   ├── service/             # 业务逻辑层
│   ├── go.mod               # Go 依赖
│   └── main.go              # 入口文件
├── frontend/                # 前端代码
│   ├── src/
│   │   ├── api/             # API 请求
│   │   ├── router/          # 路由配置
│   │   ├── utils/           # 工具函数
│   │   ├── views/           # 页面组件
│   │   ├── App.vue          # 根组件
│   │   └── main.js          # 入口文件
│   ├── index.html           # HTML 模板
│   ├── package.json         # 前端依赖
│   └── vite.config.js       # Vite 配置
├── README.md                # 项目说明
└── prd01.md                 # 需求文档
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- npm 或 yarn

### 启动后端服务

```bash
cd backend

go mod download

go run main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 启动前端服务

```bash
cd frontend

npm install

npm run dev
```

前端服务将在 `http://localhost:5173` 启动。

## API 接口

### 用户接口
- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录
- `POST /api/admin/login` - 管理员登录

### 报修接口
- `POST /api/workorder` - 创建报修
- `GET /api/workorders?user_id=xxx` - 获取报修列表（带user_id参数获取个人记录，不带获取全部）
- `PUT /api/workorder/:id` - 更新报修状态

## 数据库

系统使用 SQLite 数据库，数据库文件为 `backend/example_db.sqlite`。

### 默认管理员账号
- 用户名: `admin`
- 密码: `admin123`

## 测试

### API 测试

可以使用 curl 或 Postman 测试 API：

```bash
# 用户注册
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'

# 用户登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'

# 管理员登录
curl -X POST http://localhost:8080/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 创建报修
curl -X POST http://localhost:8080/api/workorder \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"type":"水电维修","description":"水管漏水","images":""}'

# 获取报修列表
curl http://localhost:8080/api/workorders

# 更新报修状态
curl -X PUT http://localhost:8080/api/workorder/1 \
  -H "Content-Type: application/json" \
  -d '{"status":2}'
```

### 前端页面测试

启动服务后访问 `http://localhost:5173`，测试以下功能：
1. 用户注册和登录
2. 提交报修申请
3. 查看报修记录
4. 管理员登录
5. 处理报修请求

## 开发说明

### 数据库表结构

**users 表**
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自增 |
| username | TEXT | 用户名，唯一 |
| password | TEXT | 密码 |
| role | INTEGER | 角色（1-用户，2-管理员） |

**work_orders 表**
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自增 |
| user_id | INTEGER | 用户ID |
| type | TEXT | 报修类型 |
| description | TEXT | 报修描述 |
| images | TEXT | 图片URL |
| status | INTEGER | 状态（1-未处理，2-处理中，3-已完成） |
| created_at | TEXT | 创建时间 |
| updated_at | TEXT | 更新时间 |