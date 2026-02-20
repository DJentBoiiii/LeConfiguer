package commands

import (
	"encoding/json"
	"fmt"

	"client/internal/client"
)

type configMeta struct {
	ID, Name, Type, Environment string
}

func requireID(id string) error {
	if id == "" {
		return fmt.Errorf("id flag required: -i")
	}
	return nil
}

func requireAll(id, name, typ, env string) error {
	if id == "" || name == "" || typ == "" || env == "" {
		return fmt.Errorf("all flags required: -i, -n, -t, -e")
	}
	return nil
}

func printResult(msg string, data []byte) {
	fmt.Printf("%s\n%s\n", msg, data)
}

func getMeta(c *client.Client, id string) (*configMeta, error) {
	data, err := c.Get(fmt.Sprintf("/configs/%s", id))
	if err != nil {
		return nil, err
	}
	var m configMeta
	return &m, json.Unmarshal(data, &m)
}
