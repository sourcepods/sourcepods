// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http",
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "description": "This is the API for GitPods - git in the cloud.",
    "title": "GitPods OpenAPI",
    "license": {
      "name": "Apache-2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "1.0.0"
  },
  "basePath": "/v1",
  "paths": {
    "/users": {
      "get": {
        "tags": [
          "users"
        ],
        "summary": "List all users",
        "operationId": "listUsers",
        "responses": {
          "200": {
            "description": "An array of all users",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/user"
              }
            }
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/users/{username}": {
      "get": {
        "tags": [
          "users"
        ],
        "summary": "Get a user by their username",
        "operationId": "getUser",
        "parameters": [
          {
            "type": "string",
            "description": "The username of a user",
            "name": "username",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The user by their username",
            "schema": {
              "$ref": "#/definitions/user"
            }
          },
          "404": {
            "description": "The user is not found by their username",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      },
      "patch": {
        "tags": [
          "users"
        ],
        "summary": "Update the user's information",
        "operationId": "updateUser",
        "parameters": [
          {
            "type": "string",
            "description": "The username of the user to update",
            "name": "username",
            "in": "path",
            "required": true
          },
          {
            "description": "The updated user",
            "name": "updatedUser",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "name"
              ],
              "properties": {
                "name": {
                  "type": "string"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "The user has been updated",
            "schema": {
              "$ref": "#/definitions/user"
            }
          },
          "404": {
            "description": "The user could not be found by this username",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "422": {
            "description": "The updated user has invalid input",
            "schema": {
              "$ref": "#/definitions/validationError"
            }
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "user": {
      "type": "object",
      "required": [
        "id",
        "username"
      ],
      "properties": {
        "created_at": {
          "type": "string",
          "format": "date-time"
        },
        "email": {
          "type": "string",
          "format": "email"
        },
        "id": {
          "type": "string",
          "format": "uuid",
          "readOnly": true
        },
        "name": {
          "type": "string"
        },
        "updated_at": {
          "type": "string",
          "format": "date-time"
        },
        "username": {
          "type": "string"
        }
      }
    },
    "validationError": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "errors": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "field": {
                "type": "string"
              },
              "message": {
                "type": "string"
              }
            }
          }
        },
        "message": {
          "type": "string"
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http",
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "description": "This is the API for GitPods - git in the cloud.",
    "title": "GitPods OpenAPI",
    "license": {
      "name": "Apache-2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "1.0.0"
  },
  "basePath": "/v1",
  "paths": {
    "/users": {
      "get": {
        "tags": [
          "users"
        ],
        "summary": "List all users",
        "operationId": "listUsers",
        "responses": {
          "200": {
            "description": "An array of all users",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/user"
              }
            }
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/users/{username}": {
      "get": {
        "tags": [
          "users"
        ],
        "summary": "Get a user by their username",
        "operationId": "getUser",
        "parameters": [
          {
            "type": "string",
            "description": "The username of a user",
            "name": "username",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The user by their username",
            "schema": {
              "$ref": "#/definitions/user"
            }
          },
          "404": {
            "description": "The user is not found by their username",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      },
      "patch": {
        "tags": [
          "users"
        ],
        "summary": "Update the user's information",
        "operationId": "updateUser",
        "parameters": [
          {
            "type": "string",
            "description": "The username of the user to update",
            "name": "username",
            "in": "path",
            "required": true
          },
          {
            "description": "The updated user",
            "name": "updatedUser",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "name"
              ],
              "properties": {
                "name": {
                  "type": "string"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "The user has been updated",
            "schema": {
              "$ref": "#/definitions/user"
            }
          },
          "404": {
            "description": "The user could not be found by this username",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "422": {
            "description": "The updated user has invalid input",
            "schema": {
              "$ref": "#/definitions/validationError"
            }
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "user": {
      "type": "object",
      "required": [
        "id",
        "username"
      ],
      "properties": {
        "created_at": {
          "type": "string",
          "format": "date-time"
        },
        "email": {
          "type": "string",
          "format": "email"
        },
        "id": {
          "type": "string",
          "format": "uuid",
          "readOnly": true
        },
        "name": {
          "type": "string"
        },
        "updated_at": {
          "type": "string",
          "format": "date-time"
        },
        "username": {
          "type": "string"
        }
      }
    },
    "validationError": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "errors": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "field": {
                "type": "string"
              },
              "message": {
                "type": "string"
              }
            }
          }
        },
        "message": {
          "type": "string"
        }
      }
    }
  }
}`))
}
