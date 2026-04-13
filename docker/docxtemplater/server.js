import express from 'express'
import Docxtemplater from 'docxtemplater'
import PizZip from 'pizzip'
import multer from 'multer'
import expressions from 'angular-expressions'

// Custom filters matching the user's existing setup
expressions.filters.upper = (input) => input ? input.toUpperCase() : ''
expressions.filters.lower = (input) => input ? input.toLowerCase() : ''
expressions.filters.capitalize = (input) => {
  if (!input) return ''
  return input.charAt(0).toUpperCase() + input.slice(1).toLowerCase()
}
expressions.filters.title = (input) => {
  if (!input) return ''
  return input.replace(/\w\S*/g, (word) => {
    if (!word || word.length === 0) return word
    return (word[0]?.toUpperCase() ?? '') + word.slice(1).toLowerCase()
  })
}
expressions.filters.trim = (input) => input ? input.trim() : ''

function angularParser (tag) {
  return {
    get: function (scope) {
      return expressions.compile(tag)(scope)
    }
  }
}

const app = express()
const upload = multer({ storage: multer.memoryStorage(), limits: { fileSize: 32 * 1024 * 1024 } })

// Health check
app.get('/health', (req, res) => res.json({ status: 'ok' }))

// Render: accepts multipart form with "template" (docx file) and "data" (JSON string)
app.post('/render', upload.single('template'), (req, res) => {
  try {
    if (!req.file) {
      return res.status(400).json({ error: 'template file is required' })
    }

    const dataStr = req.body.data
    if (!dataStr) {
      return res.status(400).json({ error: 'data field is required' })
    }

    const data = JSON.parse(dataStr)

    const zip = new PizZip(req.file.buffer)
    const doc = new Docxtemplater(zip, {
      paragraphLoop: true,
      linebreaks: true,
      parser: angularParser,
      nullGetter: () => '',
    })

    doc.render(data)

    const out = doc.getZip().generate({ type: 'nodebuffer' })
    res.set('Content-Type', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document')
    res.send(out)
  } catch (err) {
    console.error('Render error:', err.message)
    if (err.properties && err.properties.errors) {
      console.error(JSON.stringify(err.properties.errors, null, 2))
      return res.status(422).json({
        error: err.message,
        details: err.properties.errors.map(e => ({
          id: e.id,
          message: e.message,
          properties: e.properties,
        })),
      })
    }
    res.status(500).json({ error: err.message })
  }
})

const port = process.env.PORT || 3000
app.listen(port, '0.0.0.0', () => {
  console.log(`docxtemplater renderer listening on port ${port}`)
})
