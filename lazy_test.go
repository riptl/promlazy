package promlazy

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLazy(t *testing.T) {
	registry := prometheus.NewRegistry()
	// Register a lazy batch.
	batch := With(registry)
	metric1 := batch.NewGauge(prometheus.GaugeOpts{Name: "my_metric_1"})
	metric2 := batch.NewCounter(prometheus.CounterOpts{Name: "my_metric_2"})
	// Gather before writing to lazy metrics.
	// We expect to gather no metrics.
	gather1, err := registry.Gather()
	require.NoError(t, err)
	assert.Len(t, gather1, 0)
	// Write a value.
	metric1.Set(42)
	// Gather metrics again. The previous write should have registered metrics.
	gather2, err := registry.Gather()
	require.NoError(t, err)
	assert.Len(t, gather2, 2)
	// Write a bunch of values again to prove idempotency.
	metric1.Inc()
	metric2.Inc()
	metric2.Add(3)
	metric1.Dec()
	gather3, err := registry.Gather()
	require.NoError(t, err)
	assert.Len(t, gather3, 2)
}
