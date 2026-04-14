package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/automation/html_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
)

var _ wfexec.ExecResponse

type (
	htmlHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h htmlHandler) register() {
	h.reg.AddFunctions(
		h.ToJson(),
		h.ToMarkdown(),
	)
}

type (
	htmlToJsonArgs struct {
		hasHtml bool
		Html    string
	}

	htmlToJsonResults struct {
		ResultJson string
	}
)

// ToJson function HTML to JSON
//
// expects implementation of toJson function:
// func (h htmlHandler) toJson(ctx context.Context, args *htmlToJsonArgs) (results *htmlToJsonResults, err error) {
//    return
// }
func (h htmlHandler) ToJson() *atypes.Function {
	return &atypes.Function{
		Ref:    "htmlToJson",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short:       "HTML to JSON",
			Description: "Converts HTML rich text into a structured JSON representation with sections, headings, paragraphs and list items",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "html",
				Types: []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "HTML content",
					Description: "Raw HTML string to parse",
				},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "resultJson",
				Types: []string{"String"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &htmlToJsonArgs{
					hasHtml: in.Has("html"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *htmlToJsonResults
			if results, err = h.toJson(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.ResultJson (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.ResultJson); err != nil {
					return
				} else if err = expr.Assign(out, "resultJson", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	htmlToMarkdownArgs struct {
		hasHtml bool
		Html    string
	}

	htmlToMarkdownResults struct {
		Markdown string
	}
)

// ToMarkdown function HTML to Markdown
//
// expects implementation of toMarkdown function:
// func (h htmlHandler) toMarkdown(ctx context.Context, args *htmlToMarkdownArgs) (results *htmlToMarkdownResults, err error) {
//    return
// }
func (h htmlHandler) ToMarkdown() *atypes.Function {
	return &atypes.Function{
		Ref:    "htmlToMarkdown",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short:       "HTML to Markdown",
			Description: "Converts HTML rich text into Markdown",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "html",
				Types: []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "HTML content",
					Description: "Raw HTML string to parse",
				},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "markdown",
				Types: []string{"String"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &htmlToMarkdownArgs{
					hasHtml: in.Has("html"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *htmlToMarkdownResults
			if results, err = h.toMarkdown(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Markdown (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.Markdown); err != nil {
					return
				} else if err = expr.Assign(out, "markdown", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
