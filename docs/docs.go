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
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/polls": {
            "get": {
                "description": "list of polls",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Poll"
                ],
                "summary": "list of polls",
                "operationId": "PollList",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "PageSize",
                        "name": "page_size",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Tag",
                        "name": "tag",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "User-ID",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.ResponseSuccess"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/poll.SinglePollList"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ResponseFailure"
                        }
                    }
                }
            },
            "post": {
                "description": "create poll",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Poll"
                ],
                "summary": "create poll",
                "operationId": "CreatePoll",
                "parameters": [
                    {
                        "description": "Create-Poll-Params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.CreatePollReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ResponseFailure"
                        }
                    }
                }
            }
        },
        "/api/v1/polls/{id}/skip": {
            "post": {
                "description": "skip poll",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Poll"
                ],
                "summary": "skip poll",
                "operationId": "Skip",
                "parameters": [
                    {
                        "description": "Vote-Params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.VoteReq"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "Vote ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ResponseFailure"
                        }
                    }
                }
            }
        },
        "/api/v1/polls/{id}/stats": {
            "get": {
                "description": "list of poll stats",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "PollStats"
                ],
                "summary": "list of poll stats",
                "operationId": "PollStats",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Vote ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.ResponseSuccess"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/statistic.StatsResult"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ResponseFailure"
                        }
                    }
                }
            }
        },
        "/api/v1/polls/{id}/vote": {
            "post": {
                "description": "vote poll",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Poll"
                ],
                "summary": "vote to poll",
                "operationId": "Vote",
                "parameters": [
                    {
                        "description": "Vote-Params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.VoteReq"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "Vote ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ResponseFailure"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.CreatePollReq": {
            "type": "object",
            "required": [
                "options",
                "title"
            ],
            "properties": {
                "options": {
                    "type": "array",
                    "minItems": 2,
                    "items": {
                        "type": "string"
                    }
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "api.ErrorCode": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 404
                },
                "message": {
                    "type": "string",
                    "example": "item not found"
                }
            }
        },
        "api.ResponseFailure": {
            "type": "object",
            "properties": {
                "error": {
                    "$ref": "#/definitions/api.ErrorCode"
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "api.ResponseSuccess": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "api.VoteReq": {
            "type": "object",
            "required": [
                "user_id"
            ],
            "properties": {
                "option_index": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "poll.SinglePollList": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "options": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "statistic.StatsResult": {
            "type": "object",
            "properties": {
                "pollID": {
                    "type": "integer"
                },
                "votes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/statistic.StatsVote"
                    }
                }
            }
        },
        "statistic.StatsVote": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "option": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
