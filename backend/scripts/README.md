# 数据库初始化指南

本文档说明如何初始化华安医疗预约系统的数据库。

## 前置条件

1. ✅ 已安装 MySQL 8.0+
2. ✅ MySQL 服务正在运行
3. ✅ 已配置 `config.yaml` 中的数据库连接信息
4. ✅ 已安装 Go 1.21+

## 初始化步骤

### 方法一：一键初始化（推荐）

在 `backend` 目录下执行以下命令：

```bash
# 1. 创建数据库（使用 MySQL 客户端）
mysql -u root -p < scripts/init_db.sql

# 2. 运行程序（自动创建表结构）
go run cmd/main.go

# 程序启动后会自动执行：
# - 连接数据库
# - AutoMigrate 创建所有表
# - 启动 HTTP 服务

# 按 Ctrl+C 停止程序

# 3. 初始化基础数据
go run scripts/seed_data.go
```

### 方法二：分步执行

#### 步骤1：创建数据库

```bash
# 方式A：使用脚本
mysql -u root -p < scripts/init_db.sql

# 方式B：手动创建
mysql -u root -p
```

```sql
DROP DATABASE IF EXISTS huaan_medical;
CREATE DATABASE huaan_medical CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
exit;
```

#### 步骤2：创建表结构

项目使用 GORM AutoMigrate 自动创建表，无需手动建表。

**启动程序即可自动创建：**

```bash
cd backend
go run cmd/main.go
```

程序启动日志中会显示：

```
✅ 数据库连接成功
✅ 数据库迁移成功
```

#### 步骤3：初始化基础数据

```bash
# 在 backend 目录下执行
go run scripts/seed_data.go
```

输出示例：

```
开始初始化数据...
→ 初始化管理员...
  ✅ 成功创建 1 个管理员
→ 初始化科室...
  ✅ 成功创建 10 个科室
→ 初始化角色...
  ✅ 成功创建 3 个角色
→ 初始化权限...
  ✅ 成功创建 24 个权限
✅ 数据初始化完成!

默认管理员账号:
  用户名: admin
  密码: admin123
```

## 初始化的数据

### 1. 管理员账号

| 字段 | 值 |
|-----|-----|
| 用户名 | admin |
| 密码 | admin123 |
| 姓名 | 系统管理员 |
| 手机 | 13800138000 |
| 邮箱 | admin@huaan-medical.com |

### 2. 科室（10个）

- 内科
- 外科
- 儿科
- 妇产科
- 骨科
- 皮肤科
- 眼科
- 耳鼻喉科
- 口腔科
- 中医科

### 3. 角色（3个）

- 超级管理员
- 医院管理员
- 科室管理员

### 4. 权限（24个）

- 科室管理：查看、创建、编辑、删除
- 医生管理：查看、创建、编辑、删除
- 排班管理：查看、创建、编辑、删除
- 预约管理：查看、处理、取消、导出
- 患者管理：查看、编辑
- 系统管理：管理员增删改查
- 数据统计：查看统计、导出数据

## 数据库表结构

GORM AutoMigrate 会自动创建以下表：

| 表名 | 说明 |
|-----|------|
| users | 用户表 |
| patients | 就诊人表 |
| departments | 科室表 |
| doctors | 医生表 |
| schedules | 排班表 |
| appointments | 预约表 |
| medical_records | 就诊记录表 |
| admins | 管理员表 |
| roles | 角色表 |
| permissions | 权限表 |
| operation_logs | 操作日志表 |
| login_logs | 登录日志表 |

## 验证数据库

```bash
# 连接数据库
mysql -u root -p huaan_medical

# 查看所有表
SHOW TABLES;

# 查看管理员
SELECT * FROM admins;

# 查看科室
SELECT * FROM departments;

# 退出
exit;
```

## 常见问题

### 1. 数据库连接失败

检查 `config.yaml` 中的数据库配置：

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: "你的密码"  # 注意加引号
  dbname: huaan_medical
```

### 2. 表已存在错误

如果遇到表已存在的错误，可以：

```sql
-- 删除数据库重新初始化
DROP DATABASE huaan_medical;
CREATE DATABASE huaan_medical CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 权限不足

确保 MySQL 用户有足够权限：

```sql
GRANT ALL PRIVILEGES ON huaan_medical.* TO 'root'@'localhost';
FLUSH PRIVILEGES;
```

### 4. 数据已存在

`seed_data.go` 脚本会自动检查数据是否存在，重复执行不会报错。

## 下一步

数据库初始化完成后，可以：

1. 启动后端服务：`go run cmd/main.go`
2. 访问 API 文档：http://localhost:8080/swagger/index.html
3. 使用管理员账号登录：`admin / admin123`
4. 开始测试 API 接口

## 重置数据库

如果需要完全重置数据库：

```bash
# 1. 删除数据库
mysql -u root -p -e "DROP DATABASE IF EXISTS huaan_medical;"

# 2. 重新执行初始化
mysql -u root -p < scripts/init_db.sql
go run cmd/main.go  # 启动后按 Ctrl+C 停止
go run scripts/seed_data.go
```
