// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
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
        "/login": {
            "post": {
                "description": "Validates user credentials and logs the user in if the credentials are correct.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "User Login",
                "parameters": [
                    {
                        "description": "User credentials (username and password)",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success message indicating user login",
                        "schema": {
                            "$ref": "#/definitions/entity.MsgResponse"
                        }
                    },
                    "400": {
                        "description": "Error message indicating invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Error message indicating invalid JSON format",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Registers a new user by saving the provided user credentials to the database.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "User Registration",
                "parameters": [
                    {
                        "description": "User credentials (username, password, etc.)",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success message indicating successful registration",
                        "schema": {
                            "$ref": "#/definitions/entity.MsgResponse"
                        }
                    },
                    "400": {
                        "description": "Error message indicating invalid username or password",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Error message indicating invalid JSON format",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Retrieves all users from the database.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get Users",
                "responses": {
                    "200": {
                        "description": "List of users in the database",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.User"
                            }
                        }
                    },
                    "422": {
                        "description": "Error message indicating no users found",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error message indicating internal server error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Invalid Username and Password"
                }
            }
        },
        "entity.MsgResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "user123 login successfully"
                }
            }
        },
        "entity.User": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
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
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Tag Example Monitoring Service",
	Description:      "A Tag service API in Go using Gin framework",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
