// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Sunbula Hasanova",
            "url": "https://t.me/Neverland_MJ",
            "email": "khasanovasumbula@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/add/card": {
            "post": {
                "description": "user o'zining plastik kartasini kiritishi",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "plastik karta qo'shish",
                "parameters": [
                    {
                        "description": "Card info",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.Card"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    }
                }
            }
        },
        "/add/list/product": {
            "post": {
                "description": "bir nechta produktalarni qo'shish",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "produktalar qo'shish",
                "parameters": [
                    {
                        "description": "Product info",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/product.Product"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    }
                }
            }
        },
        "/add/product": {
            "post": {
                "description": "bitta produkta qo'shish",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "produkta qo'shish",
                "parameters": [
                    {
                        "description": "Product info",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/product.Product"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    }
                }
            }
        },
        "/buy/": {
            "get": {
                "description": "produkta sotib olish",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "buy product",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Buy product",
                        "name": "quantity",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "user login qilishi",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "sign in",
                "parameters": [
                    {
                        "description": "User info",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.PreLoginUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    }
                }
            }
        },
        "/product/list": {
            "get": {
                "description": "barcha produktalarni ko'rsatish",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "hamma produktalarni listi",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/product.Product"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    }
                }
            }
        },
        "/product/{id}": {
            "get": {
                "description": "bitta produkta haqida ma'lumotlarni olish",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "produkta",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/product.Product"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "user registratsiyadan o'tishi",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "sign up",
                "parameters": [
                    {
                        "description": "User info",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.PreSignUpUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "hamma userlar ro'yxatini chiqarish",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "hamma userlar ro'yxati",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/user.UserCard"
                            }
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.message"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.message": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "product.Product": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "image_link": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "original_price": {
                    "type": "integer"
                },
                "price": {
                    "type": "integer"
                },
                "quantity": {
                    "type": "integer"
                }
            }
        },
        "user.Card": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "integer"
                },
                "card_number": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "owner_id": {
                    "type": "string"
                }
            }
        },
        "user.PreLoginUser": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "phone_number": {
                    "type": "string"
                }
            }
        },
        "user.PreSignUpUser": {
            "type": "object",
            "properties": {
                "full_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone_number": {
                    "type": "string"
                }
            }
        },
        "user.User": {
            "type": "object",
            "properties": {
                "full_name": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_admin": {
                    "type": "boolean"
                },
                "password": {
                    "type": "string"
                },
                "phone_number": {
                    "type": "string"
                }
            }
        },
        "user.UserCard": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "integer"
                },
                "card_number": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Arzon-market API",
	Description:      "online meva va poliz mahsulotlari sotiladigan magazen APIsi",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
