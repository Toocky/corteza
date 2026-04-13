<template>
  <b-card
    data-test-id="card-template-toolbox"
    no-body
    header-class="border-bottom"
    class="shadow-sm h-100"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <b-card-body>
      <span
        v-for="sec in sections"
        :key="sec.key"
      >
        <b-btn
          :data-test-id="toolboxSectionLabelCypressId(sec.key)"
          variant="light"
          class="mb-2"
          block
          @click="openSection(sec.key)"
        >
          {{ $t(sec.key) }}
        </b-btn>
        <b-collapse
          :visible="expandedSections[sec.key]"
          role="tabpanel"
          class="pb-2 px-0"
        >
          <b-list-group flush>
            <b-list-group-item
              v-for="(opt, i) in sec.options"
              :key="opt.label + i"
              :data-test-id="toolboxOptionLabelCypressId(opt.label)"
              class="px-0 text-wrap"
              @click="opt.onClick || (() => {})"
            >
              {{ opt.label }}
              <b-btn
                data-test-id="button-copy"
                variant="link"
                class="pr-0 float-right"
                @click="copyToCb(opt.copyValue())"
              >
                <font-awesome-icon
                  v-if="opt.copyValue"
                  :icon="['far', 'copy']"
                />
              </b-btn>
            </b-list-group-item>
          </b-list-group>
        </b-collapse>

      </span>
    </b-card-body>
  </b-card>
</template>

<script>
import copy from 'copy-to-clipboard'

export default {
  i18nOptions: {
    namespaces: 'system.templates',
    keyPrefix: 'editor.content.toolbox',
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
  },

  data () {
    return {
      expandedSections: {},
    }
  },

  computed: {
    isDocx () {
      return this.template.type === 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    },

    cvSampleData () {
      return {
        variables: {
          contact: {
            details: {
              name: 'Eric Banner',
              gender: 'Male',
              nationality: 'Australian',
              marital_status: 'Married',
              current_location: 'Sydney, NSW',
              date_of_birth: 'April 15, 1968',
            },
            education: [
              {
                name: 'Bachelor of Science',
                institution: 'University of Sydney',
                started: 'February 1986',
                ended: 'November 1989',
                gpa: '3.8',
                major: 'Computer Science',
              },
              {
                name: 'Master of Business Administration',
                institution: 'UNSW Sydney',
                started: 'March 1992',
                ended: 'December 1993',
                gpa: '3.5',
                major: 'Technology Management',
              },
            ],
            languages: [
              { language: 'English', proficiency: 'Native' },
              { language: 'Japanese', proficiency: 'Conversational' },
            ],
            current_company: 'Meridian Technologies',
            current_position: 'Director of Engineering',
            availability: '4 weeks notice',
            current_salary: '185,000 AUD',
            expected_salary: '210,000 AUD',
            work_experience: [
              {
                company: 'DataFlow Systems',
                start_date: 'January 1990',
                end_date: 'August 1995',
                title: 'Software Developer',
                experience_block: 'Developed backend services for financial data processing. Built real-time market data feeds and automated reporting systems. Mentored junior developers and established code review practices.',
              },
              {
                company: 'Pacific Digital',
                start_date: 'September 1995',
                end_date: 'March 2005',
                title: 'Senior Software Architect',
                experience_block: 'Led architecture of the company\'s cloud migration initiative. Designed microservices platform serving 2M+ daily users. Managed a team of 12 engineers across Sydney and Melbourne offices.',
              },
              {
                company: 'Meridian Technologies',
                start_date: 'April 2005',
                end_date: 'Present',
                title: 'Director of Engineering',
                experience_block: 'Oversees engineering department of 45 staff across three product lines. Introduced CI/CD pipelines reducing deployment time by 80%. Drove adoption of agile methodologies and cross-functional team structure.',
              },
            ],
            certification: [
              {
                name: 'AWS Solutions Architect Professional',
                issuer: 'Amazon Web Services',
                date: 'March 2019',
              },
              {
                name: 'Certified ScrumMaster (CSM)',
                issuer: 'Scrum Alliance',
                date: 'June 2012',
              },
            ],
          },
        },
        options: {
          documentSize: 'A4',
          contentScale: '1',
          orientation: 'portrait',
          margin: '0.3',
        },
      }
    },

    sections () {
      const partials = this.partials.map(p => ({
        label: p.meta.short || p.handle,
        copyValue: () => `{{template "${p.handle}" }}`,
      }))

      const rr = []
      if (partials.length) {
        rr.push({
          key: 'partials',
          options: partials,
        })
      }

      if (this.isDocx) {
        rr.push({
          key: 'snippets.label',
          options: [
            {
              label: this.$t('snippets.docxVariable'),
              copyValue: () => '{{variableName}}',
            },
            {
              label: this.$t('snippets.docxLoop'),
              copyValue: () => '{{#each listName}}\n  {{fieldName}}\n{{/each}}',
            },
            {
              label: this.$t('snippets.docxConditional'),
              copyValue: () => '{{#if hasListName}}\n  Content here\n{{/if}}',
            },
            {
              label: this.$t('snippets.docxCurrentItem'),
              copyValue: () => '{{this}}',
            },
          ],
        },
        {
          key: 'samples.label',
          options: [
            {
              label: this.$t('samples.cvTestData'),
              copyValue: () => JSON.stringify(this.cvSampleData, null, 2),
            },
          ],
        })
      } else {
        rr.push({
          key: 'snippets.label',
          options: [
            {
              label: this.$t('snippets.interpolate'),
              copyValue: () => '{{.parameter}}',
            },
            {
              label: this.$t('snippets.iterator'),
              copyValue: () => '{{range $index, $element := .ListOfItems}}\n\n{{end}}',
            },
            {
              label: this.$t('snippets.funcCall'),
              copyValue: () => '{{funcName param1 param2 paramN}}',
            },
          ],
        },
        {
          key: 'samples.label',
          options: [
            {
              label: this.$t('samples.defaultHTML'),
              copyValue: () => `<!DOCTYPE html>
<html>
<head>
  <meta charset='utf-8'>
  <meta http-equiv='X-UA-Compatible' content='IE=edge'>
  <title>Title</title>
  <meta name='viewport' content='width=device-width, initial-scale=1'>
</head>
<body>
  <h1>Hello, world!</h1>
</body>
</html>`,
            },
          ],
        })
      }

      return rr
    },
  },

  methods: {
    openSection (sec) {
      this.$set(this.expandedSections, sec, !this.expandedSections[sec])
    },

    copyToCb: copy,

    toolboxSectionLabelCypressId (section) {
      return `button-${this.$t(section).toLowerCase()}`
    },

    toolboxOptionLabelCypressId (label) {
      return label.toLowerCase().split(' ').join('-')
    },
  },

}
</script>
