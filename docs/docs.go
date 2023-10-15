// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate_swagger = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/handle_op": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "handle op operation",
                "tags": [
                    "swapbotAA业务接口"
                ],
                "parameters": [
                    {
                        "description": "快速交易参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/vo.HandleOpVO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/limit_order": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "limit order operation",
                "tags": [
                    "swapbotAA业务接口"
                ],
                "parameters": [
                    {
                        "description": "限价单参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/vo.LimitOrderVO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "用户业务接口"
                ],
                "summary": "登录业务",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "description": "登录请求参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/refresh_token": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "刷新accessToken",
                "tags": [
                    "用户业务接口"
                ],
                "summary": "刷新accessToken",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "可以为空",
                        "name": "community_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "score",
                        "description": "排序依据",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "关键字搜索",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "每页数量",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/signup": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "注册业务",
                "tags": [
                    "用户业务接口"
                ],
                "summary": "注册业务",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "description": "注册请求参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RegisterForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.ResponseData": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "description": "omitempty当data为空时,不展示这个字段"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.LoginForm": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.RegisterForm": {
            "type": "object",
            "required": [
                "confirm_password",
                "email",
                "password",
                "username"
            ],
            "properties": {
                "confirm_password": {
                    "type": "string"
                },
                "email": {
                    "description": "邮箱",
                    "type": "string"
                },
                "gender": {
                    "description": "性别 0:未知 1:男 2:女",
                    "type": "integer",
                    "enum": [
                        0,
                        1,
                        2
                    ]
                },
                "password": {
                    "description": "密码",
                    "type": "string"
                },
                "username": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        },
        "vo.HandleOpVO": {
            "type": "object",
            "properties": {
                "beneficiary_addr": {
                    "type": "string"
                },
                "user_op": {
                    "$ref": "#/definitions/vo.UserOperationVO"
                }
            }
        },
        "vo.LimitOrderVO": {
            "type": "object",
            "required": [
                "amount_in",
                "amount_out",
                "amount_out_minimum",
                "fee",
                "token_in",
                "token_out",
                "user"
            ],
            "properties": {
                "amount_in": {
                    "type": "string"
                },
                "amount_out": {
                    "type": "string"
                },
                "amount_out_minimum": {
                    "type": "string"
                },
                "fee": {
                    "type": "integer"
                },
                "token_in": {
                    "type": "string"
                },
                "token_out": {
                    "type": "string"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "vo.UserOperationVO": {
            "type": "object",
            "properties": {
                "call_data": {
                    "type": "string"
                },
                "call_gas_limit": {
                    "type": "string"
                },
                "init_code": {
                    "type": "string"
                },
                "max_fee_per_gas": {
                    "type": "string"
                },
                "max_priority_fee_per_gas": {
                    "type": "string"
                },
                "nonce": {
                    "type": "string"
                },
                "paymaster_and_data": {
                    "type": "string"
                },
                "pre_verification_gas": {
                    "type": "string"
                },
                "sender": {
                    "type": "string"
                },
                "signature": {
                    "type": "string"
                },
                "verification_gas_limit": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "127.0.0.1:8081",
	BasePath:         "/api/v1/",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate_swagger,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
