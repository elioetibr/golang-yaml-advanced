# TransformDSL Flexible Input Feature

## Overview

The `TransformDSL` now supports flexible input types through the new `NewTransformDSLFrom()` function, similar to how `MergeFlexible()` works. This allows you to create transformation pipelines from various input formats without manual conversion.

## What Was Added

### 1. New Function: `NewTransformDSLFrom(input interface{}) (*TransformDSL, error)`

Creates a TransformDSL with flexible input that accepts:
- `*NodeTree` - Direct node tree
- `*Node` - Single node (wrapped in NodeTree)
- `*Document` - YAML document (wrapped in NodeTree)
- `[]byte` - YAML bytes (parsed)
- `string` - YAML string (parsed)
- `interface{}` - Go structs, maps, slices (marshaled to YAML then parsed)

### 2. Updated `TransformDSL` Structure

```go
type TransformDSL struct {
    transforms []Transform
    errors     []error
    inputTree  *NodeTree // Optional: stores input when created with NewTransformDSLFrom
}
```

### 3. Updated `Apply()` Method

```go
func (dsl *TransformDSL) Apply(tree *NodeTree) (*NodeTree, error)
```

Now supports two patterns:
- **Old pattern (backward compatible)**: `Apply(tree)` - explicitly pass a tree
- **New pattern**: `Apply(nil)` - uses the stored inputTree from `NewTransformDSLFrom()`

### 4. Helper Function: `convertInputToNodeTree(input interface{}) (*NodeTree, error)`

Internal function that handles conversion from various input types to NodeTree.

## Usage Examples

### Pattern 1: Old Pattern (Still Supported)

```go
// Create DSL without input
dsl := yaml.NewTransformDSL()

// Load tree separately
tree, _ := yaml.UnmarshalYAML([]byte("name: John"))

// Apply with explicit tree
result, _ := dsl.RemoveKey("email").Apply(tree)
```

### Pattern 2: New Pattern - From YAML String

```go
yamlStr := `
name: Alice
role: admin
email: alice@example.com
`

// Create DSL with input
dsl, _ := yaml.NewTransformDSLFrom(yamlStr)

// Apply without passing tree (uses stored input)
result, _ := dsl.RemoveKey("email").Apply(nil)
```

### Pattern 3: From YAML Bytes

```go
yamlBytes := []byte("name: Bob\nage: 30")

dsl, _ := yaml.NewTransformDSLFrom(yamlBytes)
result, _ := dsl.RenameKey("age", "years").Apply(nil)
```

### Pattern 4: From Go Map

```go
config := map[string]interface{}{
    "host": "localhost",
    "port": 8080,
}

dsl, _ := yaml.NewTransformDSLFrom(config)
result, _ := dsl.SortKeys().Apply(nil)
```

### Pattern 5: From Node

```go
node := yaml.NewMappingNode()
node.AddKeyValue(yaml.NewScalarNode("key"), yaml.NewScalarNode("value"))

dsl, _ := yaml.NewTransformDSLFrom(node)
result, _ := dsl.AddComment("My comment").Apply(nil)
```

### Pattern 6: From NodeTree

```go
tree, _ := yaml.UnmarshalYAML([]byte("app: myapp"))

dsl, _ := yaml.NewTransformDSLFrom(tree)
result, _ := dsl.RenameKey("app", "application").Apply(nil)
```

### Pattern 7: Complex Chaining

```go
yamlStr := `
application:
  name: myapp
  temp_field: remove_me
  password: secret123
`

dsl, _ := yaml.NewTransformDSLFrom(yamlStr)

result, _ := dsl.
    RemoveKey("temp_field").
    RemoveKey("password").
    SortKeys().
    Apply(nil)
```

## Benefits

1. **Convenience**: No need to manually convert inputs to NodeTree
2. **Flexibility**: Accepts any common YAML-related type
3. **Backward Compatible**: Old pattern with `NewTransformDSL()` + `Apply(tree)` still works
4. **Type Safety**: Returns error if conversion fails
5. **Consistent API**: Follows same pattern as `MergeFlexible()`

## Error Handling

```go
// Returns error if input is nil
dsl, err := yaml.NewTransformDSLFrom(nil)
if err != nil {
    // Handle error: "input cannot be nil"
}

// Returns error if Apply called without tree and no inputTree
dsl := yaml.NewTransformDSL()
result, err := dsl.Apply(nil)
if err != nil {
    // Handle error: "tree is nil and no input tree was provided during initialization"
}

// Returns error if YAML is invalid
dsl, err := yaml.NewTransformDSLFrom([]byte("{invalid yaml"))
if err != nil {
    // Handle YAML parsing error
}
```

## Testing

Comprehensive tests added in `advanced_test.go`:
- `TestNewTransformDSLFrom` - Tests all input types
- `TestTransformDSLApplyWithStoredTree` - Tests both old and new patterns
- `TestTransformDSLFromWithTransformations` - Tests transformations with stored input
- `TestConvertInputToNodeTree` - Tests conversion function
- `TestTransformDSLFromChaining` - Tests complex transformation chains

All tests pass with 100% coverage of new functionality.

## Design Decisions

1. **Separate Function**: Created `NewTransformDSLFrom()` instead of modifying `NewTransformDSL()` to maintain backward compatibility
2. **Optional Field**: Added `inputTree` field to TransformDSL as optional (nil for old pattern)
3. **Apply Flexibility**: Modified `Apply()` to accept nil and use stored tree, maintaining backward compatibility
4. **Error First**: Returns error immediately if conversion fails, following Go conventions
5. **Node Wrapping**: Single nodes are wrapped in Document nodes for consistency

## Future Enhancements

Potential improvements:
1. Add `MustNewTransformDSLFrom()` variant that panics on error
2. Add method to change the stored inputTree after creation
3. Add method to get the stored inputTree
4. Support for multiple document inputs

## Implementation Files

- `yaml/advanced.go` - Core implementation (lines 568-652)
- `yaml/advanced_test.go` - Comprehensive tests (lines 1502-1879)
- `examples/transform_dsl_flexible/main.go` - Usage examples

## Related Features

- `MergeFlexible()` in `core.go` - Similar flexible input pattern for merging
- `ConvertToNodeTree()` in `core.go` - Similar conversion logic for merging operations
- `interfaces.go` - TreeTransformer interface for dependency injection
