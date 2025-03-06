package remover

// Remover defines the interface for comment removal.
type Remover interface {
	// Process removes comments from the given lines of code.
	// It returns the processed lines and a list of issues found.
	Process(lines []string) ([]string, []string)
}

// registry of removers by file type
var registry = map[string]Remover{
	"go":     &GoRemover{},
	"python": &PythonRemover{},
}

// ForType returns a Remover for the given file type.
func ForType(fileType string) Remover {
	return registry[fileType]
}

// Register adds a remover to the registry.
func Register(fileType string, remover Remover) {
	registry[fileType] = remover
}
