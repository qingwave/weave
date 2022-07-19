package trace

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"
)

var (
	logMsg     = ""
	testLogger = funcr.New(func(prefix, args string) {
		logMsg = fmt.Sprintln(prefix, args)
	}, funcr.Options{})
)

func TestTrace(t *testing.T) {
	trace := New("test", testLogger, Field{Key: "test", Value: "val"})

	trace.Step("step1")
	trace.Step("step2")

	trace.Log()

	assert.Equal(t, 2, len(trace.traceItems))
	assert.Equal(t, "step1", trace.traceItems[0].(traceStep).msg)

	assert.NotZero(t, trace.TotalTime())
	assert.Contains(t, logMsg, "step1", "step2")

	// log long not output
	{
		logMsg = ""
		trace.LogIfLong(1 * time.Second)
		assert.Empty(t, logMsg)
	}

	// log long
	{
		time.Sleep(10 * time.Millisecond)
		trace.Step("long running step")
		logMsg = ""
		trace.LogIfLong(10 * time.Millisecond)

		assert.NotEmpty(t, logMsg)
	}
}

func TestNestTrace(t *testing.T) {
	trace := New("test", testLogger)

	new := trace.Nest("nest1")

	logMsg = ""
	trace.Log()

	assert.Len(t, trace.traceItems, 1)
	assert.Equal(t, trace.traceItems[0], new)
	assert.NotEmpty(t, logMsg)
}
