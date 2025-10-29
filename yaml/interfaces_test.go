package yaml

import (
	"errors"
	"strings"
	"testing"
)

// TestInterfacesImplementation verifies that concrete types implement their interfaces
func TestInterfacesImplementation(t *testing.T) {
	// Test Schema implements SchemaValidator
	schema := &Schema{Type: "string"}
	var _ SchemaValidator = schema

	// Test TransformDSL implements TreeTransformer
	dsl := NewTransformDSL()
	var _ TreeTransformer = dsl

	// Test StreamParser implements DocumentParser
	parser := NewStreamParser(strings.NewReader(""))
	var _ DocumentParser = parser
}

// MockSchemaValidator demonstrates dependency injection with mocked validator
type MockSchemaValidator struct {
	ValidateFn func(node *Node, path string) []ValidationError
	CallCount  int
}

func (m *MockSchemaValidator) Validate(node *Node, path string) []ValidationError {
	m.CallCount++
	if m.ValidateFn != nil {
		return m.ValidateFn(node, path)
	}
	return nil
}

// TestSchemaValidatorInterface demonstrates using SchemaValidator interface for DI
func TestSchemaValidatorInterface(t *testing.T) {
	// Create a function that depends on SchemaValidator interface, not concrete *Schema
	validateWithInterface := func(validator SchemaValidator, node *Node) []ValidationError {
		return validator.Validate(node, "$")
	}

	t.Run("with real Schema implementation", func(t *testing.T) {
		schema := &Schema{
			Type: "string",
		}
		node := NewScalarNode("test")

		errors := validateWithInterface(schema, node)
		if len(errors) != 0 {
			t.Errorf("Expected no errors, got %d", len(errors))
		}
	})

	t.Run("with mock implementation", func(t *testing.T) {
		mock := &MockSchemaValidator{
			ValidateFn: func(node *Node, path string) []ValidationError {
				return []ValidationError{
					{Path: path, Message: "mock error"},
				}
			},
		}

		node := NewScalarNode("test")
		errors := validateWithInterface(mock, node)

		if mock.CallCount != 1 {
			t.Errorf("Expected 1 call, got %d", mock.CallCount)
		}

		if len(errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(errors))
		}

		if errors[0].Message != "mock error" {
			t.Errorf("Expected 'mock error', got %q", errors[0].Message)
		}
	})
}

// MockTreeTransformer demonstrates dependency injection with mocked transformer
type MockTreeTransformer struct {
	ApplyFn   func(tree *NodeTree) (*NodeTree, error)
	CallCount int
}

func (m *MockTreeTransformer) Apply(tree *NodeTree) (*NodeTree, error) {
	m.CallCount++
	if m.ApplyFn != nil {
		return m.ApplyFn(tree)
	}
	return tree, nil
}

// TestTreeTransformerInterface demonstrates using TreeTransformer interface for DI
func TestTreeTransformerInterface(t *testing.T) {
	// Create a function that depends on TreeTransformer interface, not concrete *TransformDSL
	transformWithInterface := func(transformer TreeTransformer, tree *NodeTree) (*NodeTree, error) {
		return transformer.Apply(tree)
	}

	t.Run("with real TransformDSL implementation", func(t *testing.T) {
		dsl := NewTransformDSL().RemoveKey("remove_me")

		tree := NewNodeTree()
		doc := &Document{Root: NewMappingNode()}
		doc.Root.AddKeyValue(NewScalarNode("keep"), NewScalarNode("value1"))
		doc.Root.AddKeyValue(NewScalarNode("remove_me"), NewScalarNode("value2"))
		tree.Documents = append(tree.Documents, doc)

		result, err := transformWithInterface(dsl, tree)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result == nil {
			t.Fatal("Expected non-nil result")
		}
	})

	t.Run("with mock implementation", func(t *testing.T) {
		mock := &MockTreeTransformer{
			ApplyFn: func(tree *NodeTree) (*NodeTree, error) {
				return nil, errors.New("mock transformation error")
			},
		}

		tree := NewNodeTree()
		_, err := transformWithInterface(mock, tree)

		if mock.CallCount != 1 {
			t.Errorf("Expected 1 call, got %d", mock.CallCount)
		}

		if err == nil {
			t.Error("Expected error from mock")
		}

		if err.Error() != "mock transformation error" {
			t.Errorf("Expected 'mock transformation error', got %q", err.Error())
		}
	})
}

// MockDocumentParser demonstrates dependency injection with mocked parser
type MockDocumentParser struct {
	ParseFn               func() error
	SetDocumentCallbackFn func(callback func(*NodeTree) error)
	CallCount             int
	Callback              func(*NodeTree) error
}

func (m *MockDocumentParser) Parse() error {
	m.CallCount++
	if m.ParseFn != nil {
		return m.ParseFn()
	}
	return nil
}

func (m *MockDocumentParser) SetDocumentCallback(callback func(*NodeTree) error) {
	m.Callback = callback
	if m.SetDocumentCallbackFn != nil {
		m.SetDocumentCallbackFn(callback)
	}
}

// TestDocumentParserInterface demonstrates using DocumentParser interface for DI
func TestDocumentParserInterface(t *testing.T) {
	// Create a function that depends on DocumentParser interface, not concrete *StreamParser
	parseWithInterface := func(parser DocumentParser, callback func(*NodeTree) error) error {
		parser.SetDocumentCallback(callback)
		return parser.Parse()
	}

	t.Run("with real StreamParser implementation", func(t *testing.T) {
		yamlData := "key: value"
		parser := NewStreamParser(strings.NewReader(yamlData))

		callbackCalled := false
		callback := func(tree *NodeTree) error {
			callbackCalled = true
			if tree == nil {
				t.Error("Expected non-nil tree in callback")
			}
			return nil
		}

		err := parseWithInterface(parser, callback)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !callbackCalled {
			t.Error("Callback was not called")
		}
	})

	t.Run("with mock implementation", func(t *testing.T) {
		mock := &MockDocumentParser{}
		mock.ParseFn = func() error {
			// Simulate calling the callback during parsing
			if mock.Callback != nil {
				tree := NewNodeTree()
				return mock.Callback(tree)
			}
			return nil
		}

		callbackCalled := false
		callback := func(tree *NodeTree) error {
			callbackCalled = true
			return errors.New("callback error")
		}

		err := parseWithInterface(mock, callback)

		if mock.CallCount != 1 {
			t.Errorf("Expected 1 call, got %d", mock.CallCount)
		}

		if err == nil {
			t.Error("Expected error from callback")
		}

		if !callbackCalled {
			t.Error("Callback was not called")
		}
	})
}

// TestDependencyInjectionBenefits demonstrates the testability benefits of interfaces
func TestDependencyInjectionBenefits(t *testing.T) {
	// Example: A pipeline that uses multiple interfaces
	type ValidationPipeline struct {
		Validator   SchemaValidator
		Transformer TreeTransformer
		Parser      DocumentParser
	}

	executePipeline := func(pipeline ValidationPipeline, yamlData string) error {
		// Parse
		err := pipeline.Parser.Parse()
		if err != nil {
			return err
		}

		// Validate (would need tree from parser in real scenario)
		node := NewScalarNode("test")
		errors := pipeline.Validator.Validate(node, "$")
		if len(errors) > 0 {
			return errors[0]
		}

		// Transform (would use tree from parser in real scenario)
		tree := NewNodeTree()
		_, err = pipeline.Transformer.Apply(tree)
		return err
	}

	t.Run("with all mock implementations", func(t *testing.T) {
		// This test demonstrates we can test the pipeline logic
		// without depending on real YAML parsing, validation, or transformation
		pipeline := ValidationPipeline{
			Validator: &MockSchemaValidator{
				ValidateFn: func(node *Node, path string) []ValidationError {
					return nil // Mock: always valid
				},
			},
			Transformer: &MockTreeTransformer{
				ApplyFn: func(tree *NodeTree) (*NodeTree, error) {
					return tree, nil // Mock: pass-through
				},
			},
			Parser: &MockDocumentParser{
				ParseFn: func() error {
					return nil // Mock: successful parse
				},
			},
		}

		err := executePipeline(pipeline, "mock: data")
		if err != nil {
			t.Errorf("Pipeline failed: %v", err)
		}
	})

	t.Run("with mixed real and mock implementations", func(t *testing.T) {
		// This demonstrates we can test with some real and some mock components
		pipeline := ValidationPipeline{
			Validator:   &Schema{Type: "string"}, // Real implementation
			Transformer: NewTransformDSL(),       // Real implementation
			Parser: &MockDocumentParser{ // Mock implementation
				ParseFn: func() error {
					return nil
				},
			},
		}

		err := executePipeline(pipeline, "mock: data")
		if err != nil {
			t.Errorf("Pipeline failed: %v", err)
		}
	})
}

// BenchmarkInterfaceVsConcrete verifies interface overhead is minimal
func BenchmarkInterfaceVsConcrete(b *testing.B) {
	schema := &Schema{Type: "string"}
	node := NewScalarNode("test")

	b.Run("direct_call", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(node, "$")
		}
	})

	b.Run("interface_call", func(b *testing.B) {
		var validator SchemaValidator = schema
		for i := 0; i < b.N; i++ {
			_ = validator.Validate(node, "$")
		}
	})
}
