/*
 * CZERTAINLY Discovery Provider API
 *
 * REST API for implementations of custom Discovery Provider
 *
 * API version: 2.11.0
 * Contact: getinfo@czertainly.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package discovery

import (
	"log"
	"net/http"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}