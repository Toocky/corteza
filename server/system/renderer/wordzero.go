package renderer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/mr-pmillz/wordZero/pkg/document"
	"github.com/cortezaproject/corteza/server/system/types"
)

func init() {
	document.SetGlobalLevel(document.LogLevelSilent)
}

type (
	wordZeroDocx struct {
		def DriverDefinition
	}
	wordZeroDocxDriver struct{}
)

func newWordZeroDocx() driverFactory {
	return &wordZeroDocx{
		def: DriverDefinition{
			Name: "wordZeroDocx",
			InputTypes: []types.DocumentType{
				types.DocumentTypeDocx,
			},
			OutputTypes: []types.DocumentType{
				types.DocumentTypeDocx,
			},
		},
	}
}

func (d *wordZeroDocx) Define() DriverDefinition {
	return d.def
}

func (d *wordZeroDocx) CanRender(t types.DocumentType) bool {
	for _, i := range d.def.InputTypes {
		if i == t {
			return true
		}
	}
	return false
}

func (d *wordZeroDocx) CanProduce(t types.DocumentType) bool {
	for _, o := range d.def.OutputTypes {
		if o == t {
			return true
		}
	}
	return false
}

func (d *wordZeroDocx) Driver() driver {
	return &wordZeroDocxDriver{}
}

func (d *wordZeroDocxDriver) Render(ctx context.Context, pl *driverPayload) (io.ReadSeeker, error) {
	templateBytes, err := io.ReadAll(pl.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to read docx template: %w", err)
	}

	doc, err := document.OpenFromMemory(io.NopCloser(bytes.NewReader(templateBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to open docx template: %w", err)
	}

	// Merge template placeholders that Word split across multiple runs
	consolidateRuns(doc)

	engine := document.NewTemplateEngine()

	_, err = engine.LoadTemplateFromDocument("tpl", doc)
	if err != nil {
		return nil, fmt.Errorf("failed to load docx template: %w", err)
	}

	data := buildTemplateData(pl.Variables)

	rendered, err := engine.RenderTemplateToDocument("tpl", data)
	if err != nil {
		return nil, fmt.Errorf("failed to render docx template: %w", err)
	}

	out, err := rendered.ToBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize rendered docx: %w", err)
	}

	return bytes.NewReader(out), nil
}

// consolidateRuns walks all paragraphs in the document (including inside tables)
// and merges runs where template placeholders ({{ }}) have been split across
// multiple runs by Word. Operates on the parsed structs, not raw XML.
func consolidateRuns(doc *document.Document) {
	for _, elem := range doc.Body.Elements {
		switch e := elem.(type) {
		case *document.Paragraph:
			consolidateParagraphRuns(e)
		case *document.Table:
			consolidateTableRuns(e)
		}
	}
}

func consolidateTableRuns(t *document.Table) {
	for ri := range t.Rows {
		for ci := range t.Rows[ri].Cells {
			for pi := range t.Rows[ri].Cells[ci].Paragraphs {
				consolidateParagraphRuns(&t.Rows[ri].Cells[ci].Paragraphs[pi])
			}
		}
	}
}

func consolidateParagraphRuns(p *document.Paragraph) {
	if len(p.Runs) < 2 {
		return
	}

	i := 0
	for i < len(p.Runs) {
		text := p.Runs[i].Text.Content
		opens := strings.Count(text, "{{")
		closes := strings.Count(text, "}}")

		if opens <= closes {
			i++
			continue
		}

		// Unmatched {{ — pull text from subsequent runs until balanced
		j := i + 1
		for j < len(p.Runs) && opens > closes {
			nextText := p.Runs[j].Text.Content
			opens += strings.Count(nextText, "{{")
			closes += strings.Count(nextText, "}}")
			p.Runs[i].Text.Content += nextText
			p.Runs[j].Text.Content = ""
			j++
		}

		p.Runs[i].Text.Space = "preserve"
		i = j
	}
}

func buildTemplateData(vars map[string]interface{}) *document.TemplateData {
	data := document.NewTemplateData()
	flattenVariables("", vars, data)
	return data
}

// flattenVariables walks the variables map recursively.
// Scalar values are flattened to camelCase keys (e.g. profile.contact.location -> profileContactLocation).
// Slices are passed as lists for {{#each}} loops. Boolean conditions (hasX) are auto-derived from non-empty slices.
func flattenVariables(prefix string, vars map[string]interface{}, data *document.TemplateData) {
	for key, val := range vars {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + capitalizeFirst(key)
		}

		switch v := val.(type) {
		case map[string]interface{}:
			flattenVariables(fullKey, v, data)

		case []interface{}:
			data.SetList(fullKey, v)
			data.SetCondition("has"+capitalizeFirst(fullKey), len(v) > 0)

		case bool:
			data.SetCondition(fullKey, v)
			data.SetVariable(fullKey, fmt.Sprintf("%v", v))

		default:
			data.SetVariable(fullKey, val)
		}
	}
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
