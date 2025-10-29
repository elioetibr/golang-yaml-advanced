// Package yaml provides advanced YAML processing capabilities for Go that
// extend gopkg.in/yaml.v3 with comment preservation, complex structure handling,
// and advanced manipulation capabilities.
//
// This library is designed for enterprise-grade YAML processing where comment
// and formatting preservation is critical, such as Kubernetes manifests,
// Helm charts, and configuration management.
//
// # Features
//
//   - Full comment preservation (head, line, and foot comments)
//   - Integer preservation (prevents scientific notation for large numbers)
//   - Multi-document support (handles multiple YAML documents in a single file)
//   - Anchors and aliases support
//   - Schema validation with JSON Schema-like definitions
//   - Streaming parser for large files
//   - Transformation DSL for fluent YAML manipulations
//   - Flexible merging with multiple input types
//   - Diff comparison between YAML trees
//   - Empty line handling with configurable policies
//
// # Basic Usage
//
// Parse YAML with comment preservation:
//
//	tree, err := yaml.UnmarshalYAML(yamlData)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Access the first document
//	doc := tree.Documents[0]
//
//	// Serialize back to YAML (preserves comments)
//	output, err := tree.ToYAML()
//
// # Merging YAML Trees
//
// Merge two YAML trees while preserving comments from both:
//
//	base, _ := yaml.UnmarshalYAML(baseYAML)
//	override, _ := yaml.UnmarshalYAML(overrideYAML)
//	merged := yaml.MergeTrees(base, override)
//
// # Flexible Input Types
//
// Use MergeFlexible to work with various input types:
//
//	// Accepts: *NodeTree, []byte, string, maps, structs
//	result, err := yaml.MergeFlexible(base, override)
//
// # Transformation DSL
//
// Transform YAML using a fluent interface:
//
//	// Old pattern
//	tree, _ := yaml.UnmarshalYAML(data)
//	dsl := yaml.NewTransformDSL()
//	result, _ := dsl.
//	    RemoveKey("temp_field").
//	    RenameKey("old", "new").
//	    SortKeys().
//	    Apply(tree)
//
//	// New pattern with flexible input
//	dsl, _ := yaml.NewTransformDSLFrom(data)
//	result, _ := dsl.
//	    RemoveKey("temp_field").
//	    RenameKey("old", "new").
//	    Apply(nil) // Uses stored input
//
// # Schema Validation
//
// Validate YAML against a schema:
//
//	schema := &yaml.Schema{
//	    Type: "object",
//	    Properties: map[string]*yaml.Schema{
//	        "name": {Type: "string", MinLength: intPtr(1)},
//	        "age":  {Type: "integer", Minimum: float64Ptr(0)},
//	    },
//	    Required: []string{"name"},
//	}
//
//	errors := schema.Validate(node, "$")
//	if len(errors) > 0 {
//	    // Handle validation errors
//	}
//
// # Node Traversal
//
// Walk the YAML tree:
//
//	root.Walk(func(n *yaml.Node) bool {
//	    if n.Kind == yaml.ScalarNode {
//	        fmt.Printf("Found scalar: %v\n", n.Value)
//	    }
//	    return true // Continue traversal
//	})
//
// Find specific nodes:
//
//	// Find first matching node
//	node := root.Find(func(n *yaml.Node) bool {
//	    return n.Kind == yaml.ScalarNode && n.Value == "target"
//	})
//
//	// Find all matching nodes
//	nodes := root.FindAll(func(n *yaml.Node) bool {
//	    return n.Tag == "!!str"
//	})
//
// # Streaming Parser
//
// For large files (>100MB), use the streaming parser:
//
//	parser := yaml.NewStreamParser(reader)
//	parser.SetDocumentCallback(func(tree *yaml.NodeTree) error {
//	    // Process each document as it's parsed
//	    return processDocument(tree)
//	})
//	err := parser.Parse()
//
// # Diff Comparison
//
// Compare two YAML trees and get detailed differences:
//
//	diffs := yaml.DiffTrees(oldTree, newTree)
//	for _, diff := range diffs {
//	    fmt.Printf("%s at %s: %v -> %v\n",
//	        diff.Type, diff.Path, diff.OldValue, diff.NewValue)
//	}
//
// # Empty Line Handling
//
// Configure how empty lines are handled:
//
//	tree.EmptyLineConfig = yaml.NoEmptyLinesConfig()        // Remove all
//	tree.EmptyLineConfig = yaml.NormalizedEmptyLineConfig(1) // Normalize to 1
//	tree.EmptyLineConfig = yaml.DefaultEmptyLineConfig()     // Preserve original
//
// # Dependency Injection
//
// Use interfaces for better testability:
//
//	func ProcessYAML(validator yaml.SchemaValidator, node *yaml.Node) error {
//	    errors := validator.Validate(node, "$")
//	    // ... handle validation
//	}
//
// Available interfaces: SchemaValidator, TreeTransformer, DocumentParser
//
// # Common Use Cases
//
// Configuration Management:
//
//	baseConfig := loadConfig("config.yaml")
//	envConfig := loadConfig(fmt.Sprintf("config.%s.yaml", env))
//	finalConfig := yaml.MergeTrees(baseConfig, envConfig)
//
// Helm Values Processing:
//
//	values := loadValues("values.yaml")
//	overrides := loadValues("values.override.yaml")
//	merged := yaml.MergeTrees(values, overrides)
//
// # Performance Considerations
//
//   - For files <10MB: Use UnmarshalYAML
//   - For files >100MB: Use StreamParser
//   - For multiple operations: Clone nodes to preserve originals
//   - For large merges: Consider memory usage and use streaming if needed
//
// # Thread Safety
//
// This library is NOT thread-safe. Use external synchronization for
// concurrent access to the same NodeTree or Node instances.
//
// # Version Compatibility
//
// Tested with Go versions: 1.21, 1.22, 1.23, 1.24
// Minimum Go version: 1.21
//
// For more examples and detailed documentation, see:
//   - examples/ directory in the repository
//   - docs/ directory for architecture and advanced features
//   - README.md for quick start and overview
//
// Source: https://github.com/elioetibr/golang-yaml-advanced
package yaml
