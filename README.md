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
│   ├── handlers_worker.go # 工人相关接口
│   ├── handlers_test.go # 后端测试用例
│   ├── go.mod           # Go 模块依赖
│   ├── go.sum           # Go 模块锁定
│   └── property_work_order.db  # SQLite 数据库文件
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
│   │   │   │   ├── Home.vue
│   │   │   │   ├── Workers.vue
│   │   │   │   └── Stats.vue
│   │   │   ├── user/    # 用户页面（H5移动端）
│   │   │   │   ├── Login.vue
│   │   │   │   ├── Register.vue
│   │   │   │   ├── Home.vue
│   │   │   │   └── Repair.vue
│   │   │   └── worker/  # 工人页面
│   │   │       ├── Login.vue
│   │   │       ├── Home.vue
│   │   │       ├── OrderDetail.vue
│   │   │       └── Profile.vue
│   │   ├── App.vue      # 根组件
│   │   └── main.js      # 前端入口
│   ├── dist/            # 构建产物
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
│
├── prd-01.md              # 产品需求文档
├── prd-02.md              # 迭代需求文档
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
- 查看报修状态和维修结果

### 管理员功能
- 管理员登录（admin/123456）
- 查看所有用户的报修记录
- 派单给工人
- 更新报修状态（未处理、已派单、处理中、已完成）
- 工人管理（新增、编辑、禁用/启用）
- 工人工作量统计

### 工人功能
- 工人登录
- 查看分配给自己的工单
- 接单/拒单操作
- 提交维修结果和图片
- 个人中心和统计信息

## 业务流程

1. 用户提交报修 → 状态【未处理】
2. 管理员选择工人派单 → 状态【已派单】
3. 工人接单 → 状态【处理中】；工人可拒单，工单退回未处理
4. 工人维修完成，填写维修描述、上传图片 → 状态【已完成】
5. 用户可查看完整维修记录、维修结果

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
- **工人页面**: http://localhost:8080/worker

### 4. 开发模式（前端）

如果需要单独开发前端，可以运行：

```bash
cd web
npm install
npm run dev
```

前端开发服务器将在 http://localhost:5173 启动

## API 接口

### 用户认证
- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录

### 报修管理
- `POST /api/repair/create` - 创建报修
- `GET /api/repair/user` - 获取用户报修列表
- `GET /api/repair/all` - 获取所有报修（管理员）
- `PUT /api/repair/status` - 更新报修状态
- `POST /api/repair/assign` - 派单给工人（管理员）
- `GET /api/repair/worker` - 获取工人工单列表
- `POST /api/repair/accept` - 工人接单
- `POST /api/repair/reject` - 工人拒单
- `POST /api/repair/result` - 提交维修结果
- `GET /api/repair/stats` - 工人工作量统计

### 工人管理
- `POST /api/worker/login` - 工人登录
- `GET /api/worker/list` - 获取工人列表（管理员）
- `POST /api/worker/create` - 新增工人（管理员）
- `PUT /api/worker/update` - 更新工人信息（管理员）
- `PUT /api/worker/status` - 禁用/启用工人（管理员）

## 默认账户

### 管理员
- 用户名: `admin`
- 密码: `123456`

### 测试工人（需管理员添加）
管理员可以在工人管理页面添加测试工人账户：
- 工号: 唯一标识
- 姓名: 工人姓名
- 手机号: 唯一标识
- 密码: 登录密码
- 工种: 水电/木工/保洁/综合维修

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
- `workers` - 工人表

### 报修状态
- 未处理
- 已派单
- 处理中
- 已完成

### 工种类型
- 水电
- 木工
- 保洁
- 综合维修

## 注意事项

1. 首次运行后端会自动创建管理员账户
2. 前端构建后，静态文件由后端服务托管
3. 所有 API 请求支持 CORS
4. 密码使用 bcrypt 加密存储
5. 用户页面采用 H5 响应式设计，适配移动端
6. 工人端路由受权限控制，需登录后访问