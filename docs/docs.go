// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-05-21 16:11:27.6606329 +0700 +07 m=+0.042009101

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
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/image": {
            "put": {
                "consumes": [
                    "multipart/form-data"
                ],
                "summary": "Upload an image",
                "operationId": "UploadImage",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Upload file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "File name",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorRes"
                        }
                    }
                }
            }
        },
        "/admin/image/{id}": {
            "get": {
                "summary": "Get an image information",
                "operationId": "GetImageByID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of image",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ImageInfoRes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorRes"
                        }
                    }
                }
            },
            "delete": {
                "summary": "Delete an image",
                "operationId": "DeleteImage",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of image",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorRes"
                        }
                    }
                }
            }
        },
        "/admin/image/{id}/rename": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "summary": "Rename an image",
                "operationId": "RenameImage",
                "parameters": [
                    {
                        "type": "string",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "ID of image",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorRes"
                        }
                    }
                }
            }
        },
        "/admin/image/{id}/replace": {
            "post": {
                "description": "replace and image",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Replace an image",
                "operationId": "ReplaceImage",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of image",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Replaced file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorRes"
                        }
                    }
                }
            }
        },
        "/admin/image/{id}/tag/{tag}": {
            "put": {
                "summary": "Add a tag to an image",
                "operationId": "AddImageTag",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of image",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Added tag",
                        "name": "tag",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorRes"
                        }
                    }
                }
            },
            "delete": {
                "summary": "Remove a tag from an image",
                "operationId": "RemoveImageTag",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of image",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Added tag",
                        "name": "tag",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorRes"
                        }
                    }
                }
            }
        },
        "/admin/images": {
            "get": {
                "description": "Get list of images information",
                "produces": [
                    "application/json"
                ],
                "summary": "Get list of images information",
                "operationId": "GetImages",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "pageCurrent",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "name": "orderBy",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "name": "orderDir",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "name": "tags",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.ImageInfoRes"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorRes"
                        }
                    }
                }
            }
        },
        "/images/size/{width}/{height}/{/name}": {
            "get": {
                "summary": "Get a resized image",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Width of image. Zero if resize scaled on its height",
                        "name": "width",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Height of image. Zero if resize scaled on its width",
                        "name": "height",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Image local path",
                        "name": "/name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorRes"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ErrorRes": {
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                }
            }
        },
        "models.ImageInfoRes": {
            "type": "object",
            "properties": {
                "fullname": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "models.ImageRenameReq": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "models.ImagesReq": {
            "type": "object",
            "properties": {
                "orderBy": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "orderDir": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "pageCurrent": {
                    "type": "integer"
                },
                "pageSize": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
        },
        "OAuth2AccessCode": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information"
            }
        },
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Implicit": {
            "type": "oauth2",
            "flow": "implicit",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Password": {
            "type": "oauth2",
            "flow": "password",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "read": " Grants read access",
                "write": " Grants write access"
            }
        }
    },
    "x-extension-openapi": {
        "example": "value on a json format"
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
	Version:     "1.0",
	Host:        "localhost:5000",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "Swagger Example API",
	Description: "This is a sample server celler server.",
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
