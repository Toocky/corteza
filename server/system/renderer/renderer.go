package renderer

import (
	"context"
	"fmt"
	"io"

	"github.com/cortezaproject/corteza/server/pkg/options"
)

type (
	renderer struct {
		factories []driverFactory
	}
)

func Renderer(cfg options.TemplateOpt) *renderer {
	ff := make([]driverFactory, 0, 5)
	ff = append(ff, newGenericText(), newGenericHTML())

	// If a docxtemplater address is configured, use the sidecar for DOCX rendering.
	// Otherwise fall back to the in-process wordZero driver.
	if cfg.RendererDocxtemplaterAddress != "" {
		ff = append(ff, newDocxtemplaterDocx(cfg.RendererDocxtemplaterAddress))
	} else {
		ff = append(ff, newWordZeroDocx())
	}

	if cfg.RendererGotenbergEnabled {
		ff = append(ff, newGotenbergPDF(cfg.RendererGotenbergAddress))
	}

	return &renderer{
		factories: ff,
	}
}

func (r *renderer) Render(ctx context.Context, pl *RendererPayload) (io.ReadSeeker, error) {
	for _, f := range r.factories {
		if f.CanRender(pl.TemplateType) && f.CanProduce(pl.TargetType) {
			pp := make(map[string]io.Reader)
			for _, prt := range pl.Partials {
				pp[prt.Handle] = prt.Template
			}
			dpl := &driverPayload{
				Template:    pl.Template,
				Variables:   pl.Variables,
				Options:     pl.Options,
				Partials:    pp,
				Attachments: pl.Attachments,
			}

			return f.Driver().Render(ctx, dpl)
		}
	}

	return nil, fmt.Errorf("rendering failed: driver not found")
}

func (r *renderer) Drivers() []DriverDefinition {
	dd := make([]DriverDefinition, len(r.factories))
	for i, f := range r.factories {
		dd[i] = f.Define()
	}
	return dd
}
