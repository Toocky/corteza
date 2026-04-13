<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      v-if="!isConfigured"
      class="p-3 text-center text-muted"
    >
      {{ $t('fileUpload.notConfigured') }}
    </div>

    <div
      v-else
      class="p-3"
    >
      <!-- Extra fields rendered above the button -->
      <template v-if="tempRecord && configuredFields.length">
        <div
          v-for="field in configuredFields"
          :key="field.fieldID"
          class="mb-2"
        >
          <field-editor
            :namespace="namespace"
            :field="field"
            :record="tempRecord"
            :errors="fieldErrors(field.name)"
          />
        </div>
      </template>

      <!-- Upload button -->
      <b-button
        :variant="options.buttonVariant || 'primary'"
        :disabled="uploading"
        class="w-100"
        @click="triggerFileSelect"
      >
        <b-spinner
          v-if="uploading"
          small
          class="mr-1"
        />
        {{ uploading ? $t('fileUpload.uploading') : (options.buttonText || $t('fileUpload.defaultButtonText')) }}
      </b-button>

      <!-- Hidden CUploader — handles the actual upload exactly like the native file field -->
      <div class="d-none">
        <c-uploader
          ref="uploader"
          :endpoint="uploadEndpoint"
          :accepted-files="acceptedFiles"
          :form-data="uploaderFormData"
          @upload="onUploadSuccess"
          @error="onUploadError"
        />
      </div>

      <!-- Hidden file input for native OS picker -->
      <input
        ref="fileInput"
        type="file"
        :accept="options.acceptedMimeTypes || '*/*'"
        multiple
        class="d-none"
        @change="onFilesSelected"
      >
    </div>
  </wrap>
</template>

<script>
import base from './base'
import FieldEditor from 'corteza-webapp-compose/src/components/ModuleFields/Editor'
import { compose, NoID, validator } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
import { mapGetters } from 'vuex'

const { CUploader } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldEditor,
    CUploader,
  },

  extends: base,

  data () {
    return {
      uploading: false,
      tempRecord: null,
      validated: new validator.Validated(),
      pendingFiles: [],
      uploadedAttachments: [],
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
    }),

    isConfigured () {
      return this.options.moduleID &&
        this.options.moduleID !== NoID &&
        this.options.fileFieldName
    },

    targetModule () {
      if (!this.options.moduleID || this.options.moduleID === NoID) return null
      return this.getModuleByID(this.options.moduleID)
    },

    configuredFields () {
      if (!this.targetModule || !this.options.extraFields) return []

      return this.options.extraFields
        .map(ef => this.targetModule.fields.find(f => f.name === ef.name))
        .filter(f => !!f)
    },

    uploadEndpoint () {
      const { moduleID, fileFieldName } = this.options
      const namespaceID = this.namespace.namespaceID

      return this.$ComposeAPI.baseURL + this.$ComposeAPI.recordUploadEndpoint({
        namespaceID,
        moduleID,
      })
    },

    uploaderFormData () {
      return {
        fieldName: this.options.fileFieldName,
      }
    },

    acceptedFiles () {
      const mt = (this.options.acceptedMimeTypes || '').trim()
      if (!mt || mt === '*/*') {
        return ['*/*']
      }
      return mt.split(',').map(s => s.trim())
    },
  },

  watch: {
    targetModule: {
      immediate: true,
      handler (mod) {
        if (mod) {
          this.initTempRecord(mod)
        }
      },
    },
  },

  methods: {
    initTempRecord (mod) {
      this.tempRecord = new compose.Record(mod, {})
    },

    fieldErrors (fieldName) {
      return this.validated.filterByMeta('field', fieldName)
    },

    triggerFileSelect () {
      if (this.uploading) return
      this.$refs.fileInput.value = ''
      this.$refs.fileInput.click()
    },

    onFilesSelected (e) {
      const files = Array.from(e.target.files || [])
      if (!files.length) return

      this.uploading = true
      this.pendingFiles = [...files]
      this.uploadedAttachments = []

      // Feed files into the hidden CUploader one at a time
      // CUploader uses Dropzone which processes them sequentially
      const dropzone = this.$refs.uploader.$refs.dropzone
      for (const file of files) {
        dropzone.addFile(file)
      }
    },

    async onUploadSuccess (response) {
      this.uploadedAttachments.push(response)

      // Check if all files have been uploaded
      if (this.uploadedAttachments.length < this.pendingFiles.length) {
        return
      }

      // All uploads done — now create the record(s)
      const { moduleID, fileFieldName, multiFileMode, successMessage, clearOnSubmit } = this.options
      const namespaceID = this.namespace.namespaceID

      try {
        const attachmentIDs = this.uploadedAttachments.map(a => a.attachmentID)

        if (multiFileMode === 'multiValue') {
          // All files on one record
          const values = this.buildRecordValues()
          values[fileFieldName] = attachmentIDs

          await this.$ComposeAPI.recordCreate({
            namespaceID,
            moduleID,
            values: this.formatValues(values),
          })
        } else {
          // New record per file
          for (const attID of attachmentIDs) {
            const values = this.buildRecordValues()
            values[fileFieldName] = [attID]

            await this.$ComposeAPI.recordCreate({
              namespaceID,
              moduleID,
              values: this.formatValues(values),
            })
          }
        }

        this.$root.$bvToast.toast(successMessage || this.$t('fileUpload.defaultSuccessMessage'), {
          variant: 'success',
          solid: true,
          toaster: 'b-toaster-bottom-right',
        })

        if (clearOnSubmit) {
          this.initTempRecord(this.targetModule)
        }
      } catch (err) {
        console.error('FileUpload record create error:', err)
        this.$root.$bvToast.toast(err.message || this.$t('fileUpload.uploadFailed'), {
          title: this.$t('fileUpload.errorTitle'),
          variant: 'danger',
          solid: true,
          toaster: 'b-toaster-bottom-right',
        })
      } finally {
        this.uploading = false
        this.pendingFiles = []
        this.uploadedAttachments = []
      }
    },

    onUploadError (e, message) {
      console.error('FileUpload upload error:', message)
      this.uploading = false
      this.pendingFiles = []
      this.uploadedAttachments = []

      this.$root.$bvToast.toast(message || this.$t('fileUpload.uploadFailed'), {
        title: this.$t('fileUpload.errorTitle'),
        variant: 'danger',
        solid: true,
        toaster: 'b-toaster-bottom-right',
      })
    },

    buildRecordValues () {
      const values = {}
      if (!this.tempRecord || !this.configuredFields.length) return values

      for (const field of this.configuredFields) {
        const val = this.tempRecord.values[field.name]
        if (val !== undefined && val !== null) {
          values[field.name] = Array.isArray(val) ? val : [val]
        }
      }

      return values
    },

    formatValues (values) {
      const formatted = []
      for (const [name, value] of Object.entries(values)) {
        if (Array.isArray(value)) {
          value.forEach(v => {
            formatted.push({ name, value: (v !== undefined && v !== null) ? String(v) : '' })
          })
        } else {
          formatted.push({ name, value: (value !== undefined && value !== null) ? String(value) : '' })
        }
      }
      return formatted
    },
  },
}
</script>
