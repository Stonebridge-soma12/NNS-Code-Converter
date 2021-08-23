package CodeGenerator

type Layer struct {
	Category string  `json:"category"`
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Input    *string `json:"input"`
	Output   *string `json:"output"`
	Param    Param   `json:"param"`
}

func (l *Layer) ToCode() (string, error) {
	var result string
	param, err := l.Param.ToCode(l.Type)

	result += l.Name
	result += " = "

	result += tf + keras + layers + "." + param
	if l.Input != nil {
		result += "(" + *l.Input + ")\n"
	} else {
		result += "\n"
	}

	return result, err
}