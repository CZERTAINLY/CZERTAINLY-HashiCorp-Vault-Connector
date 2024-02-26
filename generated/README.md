# Go API Server for openapi

REST API for implementations of custom Discovery Provider

## Overview
This server was generated by the [openapi-generator]
(https://openapi-generator.tech) project.
By using the [OpenAPI-Spec](https://github.com/OAI/OpenAPI-Specification) from a remote server, you can easily generate a server stub.
-

To see how to make this your own, look here:

[README](https://openapi-generator.tech)

- API version: 2.11.0
- Build date: 2024-02-21T14:22:29.123341897+01:00[Europe/Prague]
For more information, please visit [https://www.czertainly.com](https://www.czertainly.com)


### Running the server
To run the server, follow these simple steps:

```
go run main.go
```

To run the server in a docker container
```
docker build --network=host -t openapi .
```

Once image is built use
```
docker run --rm -it openapi
```