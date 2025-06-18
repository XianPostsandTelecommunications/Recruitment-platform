# API 接口文档

## 基础信息

- **基础URL**: `http://localhost:8080/api`
- **认证方式**: JWT Bearer Token
- **数据格式**: JSON
- **字符编码**: UTF-8

## 通用响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {}
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "错误信息",
  "errors": [
    {
      "field": "email",
      "message": "邮箱格式不正确"
    }
  ]
}
```

## 认证相关接口

### 用户注册

**接口地址**: `POST /auth/register`

**请求参数**:
```json
{
  "username": "string",     // 用户名，必填，长度3-20
  "email": "string",        // 邮箱，必填，格式正确
  "password": "string",     // 密码，必填，长度6-20
  "role": "student"         // 角色，可选，默认student
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "id": 1,
    "username": "张三",
    "email": "zhangsan@example.com",
    "role": "student",
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

### 用户登录

**接口地址**: `POST /auth/login`

**请求参数**:
```json
{
  "email": "string",        // 邮箱，必填
  "password": "string"      // 密码，必填
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "张三",
      "email": "zhangsan@example.com",
      "role": "student",
      "avatar": "https://example.com/avatar.jpg"
    }
  }
}
```

### 获取用户信息

**接口地址**: `GET /auth/profile`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "username": "张三",
    "email": "zhangsan@example.com",
    "role": "student",
    "avatar": "https://example.com/avatar.jpg",
    "phone": "13800138000",
    "studentId": "2021001",
    "major": "计算机科学与技术",
    "grade": "2021级",
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

## 实验室管理接口

### 获取实验室列表

**接口地址**: `GET /labs`

**查询参数**:
- `page`: 页码，默认1
- `size`: 每页数量，默认10
- `search`: 搜索关键词
- `tags`: 标签筛选，多个用逗号分隔
- `status`: 状态筛选，active/inactive

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "人工智能实验室",
        "description": "专注于AI技术研究和应用开发",
        "requirements": "熟悉Python，有机器学习基础，对AI技术有浓厚兴趣",
        "maxMembers": 15,
        "currentMembers": 8,
        "contactEmail": "ai@university.edu",
        "contactPhone": "010-12345678",
        "location": "计算机学院A楼301",
        "tags": ["AI", "机器学习", "深度学习"],
        "coverImage": "https://example.com/lab1.jpg",
        "status": "active",
        "createdBy": 1,
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 10,
    "totalPages": 10
  }
}
```

### 获取实验室详情

**接口地址**: `GET /labs/:id`

**路径参数**:
- `id`: 实验室ID

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "name": "人工智能实验室",
    "description": "专注于AI技术研究和应用开发，致力于推动人工智能技术的发展和应用。",
    "requirements": "1. 熟悉Python编程语言\n2. 有机器学习基础\n3. 对AI技术有浓厚兴趣\n4. 有良好的团队协作能力",
    "maxMembers": 15,
    "currentMembers": 8,
    "contactEmail": "ai@university.edu",
    "contactPhone": "010-12345678",
    "location": "计算机学院A楼301",
    "tags": ["AI", "机器学习", "深度学习"],
    "coverImage": "https://example.com/lab1.jpg",
    "status": "active",
    "createdBy": 1,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

### 创建实验室

**接口地址**: `POST /labs`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:
```json
{
  "name": "string",           // 实验室名称，必填
  "description": "string",    // 实验室描述，必填
  "requirements": "string",   // 招新要求，必填
  "maxMembers": 10,          // 最大成员数，必填
  "contactEmail": "string",   // 联系邮箱，必填
  "contactPhone": "string",   // 联系电话，可选
  "location": "string",       // 实验室位置，可选
  "tags": ["tag1", "tag2"]    // 标签数组，可选
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "实验室创建成功",
  "data": {
    "id": 2,
    "name": "区块链实验室",
    "description": "研究区块链技术和应用",
    "requirements": "熟悉Go语言，了解密码学基础",
    "maxMembers": 12,
    "currentMembers": 0,
    "contactEmail": "blockchain@university.edu",
    "contactPhone": "010-87654321",
    "location": "计算机学院B楼205",
    "tags": ["区块链", "密码学", "分布式系统"],
    "status": "active",
    "createdBy": 1,
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

### 更新实验室

**接口地址**: `PUT /labs/:id`

**请求头**: `Authorization: Bearer <token>`

**路径参数**:
- `id`: 实验室ID

**请求参数**: 同创建实验室，所有字段可选

**响应示例**:
```json
{
  "code": 200,
  "message": "实验室更新成功",
  "data": {
    "id": 1,
    "name": "人工智能实验室（更新）",
    "description": "更新后的描述",
    "maxMembers": 20,
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

### 删除实验室

**接口地址**: `DELETE /labs/:id`

**请求头**: `Authorization: Bearer <token>`

**路径参数**:
- `id`: 实验室ID

**响应示例**:
```json
{
  "code": 200,
  "message": "实验室删除成功"
}
```

## 申请管理接口

### 提交申请

**接口地址**: `POST /applications`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:
```json
{
  "labId": 1,                    // 实验室ID，必填
  "motivation": "string",        // 申请动机，必填
  "skills": ["skill1", "skill2"], // 技能列表，必填
  "experience": "string",        // 相关经验，可选
  "availableTime": "string"      // 可用时间，可选
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "申请提交成功",
  "data": {
    "id": 1,
    "userId": 2,
    "labId": 1,
    "motivation": "对AI技术有浓厚兴趣，希望能在实验室中深入学习",
    "skills": ["Python", "机器学习", "深度学习"],
    "experience": "参加过机器学习课程，完成过图像分类项目",
    "availableTime": "每周20小时",
    "status": "pending",
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

### 获取我的申请

**接口地址**: `GET /applications/my`

**请求头**: `Authorization: Bearer <token>`

**查询参数**:
- `page`: 页码，默认1
- `size`: 每页数量，默认10
- `status`: 状态筛选，pending/accepted/rejected

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 1,
        "labId": 1,
        "labName": "人工智能实验室",
        "motivation": "对AI技术有浓厚兴趣",
        "skills": ["Python", "机器学习"],
        "experience": "参加过相关项目",
        "availableTime": "每周20小时",
        "status": "pending",
        "feedback": "",
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "size": 10
  }
}
```

### 获取实验室申请列表

**接口地址**: `GET /applications/lab/:labId`

**请求头**: `Authorization: Bearer <token>`

**路径参数**:
- `labId`: 实验室ID

**查询参数**:
- `page`: 页码，默认1
- `size`: 每页数量，默认10
- `status`: 状态筛选

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 1,
        "userId": 2,
        "userName": "李四",
        "userEmail": "lisi@example.com",
        "motivation": "对AI技术有浓厚兴趣",
        "skills": ["Python", "机器学习"],
        "experience": "参加过相关项目",
        "availableTime": "每周20小时",
        "status": "pending",
        "feedback": "",
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "size": 10
  }
}
```

### 审核申请

**接口地址**: `PUT /applications/:id/status`

**请求头**: `Authorization: Bearer <token>`

**路径参数**:
- `id`: 申请ID

**请求参数**:
```json
{
  "status": "accepted",      // 状态，必填：pending/accepted/rejected
  "feedback": "string"       // 审核反馈，可选
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "申请审核成功",
  "data": {
    "id": 1,
    "status": "accepted",
    "feedback": "欢迎加入我们的团队！",
    "reviewedBy": 1,
    "reviewedAt": "2024-01-01T00:00:00Z"
  }
}
```

## 通知管理接口

### 获取通知列表

**接口地址**: `GET /notifications`

**请求头**: `Authorization: Bearer <token>`

**查询参数**:
- `page`: 页码，默认1
- `size`: 每页数量，默认10
- `isRead`: 是否已读，true/false

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 1,
        "title": "申请审核结果",
        "content": "您的实验室申请已通过审核，欢迎加入！",
        "type": "application",
        "isRead": false,
        "relatedId": 1,
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 20,
    "page": 1,
    "size": 10
  }
}
```

### 标记通知为已读

**接口地址**: `PUT /notifications/:id/read`

**请求头**: `Authorization: Bearer <token>`

**路径参数**:
- `id`: 通知ID

**响应示例**:
```json
{
  "code": 200,
  "message": "标记成功",
  "data": {
    "id": 1,
    "isRead": true
  }
}
```

## 统计接口

### 获取仪表板统计数据

**接口地址**: `GET /stats/dashboard`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "totalUsers": 150,
    "totalLabs": 25,
    "totalApplications": 300,
    "pendingApplications": 50,
    "acceptedApplications": 200,
    "rejectedApplications": 50,
    "recentApplications": [
      {
        "id": 1,
        "userName": "张三",
        "labName": "人工智能实验室",
        "status": "pending",
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "labStats": [
      {
        "labName": "人工智能实验室",
        "applicationCount": 45,
        "acceptanceRate": 0.8
      }
    ]
  }
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 422 | 数据验证失败 |
| 500 | 服务器内部错误 |

## 权限说明

### 学生权限
- 查看实验室列表和详情
- 提交申请
- 查看自己的申请
- 查看通知
- 管理个人资料

### 管理员权限
- 所有学生权限
- 创建、编辑、删除实验室
- 审核申请
- 查看统计数据
- 管理用户
- 发送通知 