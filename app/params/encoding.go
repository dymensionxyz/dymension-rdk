package params

import (
	"github.com/ignite/cli/ignite/pkg/cosmoscmd"
)

// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	cosmoscmd.EncodingConfig
}
