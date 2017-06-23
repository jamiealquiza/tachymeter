package tachymeter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Timeline holds a []*timelineEvents,
// which nest *Metrics for analyzing
// multiple collections of measured events.
type Timeline struct {
	timeline []*timelineEvent
}

// timelineEvent holds a *Metrics and
// time that it was added to the Timeline.
type timelineEvent struct {
	Metrics *Metrics
	Created time.Time
}

// AddEvent adds a *Metrics to the *Timeline.
func (t *Timeline) AddEvent(m *Metrics) {
	t.timeline = append(t.timeline, &timelineEvent{
		Metrics: m,
		Created: time.Now(),
	})
}

// WriteHTML takes an absolute path p and writes an
// html file to 'p/tachymeter-<timestamp>.html' of all
// histograms held by the *Timeline, in series.
func (t *Timeline) WriteHTML(p string) error {
	path, err := filepath.Abs(p)
	if err != nil {
		return err
	}
	var b bytes.Buffer

	b.WriteString(head)

	// Append graph + info entry for each timeline
	// event.
	for n := range t.timeline {
		// Graph div.
		b.WriteString(fmt.Sprintf(`%s<div class="graph">%s`, tab, nl))
		b.WriteString(fmt.Sprintf(`%s%s<canvas id="canvas-%d"></canvas>%s`, tab, tab, n, nl))
		b.WriteString(fmt.Sprintf(`%s</div>%s`, tab, nl))
		// Info div.
		b.WriteString(fmt.Sprintf(`%s<div class="info">%s`, tab, nl))
		b.WriteString(fmt.Sprintf(`%s<p><h2>Iteration %d</h2>%s`, tab, n+1, nl))
		b.WriteString(t.timeline[n].Metrics.String())
		b.WriteString(fmt.Sprintf("%s%s</p></div>%s", nl, tab, nl))
	}

	// Write graphs.
	for id, m := range t.timeline {
		s := genGraphHTML(m, id)
		b.WriteString(s)
	}

	b.WriteString(tail)

	// Write file.
	d := []byte(b.String())
	fname := fmt.Sprintf("%s/tachymeter-%d.html", path, time.Now().Unix())
	err = ioutil.WriteFile(fname, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

// genGraphHTML takes a *timelineEvent and id (used for each graph
// html element ID) and creates a chart.js graph output.
func genGraphHTML(te *timelineEvent, id int) string {
	keys := []string{}
	values := []uint64{}

	for _, b := range *te.Metrics.Histogram {
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
