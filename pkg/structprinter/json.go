package structprinter

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrintJson(data interface{}, pretty bool) error {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if pretty {
		var jsonPretty bytes.Buffer
		err = json.Indent(&jsonPretty, jsonBytes, "", "\t")

		if err != nil {
			return err
		}

		fmt.Println(jsonPretty.String())
	} else {
		fmt.Println(string(jsonBytes))
	}

	return nil
}
