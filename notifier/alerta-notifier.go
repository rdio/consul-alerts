package notifier

import (
	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/davecgh/go-spew/spew"
)

type AlertaNotifier struct {
	Url          string
	Schedules    map[string]string
	Environments map[string]string
	Nodes        map[string]string
	Services     map[string]string
}

type AlertaNotification struct {
	Resource    string            `json:"resource"`
	Event       string            `json:"event"`
	Environment string            `json:"environment"`
	Severity    string            `json:"severity"`
	Status      string            `json:"status"`
	Correlate   string            `json:"correlate"`
	Services    []string          `json:"service"`
	Group       string            `json:"group"`
	Value       string            `json:"value"`
	Text        string            `json:"text"`
	Tags        []string          `json:"tags"`
	Attributes  map[string]string `json:"attributes"`
	Origin      string            `json:"origin"`
	Type        string            `json:"type"`
}

func (al *AlertaNotifier) Notify(messages Messages) bool {
	result := true

	for _, message := range messages {
		log.Println(message)
		log.Println(al.Url)
	}

	spew.Dump(al.Schedules)
	spew.Dump(al.Environments)
	spew.Dump(al.Nodes)
	spew.Dump(al.Services)

	log.Println("Alerta notification complete")
	return result
}
