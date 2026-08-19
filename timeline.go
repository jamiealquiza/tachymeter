package tachymeter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Timeline holds a series of *timelineEvents,
// which nest *Metrics for analyzing multiple
// collections of measured events.
type Timeline struct {
	timeline []*timelineEvent
}

// timelineEvent holds a *Metrics and the
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

// WriteHTML takes a path p and writes an html file to
// 'p/tachymeter-<timestamp>.html' of all histograms
// held by the *Timeline, in series.
func (t *Timeline) WriteHTML(p string) error {
	path, err := filepath.Abs(p)
	if err != nil {
		return err
	}

	var b strings.Builder

	b.WriteString(head)

	// Append graph + info entry for
	// each timeline event.
	for n, e := range t.timeline {
		// Graph div.
		fmt.Fprintf(&b, `%s<div class="graph">%s`, tab, nl)
		fmt.Fprintf(&b, `%s%s<canvas id="canvas-%d"></canvas>%s`, tab, tab, n, nl)
		fmt.Fprintf(&b, `%s</div>%s`, tab, nl)
		// Info div.
		fmt.Fprintf(&b, `%s<div class="info">%s`, tab, nl)
		fmt.Fprintf(&b, `%s<p><h2>Iteration %d</h2>%s`, tab, n+1, nl)
		b.WriteString(e.Metrics.String())
		fmt.Fprintf(&b, "%s%s</p></div>%s", nl, tab, nl)
	}

	// Write graphs.
	for id, e := range t.timeline {
		b.WriteString(genGraphHTML(e, id))
	}

	b.WriteString(tail)

	fname := filepath.Join(path, fmt.Sprintf("tachymeter-%d.html", time.Now().Unix()))

	return os.WriteFile(fname, []byte(b.String()), 0644)
}

// genGraphHTML takes a *timelineEvent and id (used for each graph
// html element ID) and creates a chart.js graph output.
func genGraphHTML(te *timelineEvent, id int) string {
	keys := []string{}
	values := []uint64{}

	for _, bin := range *te.Metrics.Histogram {
		for k, v := range bin {
			keys = append(keys, k)
			values = append(values, v)
		}
	}

	keysj, _ := json.Marshal(keys)
	valuesj, _ := json.Marshal(values)

	r := strings.NewReplacer(
		"XCANVASID", strconv.Itoa(id),
		"XKEYS", string(keysj),
		"XVALUES", string(valuesj),
	)

	return r.Replace(graph)
}
