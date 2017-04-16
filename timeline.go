package tachymeter

import (
	"fmt"
	"time"
	"bytes"
	"strconv"
	"strings"
	"encoding/json"
	"io/ioutil"
)

type Timeline struct {
	timeline []*TimelineEvent
}

type TimelineEvent struct {
	Metrics *Metrics
	Created time.Time
}

func (t *Timeline) AddEvent(m *Metrics) {
	t.timeline = append(t.timeline, &TimelineEvent{
		Metrics: m,
		Created: time.Now(),
	})
}

func (t *Timeline) WriteHtml() {
	var b bytes.Buffer

	b.WriteString(head)

	// Create a canvast for each timeline event.
	b.WriteString(fmt.Sprintf(`%s<div id="container">%s`, tab, nl))
	for n := range t.timeline {
		b.WriteString(fmt.Sprintf(`%s%s<canvas id="canvas-%d"></canvas>%s`, tab, tab, n, nl))
	}
	b.WriteString(fmt.Sprintf(`%s<div id="container">`, tab))

	// Write graphs.
	for id, m := range t.timeline {
		s := getGraphHtml(m, id)
		b.WriteString(s)
	}

	b.WriteString(tail)

	d := []byte(b.String())
	fname := fmt.Sprintf("tachymeter-%d.html", time.Now().Unix())
	err := ioutil.WriteFile(fname, d, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func getGraphHtml(te *TimelineEvent, id int) string {
	keys := []string{}
	values := []int{}

	for _, b := range te.Metrics.Histogram {
		for k, v := range b {
			keys = append(keys, k)
			values = append(values, v)
		}
	}

	keysj, _ := json.Marshal(keys)
	valuesj, _ := json.Marshal(values)

	out := strings.Replace(graph, "XCANVASID", strconv.Itoa(id), 1)
	out = strings.Replace(out, "XKEYS", string(keysj), 1)
	out = strings.Replace(out, "XVALUES", string(valuesj), 1)

	return out
}