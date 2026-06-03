# 个人碳足迹追踪平台

基于 Go (Gin) + Vue 3 + MySQL 开发的个人碳足迹追踪系统。

## 功能特性

- **碳排记录**: 记录每日出行方式、用电量和饮食习惯
- **智能计算**: 根据碳排因子库自动计算当日总排放量
- **趋势分析**: 展示周/月碳排放趋势图表
- **区域对比**: 与所在区域平均值进行对比
- **减排目标**: 支持设定减排目标，达成时推送鼓励语
- **数据持久化**: 所有数据保存在 MySQL 数据库中

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.21 + Gin + GORM |
| 前端 | Vue 3 + Element Plus + ECharts |
| 数据库 | MySQL 8.0 |
| 认证 | JWT |

## 项目结构

```
├── server/                 # 后端代码
│   ├── main.go            # 入口文件
│   ├── config/            # 配置管理
│   ├── middleware/        # 中间件(JWT认证)
│   ├── models/            # 数据模型
│   ├── routes/            # 路由处理器
│   └── services/          # 业务逻辑
├── frontend/              # 前端代码
│   ├── src/
│   │   ├── views/         # 页面组件
│   │   ├── api/           # API 封装
│   │   ├── router/        # 路由配置
│   │   └── App.vue        # 根组件
│   └── package.json
└── database/
    └── init.sql           # 数据库初始化脚本
```

## 快速开始

### 1. 初始化数据库

```bash
mysql -u root -p < database/init.sql
```

### 2. 启动后端

```bash
cd server

# 安装依赖
go mod tidy

# 配置环境变量(可选，默认值已预设)
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=123456
export DB_NAME=carbon_tracker
export SERVER_PORT=8080

# 运行
go run main.go
```

### 3. 启动前端

```bash
cd frontend

# 安装依赖
npm install

# 开发模式运行
npm run dev
```

访问 http://localhost:5173 即可使用。

## API 接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/register | 用户注册 | 否 |
| POST | /api/login | 用户登录 | 否 |
| GET | /api/factors | 获取碳排因子 | 否 |
| POST | /api/records | 创建碳排记录 | 是 |
| GET | /api/records | 获取当日记录 | 是 |
| DELETE | /api/records/:id | 删除记录 | 是 |
| GET | /api/trend/weekly | 获取周趋势 | 是 |
| GET | /api/trend/monthly | 获取月趋势 | 是 |
| GET | /api/summary/category | 分类统计 | 是 |
| POST | /api/goals | 设定减排目标 | 是 |
| GET | /api/goals/active | 获取当前目标 | 是 |
| GET | /api/goals/history | 目标历史 | 是 |

## 碳排因子说明

| 分类 | 项目 | 因子 | 单位 |
|------|------|------|------|
| 出行 | 私家车 | 0.21 | kgCO₂/km |
| 出行 | 公交车 | 0.08 | kgCO₂/km |
| 出行 | 地铁 | 0.05 | kgCO₂/km |
| 出行 | 骑行 | 0.00 | kgCO₂/km |
| 用电 | 家庭用电 | 0.785 | kgCO₂/kWh |
| 饮食 | 牛肉 | 27.0 | kgCO₂/kg |
| 饮食 | 猪肉 | 12.1 | kgCO₂/kg |
| 饮食 | 蔬菜 | 2.0 | kgCO₂/kg |
