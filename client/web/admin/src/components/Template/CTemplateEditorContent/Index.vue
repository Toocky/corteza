<template>
  <b-row>
    <!-- Toolbox -->
    <b-col
      cols="12"
      lg="3"
      class="mb-3"
    >
      <editor-toolbox
        :template="template"
        :partials="partials"
      />
    </b-col>

    <!-- Preview configuration -->
    <b-col
      cols="12"
      lg="9"
      class="mb-3"
    >
      <b-card
        v-if="!template.partial"
        body-class="p-0"
        header-class="border-bottom"
        footer-class="border-top d-flex justify-content-end flex-wrap flex-fill-child gap-1"
      >
        <!-- Partial templates can't be previewed -->
        <c-ace-editor
          v-model="previewData"
          data-test-id="template-preview-output"
          name="preview-data"
          lang="json"
          min-height="300px"
          show-line-numbers
          highlight-active-line
          resizable
          :border="false"
        />

        <template #header>
          <h4 class="m-0">
            {{ $t('preview.title') }}
          </h4>
        </template>

        <template #footer>
          <b-btn
            v-if="canPreviewHTML"
            data-test-id="button-preview-html-template"
            variant="light"
            @click="openPreview('html')"
          >
            {{ $t('preview.html') }}
          </b-btn>

          <b-btn
            v-if="canPreviewPDF"
            data-test-id="button-preview-pdf-template"
            variant="light"
            @click="openPreview('pdf')"
          >
            {{ $t('preview.pdf') }}
          </b-btn>

          <b-btn
            v-if="canPreviewDocx"
            data-test-id="button-preview-docx-template"
            variant="light"
            @click="downloadPreview('docx')"
          >
            {{ $t('preview.docx') }}
          </b-btn>
        </template>
      </b-card>
    </b-col>

    <!-- Content editor -->
    <b-col cols="12">
      <b-card
        body-class="p-0"
        header-class="d-flex align-items-center border-bottom"
        footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
        class="shadow-sm"
      >
        <component
          :is="editor"
          :template="template"
        />

        <template #header>
          <h4 class="m-0">
            {{ $t('title') }}
          </h4>
          <b-badge
            v-if="template.partial"
            data-test-id="badge-partial-template"
            variant="primary"
            class="ml-2"
          >
            {{ $t('partial') }}
          </b-badge>
        </template>

        <template #footer>
          <c-button-submit
            :disabled="!canCreate"
            :processing="processing"
            :success="success"
            :text="$t('admin:general.label.submit')"
            class="ml-auto"
            @submit="$emit('submit', template)"
          />
        </template>
      </b-card>
    </b-col>
  </b-row>
</template>

<script>
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import EditorToolbox from './EditorToolbox'
import EditorTextHtml from './EditorTextHtml'
import EditorTextPlain from './EditorTextPlain'
import EditorDocx from './EditorDocx'
import EditorUnsupported from './EditorUnsupported'
import { components } from '@cortezaproject/corteza-vue'

const { CAceEditor } = components

export default {

  components: {
    CAceEditor,
    EditorToolbox,
  },
  mixins: [
    listHelpers,
  ],

  i18nOptions: {
    namespaces: 'system.templates',
    keyPrefix: 'editor.content',
  },

  props: {
    template: {
      type: Object,
      required: true,
      default: () => ({}),
    },

    partials: {
      type: Array,
      required: false,
      default: () => [],
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    canCreate: {
      type: Boolean,
      required: true,
    },
  },

  data () {
    return {
      // @todo i18n this?
      previewData: '{\n  "variables": {\n    "param1": "value1",\n    "param2": {\n      "nestedParam1": "value2"\n    }\n  },\n  "options": {\n    "documentSize": "A4",\n    "contentScale": "1",\n    "orientation": "portrait",\n    "margin": "0.3"\n  }\n}\n',

      previewBlob: '',

      availableDrivers: [],
    }
  },

  computed: {
    editor () {
      switch (this.template.type) {
        case 'text/html':
          return EditorTextHtml
        case 'text/plain':
          return EditorTextPlain
        case 'application/vnd.openxmlformats-officedocument.wordprocessingml.document':
          return EditorDocx
        default:
          return EditorUnsupported
      }
    },

    canPreviewDocx () {
      return this.template.type === 'application/vnd.openxmlformats-officedocument.wordprocessingml.document' &&
        this.template.sourceFileID && this.template.sourceFileID !== '0'
    },

    isDocx () {
      return this.template.type === 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    },

    canPreviewHTML () {
      return !this.isDocx && this.availableDrivers.find(({ outputTypes }) => outputTypes.includes('text/html'))
    },

    canPreviewPDF () {
      return !this.isDocx && this.availableDrivers.find(({ outputTypes }) => outputTypes.includes('application/pdf'))
    },
  },

  async created () {
    this.availableDrivers = await this.$SystemAPI.templateRenderDrivers()
      .then(rsp => rsp.set)
      .catch(this.toastErrorHandler(this.$t('notification:template.driver.error')))
  },

  methods: {
    downloadPreview (ext) {
      this.incLoader()

      const cfg = {
        method: 'post',
        responseType: 'blob',
        url: this.$SystemAPI.templateRenderEndpoint({
          templateID: this.template.templateID,
          filename: 'preview',
          ext,
        }),
        data: JSON.parse(this.previewData),
      }

      this.$SystemAPI.api().request(cfg)
        .then(r => {
          const url = window.URL.createObjectURL(r.data)
          const a = document.createElement('a')
          a.href = url
          a.download = `preview.${ext}`
          document.body.appendChild(a)
          a.click()
          document.body.removeChild(a)
          window.URL.revokeObjectURL(url)
        })
        .catch(this.toastErrorHandler(this.$t('notification:template.preview.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    openPreview (ext) {
      this.incLoader()

      const cfg = {
        method: 'post',
        responseType: 'blob',
        url: this.$SystemAPI.templateRenderEndpoint({
          templateID: this.template.templateID,
          filename: 'preview',
          ext,
        }),
        data: JSON.parse(this.previewData),
      }

      this.$SystemAPI.api().request(cfg)
        .then(r => {
          this.previewBlob = window.URL.createObjectURL(r.data)
          window.open(this.previewBlob, '_newtab')
        })
        .catch(this.toastErrorHandler(this.$t('notification:template.preview.error')))
        .finally(() => {
          this.decLoader()
        })
    },
  },
}
</script>
