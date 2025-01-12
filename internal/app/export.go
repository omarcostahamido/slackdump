package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rusq/slackdump/v2/internal/export"
)

const defDirMode = 0700

// Export performs the full export of slack workspace in slack export compatible
// format.
func (app *App) Export(ctx context.Context, dir string) error {
	if dir == "" { // dir is passed from app.cfg.ExportDirectory
		return errors.New("export directory not specified")
	}

	if err := os.MkdirAll(dir, defDirMode); err != nil {
		return fmt.Errorf("Export: failed to create the export directory %q: %w", dir, err)
	}
	cfg := export.Options{
		Oldest:       time.Time(app.cfg.Oldest),
		Latest:       time.Time(app.cfg.Latest),
		IncludeFiles: app.cfg.Options.DumpFiles,
	}
	export := export.New(dir, app.sd, cfg)
	return export.Run(ctx)
}
