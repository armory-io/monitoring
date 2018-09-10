// Package datadog is used to talk to datadog agents running in our Kubernetes cluster.
package datadog

import (
	"github.com/DataDog/datadog-go/statsd"
	"github.com/pborman/uuid"
	"net"
)

// Monitor is used to send metrics to a metric collection service.
type Monitor struct {
	id  string
	app string

	host   string
	port   string
	client *statsd.Client

	logger logger
	debug  bool
}

const (
	defaulDataDogHost  = "datadog-agent"
	defaultDataDogPort = "8125"
)

type logger interface {
	Printf(format string, args ...interface{})
}

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
	m.client, err = newClient(m.host, m.port)
	if err != nil {
		panic(err)
	}
	if m.app != "" {

	}
	m.client.Tags = append(m.client.Tags, "monitor-id:"+m.id)
	if m.app != "" {
		m.client.Tags = append(m.client.Tags, "app:"+m.app)
		t := m.app + " monitor started."
		m.client.SimpleEvent(t, t)
	}
	m.log("Created datadog monitor: %v", *m)
	return m
}

// App sets the name of the application using the monitor.
func App(name string) Option {
	return func(m *Monitor) error {
		m.app = name
		return nil
	}
}

// WithLogger sets the logger to be used by the monitor.
func WithLogger(l logger) Option {
	return func(m *Monitor) error {
		m.logger = l
		return nil
	}
}

// Debug mode for the monitor.
func Debug() Option {
	return func(m *Monitor) error {
		m.debug = true
		return nil
	}
}

func newClient(host, port string) (*statsd.Client, error) {
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

func (m *Monitor) log(format string, args ...interface{}) {
	if m.logger != nil {
		m.logger.Printf("datadog.monitor: "+format, args)
	}
}

func (m *Monitor) error(err error) {
	if m.debug == true && err != nil {
		m.log("error: %v", err)
	}
}
