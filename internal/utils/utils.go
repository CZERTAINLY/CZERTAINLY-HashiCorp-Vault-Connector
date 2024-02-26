package utils

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/generated/discovery"
	"crypto/md5"
	"encoding/hex"
	"log"

	"github.com/google/uuid"
)

func DeterministicGUID(parts ...string) string {
	// concatenate all strings
	var combined string
	for _, part := range parts {
		combined += part
	}

	md5hash := md5.New()
	md5hash.Write([]byte(combined))

	// convert the hash value to a string
	md5string := hex.EncodeToString(md5hash.Sum(nil))

	// generate the UUID from the
	// first 16 bytes of the MD5 hash
	uuid, err := uuid.FromBytes([]byte(md5string[0:16]))
	if err != nil {
		log.Fatal(err)
	}

	return uuid.String()
}

func ConvertRoutesToEndpointDto(routes discovery.Routes) []discovery.EndpointDto {
	endpoints := make([]discovery.EndpointDto, 0, len(routes))
	for name, route := range routes {
		endpoint := discovery.EndpointDto{
			Name:     name,
			Method:   route.Method,
			Uuid:     DeterministicGUID(route.Method + route.Pattern),
			Context:  route.Pattern,
			Required: true,
		}
		endpoints = append(endpoints, endpoint)
	}
	return endpoints
}
