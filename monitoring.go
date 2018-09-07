// Package monitoring is used to monitor applications running in our internal Kubernetes clusters.
package monitoring

import (
	"github.com/DataDog/datadog-go/statsd"
	"github.com/pborman/uuid"
	"net"
)

// Monitor is used to send metrics to a metric collection service.
type Monitor struct {
	id  string
	app string

	datadog dataDogAgent
}

type dataDogAgent struct {
	host   string
	port   string
	client *statsd.Client
}

const (
	defaulDataDogHost  = "datadog-agent"
	defaultDataDogPort = "8125"
)

// Option used to construct a new Monitor.
type Option func(*Monitor) error

// NewMonitor creates a monitor. It will panic if there is a problem creating it.
func NewMonitor(opts ...Option) *Monitor {
	m := &Monitor{
		id: uuid.New(),
	}
	for _, opt := range opts {
		err := opt(m)
		if err != nil {
			panic(err)
		}
	}
	var err error
	m.datadog.client, err = newDataDogClient(m.datadog.host, m.datadog.port)
	if err != nil {
		panic(err)
	}
	m.datadog.tag(m.app, m.id)
	if m.app != "" {
		t := m.app + " monitor started."
		m.datadog.client.SimpleEvent(t, t)
	}
	return m
}

// App sets the name of the application using the monitor.
func App(name string) Option {
	return func(m *Monitor) error {
		m.app = name
		return nil
	}
}

func newDataDogClient(host, port string) (*statsd.Client, error) {
	h, p := host, port
	if h == "" {
		h = defaulDataDogHost
	}
	if p == "" {
		p = defaultDataDogPort
	}
	c, err := statsd.New(net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (d *dataDogAgent) tag(app, id string) {
	if app != "" {
		d.client.Tags = append(d.client.Tags, "app:"+app)
	}
	d.client.Tags = append(d.client.Tags, "monitor-id:"+id)
}

// Count can be incremented or decremented.
func (m *Monitor) Count(name string, value int64, tags []string, rate float64) error {
	return m.datadog.client.Count(name, value, tags, rate)
}

// Decr a count.
func (m *Monitor) Decr(name string, tags []string, rate float64) error {
	return m.datadog.client.Decr(name, tags, rate)
}

// Incr a count.
func (m *Monitor) Incr(name string, tags []string, rate float64) error {
	return m.datadog.client.Incr(name, tags, rate)
}

// Event marks an event.
func (m *Monitor) Event(title, text string) error {
	return m.datadog.client.SimpleEvent(title, text)
}

// Gauge is used to set a metric to a specific value. It will stay at that value until changed.
func (m *Monitor) Gauge(name string, value float64, tags []string, rate float64) error {
	return m.datadog.client.Gauge(name, value, tags, rate)
}
