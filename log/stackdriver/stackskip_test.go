package stackdriver

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/TV4/logrus-stackdriver-formatter/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestStackSkip(t *testing.T) {
	var out bytes.Buffer

	logger := logrus.New()
	logger.Out = &out
	logger.Formatter = NewFormatter(
		WithService("test"),
		WithVersion("0.1"),
		WithStackSkip("github.com/TV4/logrus-stackdriver-formatter/test"),
	)

	mylog := test.LogWrapper{
		Logger: logger,
	}

	mylog.Error("my log entry")

	want := map[string]interface{}{
		"severity": "ERROR",
		"message":  "my log entry",
		"serviceContext": map[string]interface{}{
			"service": "test",
			"version": "0.1",
		},
		"context": map[string]interface{}{
			"reportLocation": map[string]interface{}{
				"file":     "github.com/armory-io/monitoring/log/stackdriver/formatter.go",
				"line":     264.0,
				"function": "(*Formatter).Format",
			},
		},
		"sourceLocation": map[string]interface{}{
			"file":     "github.com/armory-io/monitoring/log/stackdriver/formatter.go",
			"line":     264.0,
			"function": "(*Formatter).Format",
		},
	}


	// remove timestamp from actual result
	var actual map[string]interface{}
	json.Unmarshal(out.Bytes(), &actual)
	delete(actual, "timestamp")
	actualWithoutTimestamp, err := json.Marshal(actual)

	expected, err := json.Marshal(want)
	if err != nil {
		t.Error(err)
	}
	assert.JSONEq(t, string(expected), string(actualWithoutTimestamp))
}
