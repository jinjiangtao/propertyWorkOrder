# 物业报修管理系统

一个完整的物业报修管理前后端分离系统，使用 Go Gin 框架作为后端，Vue3 + Element Plus 作为前端。

## 项目结构

```
propertyWorkOrder/
├── server/                 # 后端代码
│   ├── main.go           # 主程序入口
│   ├── database.go       # 数据库初始化和模型定义
│   ├── handlers_auth.go # 用户认证接口
│   ├── handlers_repair.go # 报修相关接口
│   ├── handlers_test.go # 后端测试用例
│   ├── go.mod           # Go 模块依赖
│   └── go.sum           # Go 模块锁定
│
├── web/                   # 前端代码
│   ├── src/
│   │   ├── api/         # API 调用封装
│   │   │   └── index.js
│   │   ├── router/      # 路由配置
│   │   │   └── index.js
│   │   ├── views/       # 页面组件
│   │   │   ├── admin/   # 管理员页面
│   │   │   │   ├── Login.vue
│   │   │   │   └── Home.vue
│   │   │   └── user/    # 用户页面
│   │   │       ├── Login.vue
│   │   │       ├── Register.vue
│   │   │       ├── Home.vue
│   │   │       └── Repair.vue
│   │   ├── App.vue      # 根组件
│   │   └── main.js      # 前端入口
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
│
├── prd01.md              # 产品需求文档
└── README.md             # 项目说明文档
```

## 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin Web Framework
- **数据库**: SQLite
- **密码加密**: bcrypt

### 前端
- **框架**: Vue 3
- **UI 库**: Element Plus
- **构建工具**: Vite
- **HTTP 客户端**: Axios

## 功能特性

### 用户功能
- 用户注册和登录
- 提交报修申请（报修类型、描述、图片）
- 查看个人报修记录
- 查看报修状态

### 管理员功能
- 管理员登录（admin/123456）
- 查看所有用户的报修记录
- 更新报修状态（未处理、处理中、已完成）

## 快速开始

### 环境要求
- Go 1.21 或更高版本
- Node.js 16 或更高版本
- npm 或 yarn

### 1. 启动后端服务

```bash
cd server
go mod tidy
go run main.go
```

后端服务将在 http://localhost:8080 启动

### 2. 构建前端

```bash
cd web
npm install
npm run build
```

前端静态文件将生成在 `web/dist` 目录

### 3. 访问应用

启动后端后，可以通过以下地址访问：

- **用户页面**: http://localhost:8080/user
- **管理员页面**: http://localhost:8080/admin

### 4. 开发模式（前端）

如果需要单独开发前端，可以运行：

```bash
cd web
npm install
npm run dev
```

前端开发服务器将在 启动

## API 接口

### 用户认证
- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录

### 报修管理
- `POST /api/repair/create` - 创建报修
- `GET /api/repair/user` - 获取用户报修列表
- `GET /api/repair/all` - 获取所有报修（管理员）
- `PUT /api/repair/status` - 更新报修状态（管理员）

## 默认账户

### 管理员
- 用户名: `admin`
- 密码: `123456`

## 运行测试

```bash
cd server
go test -v ./...
```

## 数据库

SQLite 数据库文件 `property_work_order.db` 将在后端首次启动时自动创建。

数据库包含以下表：
- `users` - 用户表
- `repair_requests` - 报修申请表

## 注意事项

1. 首次运行后端会自动创建管理员账户
2. 前端构建后，静态文件由后端服务托管
3. 所有 API 请求支持 CORS
4. 密码使用 bcrypt 加密存储

第2轮：
交互过程请用中文展示
刚才的产物中为看到admin 的页面，请输出admin的产物，admin页面不用注册，默认用户名密码 admin 123456
用户页面用用户友好的样式展示。h5 + vue ，在手机中展示