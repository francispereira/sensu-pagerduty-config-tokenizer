package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

// Config represents the mutator plugin config.
type Config struct {
	sensu.PluginConfig
}

var (
	mutatorConfig = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-pagerduty-mutator",
			Short:    "Sensu Pagerduty Mutator",
			Keyspace: "sensu.io/plugins/sensu-pagerduty-mutator/config",
		},
	}
)

func main() {
	mutator := sensu.NewGoMutator(&mutatorConfig.PluginConfig, nil, checkArgs, executeMutator)
	mutator.Execute()
}

func checkArgs(_ *types.Event) error {
	return nil
}

func handleError(message string, err error) {
	fmt.Printf("%s: %s\n", message, err)
	os.Exit(1)
}
func executeMutator(event *types.Event) (*types.Event, error) {

	namespace, checkName := "default", event.Check.Name
	url := fmt.Sprintf("%s/api/core/v2/namespaces/%s/checks/%s", os.Getenv("SENSU_API_URL"), namespace, checkName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		handleError("Error creating request", err)
	}
	req.Header.Add("Authorization", "Key "+os.Getenv("SENSU_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		handleError("Error making request", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError("Error reading response body", err)
	}

	var result map[string]json.RawMessage
	err = json.Unmarshal(body, &result)
	if err != nil {
		handleError("Error parsing JSON response", err)
	}

	var metadataDetails map[string]interface{}
	err = json.Unmarshal(result["metadata"], &metadataDetails)
	if err != nil {
		handleError("Error extracting metadata from response", err)
	}

	annotations, ok := metadataDetails["annotations"].(map[string]interface{})
	if !ok {
		handleError("Error extracting labels from metadata", fmt.Errorf("labels field not found"))
	}
	jsonAnnotations, err := json.Marshal(annotations)
	if err != nil {
		handleError("Error encoding labels to JSON", err)
	}
	var finalAnnotations map[string]string
	err = json.Unmarshal([]byte(jsonAnnotations), &finalAnnotations)

	for key, value := range finalAnnotations {
		finalAnnotations[key] = strings.ReplaceAll(value, "||.", "{{")
		finalAnnotations[key] = strings.ReplaceAll(value, "||", "}}")
	}
	event.Entity.Annotations = finalAnnotations

	return event, nil
}
