package mongowebstat

import (
	"encoding/json"
	"io/ioutil"
)

func LoadConfig() error {
	jsonBlob, err := ioutil.ReadFile("./mongowebstat.json")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonBlob, &nodes); err != nil {
		return err
	}

	return nil
}
