package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

template: schema.#optionsGroup & {
	handle: "template"
	title:  "Rendering engine"

	options: {
		renderer_gotenberg_address: {
			defaultGoExpr: ""
			description:   "Gotenberg rendering container address."
		}

		renderer_gotenberg_enabled: {
			type:        "bool"
			description: "Is Gotenberg rendering container enabled."
		}

		renderer_docxtemplater_address: {
			defaultGoExpr: ""
			description:   "docxtemplater sidecar address. When set, DOCX rendering is delegated to this sidecar instead of the in-process wordZero driver."
		}

		renderer_docxtemplater_enabled: {
			type:        "bool"
			description: "Is docxtemplater sidecar enabled."
		}
	}
}
