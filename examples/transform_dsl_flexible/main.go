package main

import (
	"fmt"
	"log"

	"github.com/elioetibr/golang-yaml-advanced/yaml"
)

func main() {
	fmt.Println("=== TransformDSL with Flexible Input Demo ===")

	// Example 1: Creating TransformDSL from YAML string
	example1()

	// Example 2: Creating TransformDSL from YAML bytes
	example2()

	// Example 3: Creating TransformDSL from Go map
	example3()

	// Example 4: Creating TransformDSL from Node
	example4()

	// Example 5: Creating TransformDSL from NodeTree
	example5()

	// Example 6: Chaining transformations
	example6()
}

func example1() {
	fmt.Println("Example 1: From YAML String")
	yamlStr := `
name: Alice
role: admin
email: alice@example.com
department: Engineering
`

	// Create DSL from string - new pattern
	dsl, err := yaml.NewTransformDSLFrom(yamlStr)
	if err != nil {
		log.Fatalf("Failed to create DSL: %v", err)
	}

	// Apply transformation without passing tree (uses stored input)
	result, err := dsl.RemoveKey("email").Apply(nil)
	if err != nil {
		log.Fatalf("Failed to apply transformation: %v", err)
	}

	output, _ := result.ToYAML()
	fmt.Printf("Original had email field, transformed output:\n%s\n", string(output))
}

func example2() {
	fmt.Println("Example 2: From YAML Bytes")
	yamlBytes := []byte(`
config:
  database:
    host: localhost
    port: 5432
  cache:
    enabled: true
    ttl: 3600
`)

	dsl, err := yaml.NewTransformDSLFrom(yamlBytes)
	if err != nil {
		log.Fatalf("Failed to create DSL: %v", err)
	}

	result, err := dsl.
		RemoveKey("cache").
		RenameKey("database", "db").
		Apply(nil)
	if err != nil {
		log.Fatalf("Failed to apply transformation: %v", err)
	}

	output, _ := result.ToYAML()
	fmt.Printf("Removed cache, renamed database to db:\n%s\n", string(output))
}

func example3() {
	fmt.Println("Example 3: From Go Map")
	configMap := map[string]interface{}{
		"server": map[string]interface{}{
			"host": "0.0.0.0",
			"port": 8080,
		},
		"debug": true,
	}

	dsl, err := yaml.NewTransformDSLFrom(configMap)
	if err != nil {
		log.Fatalf("Failed to create DSL: %v", err)
	}

	result, err := dsl.
		AddComment("Configuration settings").
		SortKeys().
		Apply(nil)
	if err != nil {
		log.Fatalf("Failed to apply transformation: %v", err)
	}

	output, _ := result.ToYAML()
	fmt.Printf("Sorted keys and added comments:\n%s\n", string(output))
}

func example4() {
	fmt.Println("Example 4: From Node")
	// Create a mapping node manually
	node := yaml.NewMappingNode()
	node.AddKeyValue(yaml.NewScalarNode("app"), yaml.NewScalarNode("myapp"))
	node.AddKeyValue(yaml.NewScalarNode("version"), yaml.NewScalarNode("1.0.0"))
	node.AddKeyValue(yaml.NewScalarNode("author"), yaml.NewScalarNode("John Doe"))

	dsl, err := yaml.NewTransformDSLFrom(node)
	if err != nil {
		log.Fatalf("Failed to create DSL: %v", err)
	}

	result, err := dsl.
		RenameKey("author", "maintainer").
		Apply(nil)
	if err != nil {
		log.Fatalf("Failed to apply transformation: %v", err)
	}

	output, _ := result.ToYAML()
	fmt.Printf("Renamed author to maintainer:\n%s\n", string(output))
}

func example5() {
	fmt.Println("Example 5: From NodeTree")
	yamlData := `
users:
  - name: Alice
    admin: true
  - name: Bob
    admin: false
`
	tree, err := yaml.UnmarshalYAML([]byte(yamlData))
	if err != nil {
		log.Fatalf("Failed to unmarshal: %v", err)
	}

	dsl, err := yaml.NewTransformDSLFrom(tree)
	if err != nil {
		log.Fatalf("Failed to create DSL: %v", err)
	}

	result, err := dsl.
		AddComment("User list").
		Apply(nil)
	if err != nil {
		log.Fatalf("Failed to apply transformation: %v", err)
	}

	output, _ := result.ToYAML()
	fmt.Printf("Added comment to user list:\n%s\n", string(output))
}

func example6() {
	fmt.Println("Example 6: Complex Chaining")
	yamlStr := `
application:
  name: myapp
  version: 1.0.0
  temp_field: to_be_removed
  database:
    host: localhost
    port: 5432
    password: secret123
`

	dsl, err := yaml.NewTransformDSLFrom(yamlStr)
	if err != nil {
		log.Fatalf("Failed to create DSL: %v", err)
	}

	// Chain multiple transformations
	result, err := dsl.
		RemoveKey("temp_field").
		RemoveKey("password").
		RenameKey("database", "db_config").
		SortKeys().
		Apply(nil)
	if err != nil {
		log.Fatalf("Failed to apply transformation: %v", err)
	}

	output, _ := result.ToYAML()
	fmt.Printf("Cleaned up config (removed temp fields and password, renamed database):\n%s\n", string(output))
}
