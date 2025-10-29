package yaml

// Package interfaces provides minimal abstractions for dependency inversion.
// Following SOLID principles (Dependency Inversion) as specified in CLAUDE.md.
//
// Design Philosophy (Option A - Minimal):
// - Interfaces focused on behavior, not data structures
// - Enable mocking for tests without breaking existing APIs
// - Keep it simple (KISS principle from CLAUDE.md)

// SchemaValidator validates YAML structures against schemas.
// Implementations: *Schema
type SchemaValidator interface {
	Validate(node *Node, path string) []ValidationError
}

// TreeTransformer applies transformations to YAML trees.
// Implementations: *TransformDSL
type TreeTransformer interface {
	Apply(tree *NodeTree) (*NodeTree, error)
}

// DocumentParser parses YAML documents in streaming or batch mode.
// Implementations: *StreamParser
type DocumentParser interface {
	Parse() error
	SetDocumentCallback(callback func(*NodeTree) error)
}

// Ensure concrete types implement their respective interfaces at compile time
var (
	_ SchemaValidator = (*Schema)(nil)
	_ TreeTransformer = (*TransformDSL)(nil)
	_ DocumentParser  = (*StreamParser)(nil)
)
