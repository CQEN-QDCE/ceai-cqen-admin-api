package structprinter

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

func PrintYaml(data interface{}) error {
	yamlBytes, err := yaml.Marshal(&data)

	if err != nil {
		return err
	}

	fmt.Print(string(yamlBytes))

	return nil
}
