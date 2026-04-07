# Plan: Parent Fields in Record List Field Picker

## Overview
Allow users to include fields from parent records in a record list, displayed alongside child module fields for unified ordering.

## User Flow
1. User adds a Record List block to a record page
2. Toggle switch "Include parent fields" 
3. Dropdown appears to select which parent field links to the record list module
4. Parent module fields appear in field picker alongside child fields
5. User can order both parent and child fields in the same list

## Implementation Steps

### Step 1: Add UI - Toggle and Parent Field Selector
**File:** `client/web/compose/src/components/PageBlocks/RecordListConfigurator.vue`

- Add `options.includeParentFields` (boolean) - toggle switch
- Add `options.parentField` (string) - dropdown to select parent field

**UI Location:** Below the existing field picker section (around line 80-92)

```vue
<!-- New section after field picker -->
<b-col
  v-if="onRecordPage"
  cols="12"
>
  <b-form-group>
    <c-input-checkbox
      v-model="options.includeParentFields"
      :label="$t('recordList.parentFields.enable')"
    />
  </b-form-group>
  
  <b-form-group
    v-if="options.includeParentFields"
    :label="$t('recordList.parentFields.field')"
  >
    <c-input-select
      v-model="options.parentField"
      :options="availableParentFields"
      label="label"
      :reduce="f => f.name"
    />
  </b-form-group>
</b-col>
```

### Step 2: Create Computed Property for Available Parent Fields
**File:** `client/web/compose/src/components/PageBlocks/RecordListConfigurator.vue`

**Logic:** Find fields in the parent module (this.module) that link to the record list module (this.recordListModule)

```javascript
availableParentFields () {
  if (!this.module || !this.recordListModule) {
    return []
  }
  
  const targetModuleID = this.recordListModule.moduleID
  
  // Find Record-kind fields in parent module that point to record list module
  return this.module.fields
    .filter(field => {
      if (field.kind !== 'Record') return false
      if (!field.options || !field.options.moduleID) return false
      return field.options.moduleID === targetModuleID
    })
    .map(field => ({
      ...field,
      label: field.name
    }))
}
```

### Step 3: Modify Field Picker to Include Parent Fields
**File:** `client/web/compose/src/components/PageBlocks/RecordListConfigurator.vue`

**Option A - Pass combined module:**
The field picker expects a single module. We need to either:
- Create a combined/pseudo module with fields from both modules
- Or extend the FieldPicker to accept additional fields

**Option B - Create enhanced field picker:**
Modify the field picker to accept an optional `extraFields` prop that combines with module fields

**Recommended Approach - Option B:**
1. Add `extraModule` and `extraModuleFields` props to FieldPicker
2. When both are provided, merge fields with a prefix/indicator

### Step 4: Update Field Picker Component
**File:** `client/web/compose/src/components/Common/FieldPicker.vue`

Add new props:
```javascript
extraModule: {
  type: Object,
  default: null
},

extraModuleFields: {
  type: Array,
  default: () => []
}
```

Modify `options` computed to merge extra fields:
```javascript
options () {
  let mFields = [...(this.fieldSubset ? this.module.filterFields(this.fieldSubset) : this.module.fields)]
  
  // Add extra fields from parent module with prefix
  if (this.extraModule && this.extraModuleFields.length > 0) {
    const extraFields = this.extraModuleFields.map(f => ({
      ...f,
      name: `parent.${f.name}`,  // Prefix to avoid conflicts
      label: `${this.extraModule.name}: ${f.label || f.name}`,
      isParentField: true
    }))
    mFields = [...mFields, ...extraFields]
  }
  
  // ... rest of existing logic
}
```

### Step 5: Pass Parent Module Info to Field Picker
**File:** `client/web/compose/src/components/PageBlocks/RecordListConfigurator.vue`

Modify the field picker usage:
```vue
<field-picker
  :module="recordListModule"
  :extra-module="options.includeParentFields ? parentModule : null"
  :extra-module-fields="options.includeParentFields ? parentModuleFields : []"
  :fields.sync="options.fields"
  class="mb-3"
  style="height: 50vh;"
/>
```

Add computed properties:
```javascript
parentModule () {
  if (this.options.includeParentFields && this.options.parentField) {
    return this.module  // Parent module
  }
  return null
},

parentModuleFields () {
  if (!this.parentModule || !this.options.parentField) {
    return []
  }
  
  // Get all fields from parent module
  return this.parentModule.fields.map(f => ({
    ...f,
    // Store original field name for reference
    originalName: f.name
  }))
}
```

### Step 6: Add Translations
**File:** `locale/en/corteza-webapp-compose/block.yaml`

```yaml
recordList:
  parentFields:
    enable: Include parent fields
    field: Parent link field
    description: Select which field links to the parent record
```

### Step 7: Handle Field Selection Storage
When fields are selected, they need to be stored with indicator if they're parent fields:

```javascript
// In the field picker, when storing selected fields
selectedFields.map(s => {
  const field = this.options.find(({ value }) => value === s).field
  if (field.isParentField) {
    return {
      ...field,
      name: field.originalName,  // Store original name without prefix
      isParentField: true,
      moduleID: this.extraModule.moduleID
    }
  }
  return field
})
```

## Technical Considerations

### 1. Field Name Collision Handling
Parent fields need unique identifiers in the picker but original names for storage. Use prefix strategy:
- Display: `parent.fieldName` 
- Storage: original name + `isParentField: true` flag

### 2. Backend API Requirements
The record list needs to fetch and display parent field values. This requires:
- API changes to support joining with parent record
- Passing `parentField` and `parentModuleID` to the API

### 3. RecordListBase.vue Changes
Need to update the record fetching logic to:
- Include parent module fields in the query
- Join with parent records based on the selected parent field

## Files to Modify

1. **`client/web/compose/src/components/PageBlocks/RecordListConfigurator.vue`**
   - Add toggle switch and dropdown UI
   - Add `availableParentFields` computed property
   - Add `parentModule` and `parentModuleFields` computed properties
   - Pass extra module info to field picker

2. **`client/web/compose/src/components/Common/FieldPicker.vue`**
   - Add `extraModule` and `extraModuleFields` props
   - Modify `options` computed to merge extra fields

3. **`locale/en/corteza-webapp-compose/block.yaml`**
   - Add translation keys for parent fields feature

4. **`client/web/compose/src/components/PageBlocks/RecordListBase.vue`**
   - Update record fetching to include parent fields
   - Handle displaying parent field values

## Example Use Case

**Scenario:** Candidate record page showing related Contacts

1. **Page:** Candidate record page (module: Candidate)
2. **Record List:** Contacts module (has `candidateID` field linking to Candidate)
3. **Configuration:**
   - Enable "Include parent fields"
   - Select parent field: `candidateID`
4. **Result:** Field picker shows both:
   - Child fields: name, email, phone (from Contacts)
   - Parent fields: firstName, lastName (from Candidate)
5. **User Experience:** Can order fields like:
   - firstName (parent)
   - name (child)
   - email (child)
   - lastName (parent)
