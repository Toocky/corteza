<template>
  <b-tab :title="$t('fileUpload.label')">
    <h5 class="mb-3">
      {{ $t('fileUpload.config.target') }}
    </h5>

    <!-- Module selector -->
    <b-form-group
      :label="$t('general.module')"
      label-class="text-primary"
    >
      <c-input-select
        v-model="options.moduleID"
        :options="modules"
        label="name"
        :reduce="o => o.moduleID"
        :placeholder="$t('fileUpload.config.modulePlaceholder')"
        default-value="0"
      />
    </b-form-group>

    <!-- File field selector (only shows File-kind fields from the selected module) -->
    <b-form-group
      v-if="targetModule"
      :label="$t('fileUpload.config.fileField')"
      label-class="text-primary"
    >
      <c-input-select
        v-model="options.fileFieldName"
        :options="fileFields"
        label="label"
        :reduce="o => o.name"
        :placeholder="$t('fileUpload.config.fileFieldPlaceholder')"
        :get-option-label="f => f.label || f.name"
      />
    </b-form-group>

    <hr v-if="targetModule">

    <!-- Button configuration -->
    <h5
      v-if="targetModule"
      class="mb-3"
    >
      {{ $t('fileUpload.config.buttonConfig') }}
    </h5>

    <b-form-group
      v-if="targetModule"
      :label="$t('fileUpload.config.buttonText')"
      label-class="text-primary"
    >
      <b-form-input
        v-model="options.buttonText"
        :placeholder="$t('fileUpload.config.buttonTextPlaceholder')"
      />
    </b-form-group>

    <b-form-group
      v-if="targetModule"
      :label="$t('fileUpload.config.buttonVariant')"
      label-class="text-primary"
    >
      <c-input-select
        v-model="options.buttonVariant"
        :options="buttonVariants"
        :reduce="o => o.value"
        label="text"
      />
    </b-form-group>

    <hr v-if="targetModule">

    <!-- Multi-file handling -->
    <h5
      v-if="targetModule"
      class="mb-3"
    >
      {{ $t('fileUpload.config.multiFileHandling') }}
    </h5>

    <b-form-group
      v-if="targetModule"
      label-class="text-primary"
    >
      <b-form-radio-group
        v-model="options.multiFileMode"
        stacked
      >
        <b-form-radio value="newRecord">
          {{ $t('fileUpload.config.newRecordPerFile') }}
        </b-form-radio>
        <b-form-radio
          value="multiValue"
          :disabled="!isFileFieldMulti"
        >
          {{ $t('fileUpload.config.multipleFilesPerRecord') }}
          <small
            v-if="!isFileFieldMulti"
            class="text-muted d-block"
          >
            {{ $t('fileUpload.config.multipleFilesDisabledHint') }}
          </small>
        </b-form-radio>
      </b-form-radio-group>
    </b-form-group>

    <hr v-if="targetModule">

    <!-- Extra fields -->
    <h5
      v-if="targetModule"
      class="mb-3"
    >
      {{ $t('fileUpload.config.extraFields') }}
    </h5>

    <b-form-group
      v-if="targetModule"
      :description="$t('fileUpload.config.extraFieldsDescription')"
    >
      <b-form-checkbox-group
        v-model="selectedExtraFieldNames"
        stacked
      >
        <b-form-checkbox
          v-for="field in availableExtraFields"
          :key="field.name"
          :value="field.name"
        >
          {{ field.label || field.name }}
          <small class="text-muted">({{ field.kind }})</small>
        </b-form-checkbox>
      </b-form-checkbox-group>
    </b-form-group>

    <hr v-if="targetModule">

    <!-- Behavior -->
    <h5
      v-if="targetModule"
      class="mb-3"
    >
      {{ $t('fileUpload.config.behavior') }}
    </h5>

    <b-form-group v-if="targetModule">
      <b-form-checkbox v-model="options.clearOnSubmit">
        {{ $t('fileUpload.config.clearOnSubmit') }}
      </b-form-checkbox>
    </b-form-group>

    <b-form-group
      v-if="targetModule"
      :label="$t('fileUpload.config.successMessage')"
      label-class="text-primary"
    >
      <b-form-input
        v-model="options.successMessage"
        :placeholder="$t('fileUpload.config.successMessagePlaceholder')"
      />
    </b-form-group>

    <b-form-group
      v-if="targetModule"
      :label="$t('fileUpload.config.acceptedMimeTypes')"
      label-class="text-primary"
    >
      <b-form-input
        v-model="options.acceptedMimeTypes"
        :placeholder="$t('fileUpload.config.acceptedMimeTypesPlaceholder')"
      />
    </b-form-group>
  </b-tab>
</template>

<script>
import base from './base'
import { NoID } from '@cortezaproject/corteza-js'
import { mapGetters } from 'vuex'

// Field kinds that make sense as extra fields on the upload block
const allowedExtraFieldKinds = [
  'String',
  'Number',
  'DateTime',
  'Select',
  'Record',
  'Bool',
  'Email',
  'Url',
]

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'FileUpload',

  extends: base,

  data () {
    return {
      // Tracks the moduleID at mount time so the watcher can
      // distinguish "module changed by user" from "initial load"
      initialModuleID: null,
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      modules: 'module/set',
    }),

    targetModule () {
      if (!this.options.moduleID || this.options.moduleID === NoID) return null
      return this.getModuleByID(this.options.moduleID)
    },

    fileFields () {
      if (!this.targetModule) return []
      return this.targetModule.fields.filter(f => f.kind === 'File')
    },

    selectedFileField () {
      if (!this.targetModule || !this.options.fileFieldName) return null
      return this.targetModule.fields.find(f => f.name === this.options.fileFieldName)
    },

    isFileFieldMulti () {
      return this.selectedFileField ? this.selectedFileField.isMulti : false
    },

    availableExtraFields () {
      if (!this.targetModule) return []
      return this.targetModule.fields.filter(f => {
        // Exclude the file field we're uploading to
        if (f.name === this.options.fileFieldName) return false
        // Exclude system fields
        if (f.isSystem) return false
        // Only allow sensible field kinds
        return allowedExtraFieldKinds.includes(f.kind)
      })
    },

    selectedExtraFieldNames: {
      get () {
        return (this.options.extraFields || []).map(f => f.name)
      },
      set (names) {
        this.options.extraFields = names.map(name => {
          const field = this.targetModule.fields.find(f => f.name === name)
          return {
            name,
            label: field ? (field.label || field.name) : name,
          }
        })
      },
    },

    buttonVariants () {
      return [
        { value: 'primary', text: 'Primary' },
        { value: 'secondary', text: 'Secondary' },
        { value: 'success', text: 'Success' },
        { value: 'danger', text: 'Danger' },
        { value: 'warning', text: 'Warning' },
        { value: 'info', text: 'Info' },
        { value: 'light', text: 'Light' },
        { value: 'dark', text: 'Dark' },
        { value: 'outline-primary', text: 'Outline Primary' },
        { value: 'outline-secondary', text: 'Outline Secondary' },
      ]
    },
  },
  watch: {
    'options.moduleID' (newID) {
      // Skip the reset on initial load — only reset when the user
      // actually changes the module in the dropdown
      if (this.initialModuleID === newID) {
        this.initialModuleID = null
        return
      }
      this.initialModuleID = null

      this.options.fileFieldName = ''
      this.options.extraFields = []
      this.options.multiFileMode = 'newRecord'
    },

    isFileFieldMulti (multi) {
      // If file field is no longer multi, force back to newRecord
      if (!multi && this.options.multiFileMode === 'multiValue') {
        this.options.multiFileMode = 'newRecord'
      }
    },
  },

  created () {
    this.initialModuleID = this.options.moduleID || null
  },

}
</script>
