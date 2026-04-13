import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply, CortezaID } from '../../../cast'

const kind = 'FileUpload'

interface ExtraField {
  name: string;
  label: string;
}

interface Options {
  // Target module and file field
  moduleID: string;
  fileFieldName: string;

  // Button configuration
  buttonText: string;
  buttonVariant: string;

  // Multiple file handling:
  // 'newRecord' = one record per file (always available)
  // 'multiValue' = all files on one record (only when file field isMulti)
  multiFileMode: string;

  // Extra fields from the target module to display below the button
  extraFields: Array<ExtraField>;

  // Clear extra field values after successful upload
  clearOnSubmit: boolean;

  // Success message shown after upload
  successMessage: string;

  // Accepted file types (mime types)
  acceptedMimeTypes: string;
}

const defaults: Readonly<Options> = Object.freeze({
  moduleID: '',
  fileFieldName: '',
  buttonText: 'Upload File',
  buttonVariant: 'primary',
  multiFileMode: 'newRecord',
  extraFields: [],
  clearOnSubmit: true,
  successMessage: 'File uploaded successfully',
  acceptedMimeTypes: '*/*',
})

export class PageBlockFileUpload extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, CortezaID, 'moduleID')

    Apply(this.options, o, String,
      'fileFieldName',
      'buttonText',
      'buttonVariant',
      'multiFileMode',
      'successMessage',
      'acceptedMimeTypes',
    )

    Apply(this.options, o, Boolean, 'clearOnSubmit')

    if (o.extraFields) {
      this.options.extraFields = o.extraFields.map(f => ({ ...f }))
    }
  }

  validate (): Array<string> {
    const ee = super.validate()

    if (!this.options.moduleID) {
      ee.push('File upload block requires a target module')
    }

    if (!this.options.fileFieldName) {
      ee.push('File upload block requires a file field')
    }

    return ee
  }
}

Registry.set(kind, PageBlockFileUpload)
