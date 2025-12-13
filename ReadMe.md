# OpenApi-HTTP

Create rfc9110 compliant http requests from an Open API spec. These can be used in `.http` files that can document API usage and help with developer experience.

### Installation

Install using go: 

```sh
go install github.com/kalli/openapi-http
```

Or clone this repo and build

```sh
git clone git@github.com:kalli/openapi-http.git
```

## Usage 

See the available flags and parameters:

```sh
openapi-http --help
```

List the operations available in an Open Api spec:

```sh
# Use a URL:
openapi-http https://petstore3.swagger.io/api/v3/openapi.json
# Or a filepath
openapi-http test/petstore.yml
```

Create an example request:

```sh
openapi-http https://petstore3.swagger.io/api/v3/openapi.json --operation-id addPet
# or use -i for operationId
openapi-http https://petstore3.swagger.io/api/v3/openapi.json -i addPet
# or positional
openapi-http https://petstore3.swagger.io/api/v3/openapi.json addPet
```

This will create an example request like so:

```http
###
# @name addPet
# Add a new pet to the store

POST http://petstore.swagger.io/v2/pet
Content-Type: application/json

{
  "category": {
    "id": 0,
    "name": "string"
  },
  "id": 0,
  "name": "doggie",
  "photoUrls": [
    "string"
  ],
  "status": "available",
  "tags": [
    {
      "id": 0,
      "name": "string"
    }
  ]
}
```

You can also list all the operations for a given path:

```sh
openapi-http https://petstore3.swagger.io/api/v3/openapi.json -path /pet 
```

This will print out requests for the `updatePet` and `addPet` requests.