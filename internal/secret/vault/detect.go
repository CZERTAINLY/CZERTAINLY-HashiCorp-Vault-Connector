package vault

import (
	"context"
	"fmt"
	"log/slog"

	vcg "github.com/hashicorp/vault-client-go"
	"github.com/stretchr/objx"
)

const (
	mountInfoTypeKVString = "kv"
)

func DetectKVVersion(ctx context.Context, client *vcg.Client, mount string) (KVVersion, error) {
	mounts, err := client.System.InternalUiListEnabledVisibleMounts(ctx)
	if err != nil {
		return KVVersionV1, fmt.Errorf("`InternalUiListEnabledVisibleMounts()` failed: %w", toPkgErr(err))
	}

	// if there are any warnings on the response, log them as warnings if the warn log level is enabled
	if len(mounts.Warnings) > 0 && slog.Default().Enabled(context.Background(), slog.LevelWarn) {
		attrs := []slog.Attr{}
		for i, v := range mounts.Warnings {
			attrs = append(attrs, slog.String(fmt.Sprintf("warning-%d", i), v))
		}
		// TODO: maybe group them under `warnings` key?
		slog.LogAttrs(ctx, slog.LevelWarn, "Calling `InternalUiListEnabledVisibleMounts()` returned a response with warnings.", attrs...)
	}

	for engineName, engineData := range mounts.Data.Secret {
		if mount != engineName {
			continue
		}

		o := objx.New(engineData)
		if !o.Get("type").IsStr() {
			slog.WarnContext(ctx, "Unexpected mount info structure, expected type of key `type` is string.", slog.String("mount", engineName), slog.Any("info", engineData))
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
