<template>
  <div class="p-3">
    <b-form-group
      :label="$t('docx.sourceFile')"
      label-class="text-primary"
    >
      <div
        v-if="template.sourceFileID && template.sourceFileID !== '0'"
        class="mb-2 d-flex align-items-center"
      >
        <b-badge variant="success" class="mr-2">
          <font-awesome-icon :icon="['fas', 'file-word']" class="mr-1" />
          {{ $t('docx.fileAttached') }}
        </b-badge>

        <b-btn
          variant="outline-secondary"
          size="sm"
          class="mr-2"
          @click="downloadTemplate"
        >
          {{ $t('docx.downloadTemplate') }}
        </b-btn>

        <b-btn
          variant="outline-danger"
          size="sm"
          @click="removeFile"
        >
          {{ $t('docx.removeFile') }}
        </b-btn>
      </div>

      <b-form-file
        v-model="file"
        :placeholder="$t('docx.uploadPlaceholder')"
        :drop-placeholder="$t('docx.dropPlaceholder')"
        accept=".docx,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        @input="onFileSelected"
      />
    </b-form-group>

    <b-card
      no-body
      border-variant="info"
      class="mb-0"
    >
      <b-card-header
        header-bg-variant="info"
        header-text-variant="white"
      >
        <strong>{{ $t('docx.guideTitle') }}</strong>
      </b-card-header>
      <b-card-body>
        <p class="mb-2">
          {{ $t('docx.guideIntro') }}
        </p>

        <h6 class="mt-3 mb-1">
          {{ $t('docx.guideVariablesTitle') }}
        </h6>
        <p class="mb-1 text-muted small">
          {{ $t('docx.guideVariablesDesc') }}
        </p>
        <ul class="pl-3 mb-2">
          <li><code v-pre>{{variableName}}</code> — {{ $t('docx.guideVariableSimple') }}</li>
          <li>
            <code v-pre>{{contactDetailsName}}</code> — {{ $t('docx.guideVariableNested') }}
          </li>
        </ul>

        <h6 class="mt-3 mb-1">
          {{ $t('docx.guideLoopsTitle') }}
        </h6>
        <p class="mb-1 text-muted small">
          {{ $t('docx.guideLoopsDesc') }}
        </p>
        <pre
          v-pre
          class="bg-light p-2 rounded small mb-2"
        >{{#each contactEducation}}
  {{name}} - {{institution}}
  {{started}} to {{ended}}
{{/each}}</pre>

        <h6 class="mt-3 mb-1">
          {{ $t('docx.guideConditionalsTitle') }}
        </h6>
        <p class="mb-1 text-muted small">
          {{ $t('docx.guideConditionalsDesc') }}
        </p>
        <pre
          v-pre
          class="bg-light p-2 rounded small mb-2"
        >{{#if hasContactCertification}}
  Certifications:
  {{#each contactCertification}}
    {{name}} - {{issuer}}
  {{/each}}
{{/if}}</pre>

        <h6 class="mt-3 mb-1">
          {{ $t('docx.guideCurrentItemTitle') }}
        </h6>
        <p class="mb-1 text-muted small">
          {{ $t('docx.guideCurrentItemDesc') }}
        </p>
        <pre
          v-pre
          class="bg-light p-2 rounded small mb-2"
        >{{#each contactSkills}}
  {{this}}
{{/each}}</pre>

        <h6 class="mt-3 mb-1">
          {{ $t('docx.guideDataTitle') }}
        </h6>
        <p class="mb-1 text-muted small">
          {{ $t('docx.guideDataDesc') }}
        </p>
        <div class="small">
          <table class="table table-sm table-bordered mb-2">
            <thead class="thead-light">
              <tr>
                <th>{{ $t('docx.guideTableJson') }}</th>
                <th>{{ $t('docx.guideTablePlaceholder') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td><code>contact.details.name</code></td>
                <td v-pre>
                  <code>{{contactDetailsName}}</code>
                </td>
              </tr>
              <tr>
                <td><code>contact.current_company</code></td>
                <td v-pre>
                  <code>{{contactCurrent_company}}</code>
                </td>
              </tr>
              <tr>
                <td><code>contact.education[]</code></td>
                <td v-pre>
                  <code>{{#each contactEducation}}</code>
                </td>
              </tr>
              <tr>
                <td><code>contact.education[].name</code></td>
                <td v-pre>
                  <code>{{name}}</code> <span class="text-muted">(inside loop)</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <h6 class="mt-3 mb-1">
          {{ $t('docx.guideTipsTitle') }}
        </h6>
        <ul class="pl-3 mb-0 small">
          <li>{{ $t('docx.guideTip1') }}</li>
          <li>{{ $t('docx.guideTip2') }}</li>
          <li>{{ $t('docx.guideTip3') }}</li>
        </ul>
      </b-card-body>
    </b-card>
  </div>
</template>

<script>
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'

export default {
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
  },

  data () {
    return {
      file: null,
    }
  },

  methods: {
    async onFileSelected (file) {
      if (!file) return

      const formData = new FormData()
      formData.append('upload', file)
      if (this.template.templateID) {
        formData.append('templateID', this.template.templateID)
      }

      try {
        const rsp = await this.$SystemAPI.api().request({
          method: 'post',
          url: '/attachment/template/',
          data: formData,
          headers: { 'Content-Type': 'multipart/form-data' },
        })

        if (rsp.data && rsp.data.response && rsp.data.response.attachmentID) {
          this.template.sourceFileID = rsp.data.response.attachmentID
        }
      } catch (e) {
        this.toastErrorHandler(this.$t('docx.uploadError'))(e)
      }
    },

    async downloadTemplate () {
      try {
        const rsp = await this.$SystemAPI.api().request({
          method: 'get',
          responseType: 'blob',
          url: `/attachment/template/${this.template.sourceFileID}/original/template.docx`,
        })
        const url = window.URL.createObjectURL(rsp.data)
        const a = document.createElement('a')
        a.href = url
        a.download = `${this.template.handle || 'template'}.docx`
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
        window.URL.revokeObjectURL(url)
      } catch (e) {
        this.toastErrorHandler(this.$t('docx.downloadError'))(e)
      }
    },

    async removeFile () {
      const attachmentID = this.template.sourceFileID
      this.template.sourceFileID = '0'
      this.file = null

      if (attachmentID && attachmentID !== '0') {
        try {
          await this.$SystemAPI.attachmentDelete({
            kind: 'template',
            attachmentID,
          })
        } catch (e) {
          this.toastErrorHandler(this.$t('docx.removeError'))(e)
        }
      }
    },
  },
}
</script>
