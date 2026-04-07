# Updated Plan: Common Table Filter (Sibling Records)

## Summary
The parent field dropdown should support TWO modes:
1. **Direct/Hierarchical** (existing): Parent → Grandparent → Great-grandparent 
2. **Common Field** (new): Fields in target module that link to the same module as current record (sibling relationship)

---

## Example

```
Current page: Candidate
- Has: personID → links to Person module

Target module: PersonContactMethods
- Has: personID → links to Person module ✓ (same as current record)
- Has: companyID → links to Company module ✗

When selecting parent field in PersonContactMethods:
- Show personID (common field - sibling relationship)
- Hide companyID (different module)
```

---

## Implementation

### Step 1: Extend parentFields Computed Property
In [`RecordListConfigurator.vue`](client/web/compose/src/components/PageBlocks/RecordListConfigurator.vue):

Keep existing hierarchical path finding AND add common field matching:

```javascript
parentFields () {
  if (!this.recordListModule || !this.record) {
    return []
  }

  const fields = []
  const currentRecordModuleID = this.record.moduleID

  // 1. Direct/Hierarchical paths (existing logic)
  // Get fields that point to current record's parent chain
  const hierarchicalFields = this.getHierarchicalParentFields()
  
  // 2. Common Field (sibling) - NEW
  // Get ALL Record fields from target module where options.moduleID matches current record's module
  const commonFields = this.recordListModule.fields
    .filter(field => field.kind === 'Record')
    .filter(field => 
      field.options && 
      field.options.moduleID === currentRecordModuleID
    )
    .map(field => ({
      ...field,
      isCommonField: true,  // Flag to identify this type
      label: `${field.name} (Common)`
    }))

  return [...hierarchicalFields, ...commonFields]
}
```

### Step 2: Update Filter Logic in RecordListBase.vue
In [`RecordListBase.vue`](client/web/compose/src/components/PageBlocks/RecordListBase.vue):
- Detect if parent field is a "common field" (sibling relationship)
- Apply appropriate filter: filter where field value equals current record's linked value

---

## Files to Modify

1. `client/web/compose/src/components/PageBlocks/RecordListConfigurator.vue` 
   - Extend `parentFields` computed property
   
2. `client/web/compose/src/components/PageBlocks/RecordListBase.vue`
   - Add logic to handle common field (sibling) filtering

---

## Algorithm

```javascript
// In RecordListConfigurator.vue
getCommonFields () {
  if (!this.recordListModule || !this.record) {
    return []
  }
  
  const currentRecordModuleID = this.record.moduleID
  
  return this.recordListModule.fields
    .filter(field => field.kind === 'Record')
    .filter(field => 
      field.options && 
      field.options.moduleID === currentRecordModuleID
    )
}
```
