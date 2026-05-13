package schema

// Scaffold is the top-level structure parsed from a scaffold.yaml file.
type Scaffold struct {
	Name         string       `yaml:"name"`
	Description  string       `yaml:"description"`
	Version      string       `yaml:"version"`
	Vars         []Var        `yaml:"vars"`
	Files        []File       `yaml:"files"`
	Dependencies []Dependency `yaml:"dependencies"`
	Hooks        Hooks        `yaml:"hooks"`
}

type Var struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Type        string    `yaml:"type"` // string | bool | select | multiselect
	Required    bool      `yaml:"required"`
	Prompt      string    `yaml:"prompt"`
	Options     []Option  `yaml:"options"`
	Default     any       `yaml:"default"`
	Validate    *Validate `yaml:"validate,omitempty"`
}

type Option struct {
	Label string `yaml:"label"`
	Value string `yaml:"value"`
}

type Validate struct {
	Pattern string `yaml:"pattern"`
	Message string `yaml:"message"`
}

type File struct {
	Src       string `yaml:"src"`
	Dst       string `yaml:"dst"`
	Condition string `yaml:"condition,omitempty"`
}

type Dependency struct {
	Pkg       string `yaml:"pkg"`
	Condition string `yaml:"condition,omitempty"`
}

type Hooks struct {
	Pre  []Hook `yaml:"pre"`
	Post []Hook `yaml:"post"`
}

type Hook struct {
	Cmd       string `yaml:"cmd"`
	Condition string `yaml:"condition,omitempty"`
}
