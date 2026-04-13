package provision

import (
	"bytes"
	"context"
	_ "embed"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/objstore"
	"github.com/cortezaproject/corteza/server/store"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
)

//go:embed cv_template.docx
var cvTemplateDocx []byte

const cvTemplateHandle = "docx_cv"

// seedDocxCVTemplate ensures the default CV docx template is present.
// It is idempotent — if the template already exists (by handle), it is left unchanged.
func SeedDocxCVTemplate(ctx context.Context, log *zap.Logger, s store.Storer, files objstore.Store) error {
	log = log.Named("docx_cv_template")

	existing, err := store.LookupTemplateByHandle(ctx, s, cvTemplateHandle)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if existing != nil {
		log.Debug("default CV docx template already exists, skipping")
		return nil
	}

	log.Info("seeding default CV docx template")

	attID := id.Next()
	attExt := "docx"
	attURL := files.Original(attID, attExt)

	if err = files.Save(attURL, bytes.NewReader(cvTemplateDocx)); err != nil {
		return err
	}

	att := &sysTypes.Attachment{
		ID:   attID,
		Name: "cv_template.docx",
		Kind: sysTypes.AttachmentKindTemplate,
		Url:  attURL,
	}
	att.Meta.Original.Extension = attExt
	att.Meta.Original.Size = int64(len(cvTemplateDocx))
	att.Meta.Original.Mimetype = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	att.CreatedAt = time.Now().Round(time.Second)

	if err = store.CreateAttachment(ctx, s, att); err != nil {
		return err
	}

	tplID := id.Next()
	tpl := &sysTypes.Template{
		ID:           tplID,
		Handle:       cvTemplateHandle,
		Language:     "",
		Type:         sysTypes.DocumentTypeDocx,
		Partial:      false,
		SourceFileID: attID,
		CreatedAt:    time.Now().Round(time.Second),
	}
	tpl.Meta.Short = "CV Template (DOCX)"
	tpl.Meta.Description = "Default CV template using wordZero syntax. Upload a new .docx to replace."

	return store.CreateTemplate(ctx, s, tpl)
}
