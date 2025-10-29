package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/elioetibr/golang-yaml-advanced/yaml"
)

func main() {
	// Get directory of main.go
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// Step 1: Load the base values file
	baseTreeFile := filepath.Join(dir, "yaml-values.yaml")
	fmt.Printf("Step 1: Loading {%s}...\n", baseTreeFile)
	baseTree, err := loadYAML(filepath.Join(dir, "yaml-values.yaml"))
	if err != nil {
		log.Fatalf("Failed to load base values: %v", err)
	}
	fmt.Printf("  ✓ Loaded base values")

	// Step 2: Load the override values file
	overrideTreeFile := filepath.Join(dir, "yaml-values-override.yaml")
	fmt.Printf("\nStep 2: Loading {%s}...", overrideTreeFile)
	overrideTree, err := loadYAML(overrideTreeFile)
	if err != nil {
		log.Fatalf("Failed to load override values: %v", err)
	}
	fmt.Printf("  ✓ Loaded override values")

	// Step 3: Merge the files (override values take precedence)
	fmt.Println("\nStep 3: Merging YAML files...")
	mergedTree := yaml.MergeTrees(baseTree, overrideTree)
	fmt.Printf("  ✓ Merged YAML trees")

	// Show what was overridden
	showOverrides(overrideTree)

	// Step 4: Write the merged result
	mergedTreeFile := filepath.Join(dir, "yaml-values-merged.yaml")
	fmt.Printf("\nStep 4: Writing merged result to {%s}...", mergedTreeFile)
	if err := writeYAML(mergedTreeFile, mergedTree); err != nil {
		log.Fatalf("Failed to write merged values: %v", err)
	}
	fmt.Println("  ✓ Successfully wrote merged YAML")

	// Step 5: Verify the output - compare with existing merged file
	fmt.Println("\nStep 5: Verifying output...")
	existingData, err := os.ReadFile(mergedTreeFile)
	if err != nil {
		log.Fatalf("Failed to read merged file for verification: %v", err)
	}

	// Re-parse to verify
	verifiedTree, err := yaml.UnmarshalYAML(existingData)
	if err != nil {
		log.Fatalf("Failed to parse merged file for verification: %v", err)
	}

	// Compare the trees
	differences := yaml.DiffTrees(mergedTree, verifiedTree)
	if len(differences) == 0 {
		fmt.Println("  ✓ Output matches expected result - no diff!")
	} else {
		fmt.Printf("  ⚠ Found %d differences\n", len(differences))
		for i, diff := range differences {
			if i < 3 { // Show only first 3 differences
				fmt.Printf("    - Path: %s, Type: %s, Old: %v, New: %v\n",
					diff.Path, diff.Type, diff.OldValue, diff.NewValue)
			}
		}
		if len(differences) > 3 {
			fmt.Printf("    ... and %d more\n", len(differences)-3)
		}
	}

	fmt.Println("\n✨ Process completed successfully!")
}

// loadYAML loads and parses a YAML file using the custom yaml package
func loadYAML(filename string) (*yaml.NodeTree, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	tree, err := yaml.UnmarshalYAML(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return tree, nil
}

// writeYAML writes a NodeTree to a YAML file
func writeYAML(filename string, tree *yaml.NodeTree) error {
	yamlContent, err := tree.ToYAML()
	if err != nil {
		return fmt.Errorf("failed to convert to YAML: %w", err)
	}

	return os.WriteFile(filename, yamlContent, 0644)
}

// showOverrides displays what values are being overridden
func showOverrides(tree *yaml.NodeTree) {
	fmt.Println("\n  Override summary:")

	if len(tree.Documents) > 0 && tree.Documents[0].Root != nil {
		root := tree.Documents[0].Root

		// Look for specific overrides
		for _, child := range root.Children {
			if child.Key != nil {
				keyValue := fmt.Sprintf("%v", child.Key.Value)
				switch keyValue {
				case "replicaCount":
					if child.Value != nil {
						fmt.Printf("    • replicaCount: %v (increased replicas)\n", child.Value)
					}
				case "topologySpreadConstraints":
					fmt.Printf("    • topologySpreadConstraints: configured")
				case "strategy":
					fmt.Printf("    • strategy: configured")
				}
			}
		}
	}
}
