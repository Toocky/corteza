# Implementation Plan: Parent Fields in Record List

## Current State Analysis

### Already Implemented (RecordListConfigurator.vue)
- `parentModuleConfig` computed property (lines 1409-1447)
- `availableParentFields` computed property (lines 1385-1403)
- UI toggle for `includeParentFields` (lines 95-106)
- Dropdown for `parentField` selection (lines 108-125)
- Passing `extra-module` and `extra-module-fields` to FieldPicker (lines 87-88)

### Needs Implementation

#### File 1: FieldPicker.vue
**Status:** Missing `extraModule` and `extraModuleFields` props and merge logic

#### File 2: RecordListBase.vue
**Status:** Missing runtime data, computed properties, methods, and template updates

---

## Implementation Steps

### Step 1: Update FieldPicker.vue - Add Props
**File:** `client/web/compose/src/components/Common/FieldPicker.vue`
**Location:** After line 97 (after `fieldSubset` prop)

Add two new props:
- `extraModule` (Object, default: null) - parent module object with moduleID and name
- `extraModuleFields` (Array, default: []) - array of field objects from parent module

### Step 2: Update FieldPicker.vue - Modify `options` Computed
**File:** `client/web/compose/src/components/Common/FieldPicker.vue`
**Location:** Replace `options` computed (lines 127-181)

Update the `options` computed to:
1. Build `extraFieldOptions` array when `extraModule` and `extraModuleFields` are provided
2. Use `moduleID::fieldName` as unique key to avoid name collisions
3. Add `isParentField`, `parentModuleID`, `parentModuleName`, `originalName` flags to each parent field
4. Concatenate `extraFieldOptions` after the base options

### Step 3: Update FieldPicker.vue - Template Changes
**File:** `client/web/compose/src/components/Common/FieldPicker.vue`
**Location:** Template slot (lines 17-43)

Update the template to:
1. Show `field.originalName || field.name` instead of just `field.name`
2. Show `field.label` with parent module context when available
3. Add a small badge showing `field.parentModuleName` when `field.isParentField` is true

### Step 4: Update RecordListBase.vue - Add Data Property
**File:** `client/web/compose/src/components/PageBlocks/RecordListBase.vue`
**Location:** In `data()` function (around line 1221)

Add `resolvedParentRecords: {}` to store fetched parent record objects (map of parentRecordID -> Record)

### Step 5: Update RecordListBase.vue - Add Computed Properties
**File:** `client/web/compose/src/components/PageBlocks/RecordListBase.vue`
**Location:** In `computed` section (after line 1299)

Add two computed properties:
1. `parentFieldConfigs` - filters `options.fields` for fields with `isParentField: true`
2. `parentLinkFieldName` - returns `options.parentField` (the field in child module linking to parent)

### Step 6: Update RecordListBase.vue - Modify `fields` Computed
**File:** `client/web/compose/src/components/PageBlocks/RecordListBase.vue`
**Location:** Replace `fields` computed (lines 1601-1638)

Update to:
1. Separate parent fields from child fields before filtering
2. Build parent field column definitions manually (since they don't exist in child module)
3. Merge columns in the order specified in `options.fields`
4. Add `isParentField`, `parentModuleID`, `originalName` to parent field columns

### Step 7: Update RecordListBase.vue - Add `fetchParentRecords` Method
**File:** `client/web/compose/src/components/PageBlocks/RecordListBase.vue`
**Location:** In `methods` section (after line 1939)

Add method that:
1. Checks if parent fields are configured (`options.includeParentFields`, `parentLinkFieldName`, `parentFieldConfigs`)
2. Finds the link field definition in child module
3. Collects unique parent record IDs from child records
4. Builds a query to fetch all parent records in one API call
5. Creates Record objects and stores in `resolvedParentRecords` map

### Step 8: Update RecordListBase.vue - Call `fetchParentRecords` in `pullRecords`
**File:** `client/web/compose/src/components/PageBlocks/RecordListBase.vue`
**Location:** In `pullRecords` method, Promise.all block (around line 2684)

Add `this.fetchParentRecords(namespaceID, records)` to the Promise.all array

### Step 9: Update RecordListBase.vue - Add `getParentFieldValue` Helper
**File:** `client/web/compose/src/components/PageBlocks/RecordListBase.vue`
**Location:** In `methods` section

Add helper method that:
1. Gets parent record ID from child record's link field value
2. Looks up parent record in `resolvedParentRecords`
3. Returns the value from the parent record for the specified field

### Step 10: Update RecordListBase.vue - Template Changes
**File:** `client/web/compose/src/components/PageBlocks/RecordListBase.vue`
**Location:** In template, `<b-td v-for="field in fields">` section (around line 465)

Update to:
1. Check if `field.isParentField` is true
2. If parent field: render using `field-viewer` with parent record and parent module
3. If not parent field: use existing rendering logic unchanged
4. Show "—" when parent record is not loaded or link field is empty

### Step 11: Update RecordListBase.vue - Reset in `setDefaultValues`
**File:** `client/web/compose/src/components/PageBlocks/RecordListBase.vue`
**Location:** In `setDefaultValues` method (around line 3099)

Add `this.resolvedParentRecords = {}` to reset the parent records cache

---

## Data Flow Diagram

```
Configurator (edit time):
  user picks parentField (link field in child module)
    -> parentModuleConfig computes { extraModule, extraModuleFields }
    -> FieldPicker shows parent module fields in the picker
    -> user selects parent fields -> stored in options.fields with isParentField: true

RecordListBase (runtime):
  pullRecords() fetches child records
    -> fetchParentRecords() collects unique parent IDs from child records
    -> single API call fetches all parent records
    -> stored in resolvedParentRecords map (recordID -> Record)
  
  Template renders rows:
    -> for each field column:
      -> if field.isParentField: look up resolvedParentRecords[item.r.values[parentLinkFieldName]]
        -> pass that parent Record to field-viewer with the parent module
      -> else: existing rendering logic unchanged
```

---

## Key Considerations

1. **Vue 2 Reactivity:** Always replace the whole object (`this.resolvedParentRecords = { ...newData }`) rather than mutating individual keys
2. **Field Name Collisions:** Use `moduleID::fieldName` as unique key in FieldPicker to avoid conflicts
3. **Serialization:** Ensure `isParentField`, `parentModuleID`, `originalName` flags survive save/load round-trip
4. **Error Handling:** `fetchParentRecords` should fail gracefully - rows show "—" if parent record can't be loaded
5. **Performance:** Batch all parent record IDs into a single API call to minimize requests

---

## Files to Modify

1. `client/web/compose/src/components/Common/FieldPicker.vue` - Add props and merge logic
2. `client/web/compose/src/components/PageBlocks/RecordListBase.vue` - Add runtime logic and template updates

## Files Already Complete

1. `client/web/compose/src/components/PageBlocks/RecordListConfigurator.vue` - UI and computed properties ready
