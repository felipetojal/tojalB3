package engine

import (
	"github.com/felipetojal/tojalB3/internal/metadata"
	"github.com/felipetojal/tojalB3/internal/volume"
)

type Engine struct {
	v *volume.VolumeManager
	d *metadata.Database
}
