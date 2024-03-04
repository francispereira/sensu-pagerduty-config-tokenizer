package main

import (
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
			Name:     "sensu-pagerduty-config-tokenizer",
			Short:    "Sensu Pagerduty Config tokenizer",
			Keyspace: "sensu.io/plugins/sensu-pagerduty-config-tokenizer/config",
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

func executeMutator(event *types.Event) (*types.Event, error) {
	var sensuConfigProps = [3]string{"sensu.io/plugins/sensu-pagerduty-handler/config/summary-template",
		"sensu.io/plugins/sensu-pagerduty-handler/config/details-template",
		"sensu.io/plugins/sensu-pagerduty-handler/config/dedup-key-template",
	}
	for _, sensuConfigProp := range sensuConfigProps {
		if _, ok := event.Check.Annotations[sensuConfigProp]; ok {
			event.Check.Annotations[sensuConfigProp] = strings.ReplaceAll(event.Check.Annotations[sensuConfigProp], "||.", "{{.")
			event.Check.Annotations[sensuConfigProp] = strings.ReplaceAll(event.Check.Annotations[sensuConfigProp], "||", "}}")
		}
	}
	return event, nil
}
