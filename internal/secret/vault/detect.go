package vault

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	vcg "github.com/hashicorp/vault-client-go"
	"github.com/stretchr/objx"
)

const (
	mountInfoTypeKVString = "kv"
)

func DetectKVVersion(ctx context.Context, client *vcg.Client, mount string) (KVVersion, error) {
	engines, err := client.System.MountsListSecretsEngines(ctx)
	if err != nil {
		return KVVersionV1, fmt.Errorf("`MountsListSecretsEngines()` failed: %w", toPkgErr(err))
	}

	// if there are any warnings on the response, log them as warnings if the warn log level is enabled
	if len(engines.Warnings) > 0 && slog.Default().Enabled(context.Background(), slog.LevelWarn) {
		attrs := []slog.Attr{}
		for i, v := range engines.Warnings {
			attrs = append(attrs, slog.String(fmt.Sprintf("warning-%d", i), v))
		}
		// TODO: maybe group them under `warnings` key?
		slog.LogAttrs(ctx, slog.LevelWarn, "Calling `MountsListSecretEngines()` returned a response with warnings.", attrs...)
	}

	// `engines.Data` type is map[string]interface{}
	// Example:
	//   "legacy/" -> map[ ... options:<nil> type:kv ...]
	//   "something/else/" -> map[ ... options:map[version:2] type:kv ... ]
	//
	// Meaning that when KV engine is kv-v2 then options is non-nil and contains
	// key version and when KV engine is kv-v1 then options is nil.
	for k, v := range engines.Data {
		if !strings.HasPrefix(mount, k) {
			continue
		}

		o := objx.New(v)
		if !o.Get("type").IsStr() {
			slog.WarnContext(ctx, "Unexpected mount info structure, expected type of key `type` is string.", slog.String("mount", k), slog.Any("info", v))
			continue
		}

		if o.Get("type").Str() != mountInfoTypeKVString {
			continue
		}

		if o.Get("options.version").Str() == "2" {
			return KVVersionV2, nil
		} else {
			return KVVersionV1, nil
		}
	}
	return KVVersionV1, ErrNotFound
}
