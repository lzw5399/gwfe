// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/process-definitions": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "process-definitions"
                ],
                "summary": "创建流程模板",
                "parameters": [
                    {
                        "description": "request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ProcessDefinitionRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "current-user",
                        "name": "current-user",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.HttpResponse"
                        }
                    }
                }
            }
        },
        "/api/process-instances": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "process-instances"
                ],
                "summary": "获取流程实例列表",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "取的条数",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "跳过的条数",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "asc或者是desc",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "排序键的名字，在各查询实现中默认值与可用值都不同",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "类别",
                        "name": "type",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "current-user",
                        "name": "current-user",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.HttpResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "process-instances"
                ],
                "summary": "创建新的流程实例",
                "parameters": [
                    {
                        "description": "request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ProcessInstanceRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "current-user",
                        "name": "current-user",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.HttpResponse"
                        }
                    }
                }
            }
        },
        "/api/process-instances/_handle": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "process-instances"
                ],
                "summary": "处理/审批一个流程",
                "parameters": [
                    {
                        "description": "request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.HandleInstancesRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "current-user",
                        "name": "current-user",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.HttpResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "request.HandleInstancesRequest": {
            "type": "object",
            "properties": {
                "edgeId": {
                    "description": "走的流程的id",
                    "type": "string"
                },
                "processInstanceId": {
                    "description": "流程实例的id",
                    "type": "integer"
                },
                "remarks": {
                    "description": "备注",
                    "type": "string"
                }
            }
        },
        "request.ProcessDefinitionRequest": {
            "type": "object",
            "properties": {
                "classifyId": {
                    "description": "分类ID",
                    "type": "integer"
                },
                "formId": {
                    "description": "对应的表单的id(仅对外部系统做一个标记)",
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "description": "流程名称",
                    "type": "string"
                },
                "notice": {
                    "description": "绑定通知",
                    "type": "string"
                },
                "remarks": {
                    "description": "流程备注",
                    "type": "string"
                },
                "structure": {
                    "description": "流程结构",
                    "type": "string"
                },
                "task": {
                    "description": "任务ID, array, 可执行多个任务，可以当成通知任务，每个节点都会去执行",
                    "type": "string"
                }
            }
        },
        "request.ProcessInstanceRequest": {
            "type": "object",
            "properties": {
                "processDefinitionId": {
                    "description": "流程ID",
                    "type": "integer"
                },
                "title": {
                    "description": "流程实例标题",
                    "type": "string"
                }
            }
        },
        "response.HttpResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "object"
                },
                "success": {
                    "type": "boolean"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
