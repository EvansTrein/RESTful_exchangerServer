// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Evans Trein",
            "url": "https://github.com/EvansTrein",
            "email": "evanstrein@icloud.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/balance": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get the balance of all accounts",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wallet"
                ],
                "summary": "Get user balance",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.BalanceResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "504": {
                        "description": "Gateway Timeout",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    }
                }
            }
        },
        "/delete": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "user delete",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Delete",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "504": {
                        "description": "Gateway Timeout",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    }
                }
            }
        },
        "/exchange": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Exchange one currency to another for the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wallet"
                ],
                "summary": "Exchange currency",
                "parameters": [
                    {
                        "description": "Exchange request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ExchangeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ExchangeResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "402": {
                        "description": "Payment Required",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "504": {
                        "description": "Gateway Timeout",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    }
                }
            }
        },
        "/exchange/rates": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get the current exchange rates for supported currencies",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wallet"
                ],
                "summary": "Get all exchange rates",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ExchangeRatesResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "504": {
                        "description": "Gateway Timeout",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "user login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "504": {
                        "description": "Gateway Timeout",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Creating a new user with the provided data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Creating a new user",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.RegisterResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "504": {
                        "description": "Gateway Timeout",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    }
                }
            }
        },
        "/wallet/deposit": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Deposit funds into a user's account for a specific currency",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wallet"
                ],
                "summary": "Deposit funds into an account",
                "parameters": [
                    {
                        "description": "Deposit request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AccountOperationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AccountOperationResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "504": {
                        "description": "Gateway Timeout",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    }
                }
            }
        },
        "/wallet/withdraw": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Withdraw funds from a user's account for a specific currency",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wallet"
                ],
                "summary": "Withdraw funds from an account",
                "parameters": [
                    {
                        "description": "Withdraw request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AccountOperationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AccountOperationResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "402": {
                        "description": "Payment Required",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    },
                    "504": {
                        "description": "Gateway Timeout",
                        "schema": {
                            "$ref": "#/definitions/models.HandlerResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AccountOperationRequest": {
            "type": "object",
            "required": [
                "amount",
                "currency"
            ],
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 2000
                },
                "currency": {
                    "type": "string",
                    "maxLength": 6,
                    "minLength": 3,
                    "example": "USD"
                }
            }
        },
        "models.AccountOperationResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "text message"
                },
                "new_balance": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "number"
                    }
                }
            }
        },
        "models.BalanceResponse": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "number"
                    }
                }
            }
        },
        "models.ExchangeRatesResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "text message"
                },
                "rates": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "number"
                    }
                }
            }
        },
        "models.ExchangeRequest": {
            "type": "object",
            "required": [
                "amount",
                "from_currency",
                "to_currency"
            ],
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 500
                },
                "from_currency": {
                    "type": "string",
                    "maxLength": 6,
                    "minLength": 3,
                    "example": "USD"
                },
                "to_currency": {
                    "type": "string",
                    "maxLength": 6,
                    "minLength": 3,
                    "example": "CNY"
                }
            }
        },
        "models.ExchangeResponse": {
            "type": "object",
            "properties": {
                "exchange_rate": {
                    "type": "number",
                    "example": 7.424683
                },
                "message": {
                    "type": "string",
                    "example": "text message"
                },
                "new_balance": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "number"
                    }
                },
                "received_account": {
                    "$ref": "#/definitions/models.ReceivedAccount"
                },
                "spent_accoutn": {
                    "$ref": "#/definitions/models.SpentAccoutn"
                }
            }
        },
        "models.HandlerResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "text error"
                },
                "message": {
                    "type": "string",
                    "example": "text message"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "models.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john.doe@example.com"
                },
                "password": {
                    "type": "string",
                    "example": "123456"
                }
            }
        },
        "models.LoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string",
                    "example": "JWT-token"
                }
            }
        },
        "models.ReceivedAccount": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 3636.3
                },
                "currency": {
                    "type": "string",
                    "example": "CNY"
                }
            }
        },
        "models.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john.doe@example.com"
                },
                "password": {
                    "type": "string",
                    "minLength": 6,
                    "example": "123456"
                },
                "username": {
                    "type": "string",
                    "minLength": 3,
                    "example": "john"
                }
            }
        },
        "models.RegisterResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "message": {
                    "type": "string",
                    "example": "user successfully created"
                }
            }
        },
        "models.SpentAccoutn": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 500
                },
                "currency": {
                    "type": "string",
                    "example": "USD"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "Type \"Bearer\" followed by a space and the token. Example: \"Bearer your_token\"",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "Currency exchanger",
	Description:      "REST API that works with - postgres as a database, a third-party gRPC server (for currency currencies)\nand Redis for caching responses from a third-party gRPC service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
