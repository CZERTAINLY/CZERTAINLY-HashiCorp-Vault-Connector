#!/bin/sh

czertainlyHome="/opt/czertainly"
source ${czertainlyHome}/static-functions

log "INFO" "Launching Hashicorp Vault Connector"
./appbin

#exec "$@"