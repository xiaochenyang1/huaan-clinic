// Package docs 华安医疗预约系统 API
//
// 这是Swagger文档占位文件，实际文档由swag命令自动生成
//
// 生成命令：swag init -g cmd/main.go -o docs
//
// Schemes: http, https
// Host: localhost:8080
// BasePath: /api
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package docs

// 注意：此文件为占位文件
// 请运行以下命令生成实际的Swagger文档：
// 1. 安装swag: go install github.com/swaggo/swag/cmd/swag@latest
// 2. 生成文档: swag init -g cmd/main.go -o docs

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "title": "华安医疗预约系统 API",
        "description": "华安医疗门诊预约系统后端API接口文档",
        "version": "1.0.0",
        "contact": {
            "name": "API Support",
            "email": "support@huaan-medical.com"
        }
    },
    "host": "localhost:8080",
    "basePath": "/api",
    "schemes": ["http", "https"],
    "paths": {},
    "definitions": {}
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}{
	Version:     "1.0.0",
	Host:        "localhost:8080",
	BasePath:    "/api",
	Schemes:     []string{"http", "https"},
	Title:       "华安医疗预约系统 API",
	Description: "华安医疗门诊预约系统后端API接口文档",
}

func init() {
	// 初始化Swagger信息
}
