// Package datadog is used to talk to datadog agents running in our Kubernetes cluster.
package datadog

import (
	"testing"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/stretchr/testify/assert"
)

type MonitorComparer func(expected, got *Monitor) assert.Comparison

func CompareMonitorNamespace(expected, got *Monitor) assert.Comparison {
	return func() (success bool) {
		return expected.client.Namespace == got.client.Namespace
	}
}

// Just returns a client without an err
func NewStatsdNew(addr string, options ...statsd.Option) *statsd.Client {
	client, err := statsd.New(addr, options...)
	if err != nil {
		panic(err)
	}
	return client
}

func Test_newDDClient(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name        string
		args        args
		want        *Monitor
		comparisons []MonitorComparer
	}{
		{
			name: "should be able to configure a namespace",
			args: args{
				opts: []Option{App("hello.")},
			},
			want: &Monitor{
				client: NewStatsdNew("127.0.0.1:8125", statsd.WithNamespace("hello.")),
			},
			comparisons: []MonitorComparer{
				CompareMonitorNamespace,
			},
		},
		{
			name: "should be able to configure a namespace and append a . if missing",
			args: args{
				opts: []Option{App("hello")},
			},
			want: &Monitor{
				client: NewStatsdNew("127.0.0.1:8125", statsd.WithNamespace("hello.")),
			},
			comparisons: []MonitorComparer{
				CompareMonitorNamespace,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newDDClient(tt.args.opts...)

			for _, comparision := range tt.comparisons {
				assert.Condition(t, comparision(tt.want, got))
			}
		})
	}
}
