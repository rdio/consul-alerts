package notifier

import (
	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AlertaNotifier struct {
	Url             string
	DefaultSchedule string
	Environment     string
	Schedules       map[string]string
	Nodes           map[string]string
	Services        map[string]string
}

type PagerDutyKeys struct {
	Keys []string `json:"pdkeys"`
}

type AlertaNotification struct {
	Resource    string        `json:"resource"`
	Event       string        `json:"event"`
	Environment string        `json:"environment"`
	Severity    string        `json:"severity"`
	Status      string        `json:"status"`
	Services    []string      `json:"service"`
	Value       string        `json:"value"`
	Text        string        `json:"text"`
	Attributes  PagerDutyKeys `json:"attributes"`
}

func getAlertaStatusAndSeverity(status string) (string, string) {
	switch status {
	case "critical":
		return "open", "critical"
	case "passing":
		return "closed", "ok"
	case "warning":
		return "open", "warning"
	default:
		return "open", "warning"
	}
}

func (al *AlertaNotifier) ProcessSchedule(str string, src map[string]string, m map[string]string) {
	a := strings.Split(str, ",")
	for _, v := range a {
		if val, ok := src[strings.TrimSpace(v)]; ok {
			if hash, ok := al.Schedules[val]; ok {
				m[val] = hash
			}
		}
	}
}

func (al *AlertaNotifier) Notify(messages Messages) bool {
	result := true

	for _, message := range messages {
		// map of schedule to pagerduty keys
		// i.e. operations => <hash>
		m := make(map[string]string)

		service := message.Service
		node := message.Node

		// populate map
		al.ProcessSchedule(service, al.Services, m)
		al.ProcessSchedule(node, al.Nodes, m)

		// add default schedule if map is empty
		if len(m) == 0 {
			m[al.DefaultSchedule] = al.Schedules[al.DefaultSchedule]
		}

		// convert map to array of pagerduty keys
		a := []string{}
		for _, v := range m {
			a = append(a, v)
		}

		// set status and severity
		status, severity := getAlertaStatusAndSeverity(message.Status)

		// create notification
		an := AlertaNotification{
			Environment: al.Environment,
			Services:    []string{message.Service},
			Resource:    message.Node,
			Event:       message.Status,
			Value:       message.Status,
			Status:      status,
			Severity:    severity,
			Attributes: PagerDutyKeys{
				Keys: a,
			},
			Text: message.Output,
		}

		// post to server
		if jsonStr, err := json.Marshal(an); err == nil {
			fmt.Println("POSTing to:", al.Url)
			req, err := http.NewRequest("POST", al.Url, bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("ERROR: ", resp.Status)
			}
			defer resp.Body.Close()

			fmt.Println("response Status:", resp.Status)
		}
	}
	log.Println("Alerta notification complete")
	return result
}
