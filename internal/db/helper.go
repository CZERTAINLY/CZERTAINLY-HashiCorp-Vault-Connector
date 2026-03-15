package db

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/config"

	"github.com/lib/pq"
)

func tbl(name string) string {
	return pq.QuoteIdentifier(config.Get().Database.Schema) +
		"." + pq.QuoteIdentifier(name)
}
