package renderer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	docxtemplaterDocx struct {
		url string
		def DriverDefinition
	}
	docxtemplaterDocxDriver struct {
		url string
	}
)

func newDocxtemplaterDocx(url string) driverFactory {
	return &docxtemplaterDocx{
		url: url,
		def: DriverDefinition{
			Name: "docxtemplaterDocx",
			InputTypes: []types.DocumentType{
				types.DocumentTypeDocx,
			},
			OutputTypes: []types.DocumentType{
				types.DocumentTypeDocx,
			},
		},
	}
}

func (d *docxtemplaterDocx) Define() DriverDefinition {
	return d.def
}

func (d *docxtemplaterDocx) CanRender(t types.DocumentType) bool {
	for _, i := range d.def.InputTypes {
		if i == t {
			return true
		}
	}
	return false
}

func (d *docxtemplaterDocx) CanProduce(t types.DocumentType) bool {
	for _, o := range d.def.OutputTypes {
		if o == t {
			return true
		}
	}
	return false
}

func (d *docxtemplaterDocx) Driver() driver {
	return &docxtemplaterDocxDriver{url: d.url}
}

func (d *docxtemplaterDocxDriver) Render(ctx context.Context, pl *driverPayload) (io.ReadSeeker, error) {
	templateBytes, err := io.ReadAll(pl.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to read docx template: %w", err)
	}

	// Build multipart request: template file + data JSON
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add template file
	part, err := writer.CreateFormFile("template", "template.docx")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err = part.Write(templateBytes); err != nil {
		return nil, fmt.Errorf("failed to write template: %w", err)
	}

	// Add data as JSON string
	dataJSON, err := json.Marshal(pl.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal variables: %w", err)
	}
	if err = writer.WriteField("data", string(dataJSON)); err != nil {
		return nil, fmt.Errorf("failed to write data field: %w", err)
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", d.url+"/render", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("docxtemplater service unavailable: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("docxtemplater render failed (%d): %s", resp.StatusCode, string(respBytes))
	}

	return bytes.NewReader(respBytes), nil
}
