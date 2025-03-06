package remover

type Remover interface {
	Process(lines []string) ([]string, []string)
}

var registry = map[string]Remover{
	"go":         &GoRemover{},
	"python":     &PythonRemover{},
	"dockerfile": &DockerfileRemover{},
}

func ForType(fileType string) Remover {
	return registry[fileType]
}

func Register(fileType string, remover Remover) {
	registry[fileType] = remover
}
