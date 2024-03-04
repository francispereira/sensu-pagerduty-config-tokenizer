package main

import (
	"testing"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
}

func TestExecuteMutator(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	configValue := `{"foo":"||.bar||","foo1":"||.bar1||"}`
	event.Check.Annotations = map[string]string{
		"sensu.io/plugins/sensu-pagerduty-handler/config/summary-template":   configValue,
		"sensu.io/plugins/sensu-pagerduty-handler/config/details-template":   configValue,
		"sensu.io/plugins/sensu-pagerduty-handler/config/dedup-key-template": configValue,
	}
	event.Metrics = corev2.FixtureMetrics()
	ev, err := executeMutator(event)
	assert.NoError(err)
	assert.Equal(`{"foo":"{{.bar}}","foo1":"{{.bar1}}"}`, ev.Check.Annotations["sensu.io/plugins/sensu-pagerduty-handler/config/summary-template"])
	assert.Equal(`{"foo":"{{.bar}}","foo1":"{{.bar1}}"}`, ev.Check.Annotations["sensu.io/plugins/sensu-pagerduty-handler/config/details-template"])
	assert.Equal(`{"foo":"{{.bar}}","foo1":"{{.bar1}}"}`, ev.Check.Annotations["sensu.io/plugins/sensu-pagerduty-handler/config/dedup-key-template"])

}
