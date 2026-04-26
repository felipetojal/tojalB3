package engine

import (
	"github.com/felipetojal/tojalB3/internal/metadata"
	"github.com/felipetojal/tojalB3/internal/volume"
)

// Engine represents the struct containing
// the engines of the
type Engine struct {
	v *volume.VolumeManager
	d *metadata.Database
}
