// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type testConfigCollection int

const (
	testSetDefault testConfigCollection = iota
	testSetAll
	testSetNone
)

func TestMetricsBuilder(t *testing.T) {
	tests := []struct {
		name      string
		configSet testConfigCollection
	}{
		{
			name:      "default",
			configSet: testSetDefault,
		},
		{
			name:      "all_set",
			configSet: testSetAll,
		},
		{
			name:      "none_set",
			configSet: testSetNone,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			start := pcommon.Timestamp(1_000_000_000)
			ts := pcommon.Timestamp(1_000_001_000)
			observedZapCore, observedLogs := observer.New(zap.WarnLevel)
			settings := receivertest.NewNopCreateSettings()
			settings.Logger = zap.New(observedZapCore)
			mb := NewMetricsBuilder(loadConfig(t, test.name), settings, WithStartTime(start))

			expectedWarnings := 0
			if test.configSet == testSetDefault {
				assert.Equal(t, "[WARNING] Please set `enabled` field explicitly for `default.metric`: This metric will be disabled by default soon.", observedLogs.All()[expectedWarnings].Message)
				expectedWarnings++
			}
			if test.configSet == testSetDefault || test.configSet == testSetAll {
				assert.Equal(t, "[WARNING] `default.metric.to_be_removed` should not be enabled: This metric is deprecated and will be removed soon.", observedLogs.All()[expectedWarnings].Message)
				expectedWarnings++
			}
			if test.configSet == testSetAll || test.configSet == testSetNone {
				assert.Equal(t, "[WARNING] `optional.metric` should not be configured: This metric is deprecated and will be removed soon.", observedLogs.All()[expectedWarnings].Message)
				expectedWarnings++
			}
			assert.Equal(t, expectedWarnings, observedLogs.Len())

			defaultMetricsCount := 0
			allMetricsCount := 0

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordDefaultMetricDataPoint(ts, 1, "attr-val", 1, AttributeEnumAttr(1), []any{"one", "two"}, map[string]any{"onek": "onev", "twok": "twov"})

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordDefaultMetricToBeRemovedDataPoint(ts, 1)

			allMetricsCount++
			mb.RecordOptionalMetricDataPoint(ts, 1, "attr-val", true)

			metrics := mb.Emit(WithMapResourceAttr(map[string]any{"onek": "onev", "twok": "twov"}), WithOptionalResourceAttr("attr-val"), WithSliceResourceAttr([]any{"one", "two"}), WithStringEnumResourceAttrOne, WithStringResourceAttr("attr-val"))

			if test.configSet == testSetNone {
				assert.Equal(t, 0, metrics.ResourceMetrics().Len())
				return
			}

			assert.Equal(t, 1, metrics.ResourceMetrics().Len())
			rm := metrics.ResourceMetrics().At(0)
			attrCount := 0
			enabledAttrCount := 0
			attrVal, ok := rm.Resource().Attributes().Get("map.resource.attr")
			attrCount++
			assert.Equal(t, mb.resourceAttributesSettings.MapResourceAttr.Enabled, ok)
			if mb.resourceAttributesSettings.MapResourceAttr.Enabled {
				enabledAttrCount++
				assert.EqualValues(t, map[string]any{"onek": "onev", "twok": "twov"}, attrVal.Map().AsRaw())
			}
			attrVal, ok = rm.Resource().Attributes().Get("optional.resource.attr")
			attrCount++
			assert.Equal(t, mb.resourceAttributesSettings.OptionalResourceAttr.Enabled, ok)
			if mb.resourceAttributesSettings.OptionalResourceAttr.Enabled {
				enabledAttrCount++
				assert.EqualValues(t, "attr-val", attrVal.Str())
			}
			attrVal, ok = rm.Resource().Attributes().Get("slice.resource.attr")
			attrCount++
			assert.Equal(t, mb.resourceAttributesSettings.SliceResourceAttr.Enabled, ok)
			if mb.resourceAttributesSettings.SliceResourceAttr.Enabled {
				enabledAttrCount++
				assert.EqualValues(t, []any{"one", "two"}, attrVal.Slice().AsRaw())
			}
			attrVal, ok = rm.Resource().Attributes().Get("string.enum.resource.attr")
			attrCount++
			assert.Equal(t, mb.resourceAttributesSettings.StringEnumResourceAttr.Enabled, ok)
			if mb.resourceAttributesSettings.StringEnumResourceAttr.Enabled {
				enabledAttrCount++
				assert.Equal(t, "one", attrVal.Str())
			}
			attrVal, ok = rm.Resource().Attributes().Get("string.resource.attr")
			attrCount++
			assert.Equal(t, mb.resourceAttributesSettings.StringResourceAttr.Enabled, ok)
			if mb.resourceAttributesSettings.StringResourceAttr.Enabled {
				enabledAttrCount++
				assert.EqualValues(t, "attr-val", attrVal.Str())
			}
			assert.Equal(t, enabledAttrCount, rm.Resource().Attributes().Len())
			assert.Equal(t, attrCount, 5)

			assert.Equal(t, 1, rm.ScopeMetrics().Len())
			ms := rm.ScopeMetrics().At(0).Metrics()
			if test.configSet == testSetDefault {
				assert.Equal(t, defaultMetricsCount, ms.Len())
			}
			if test.configSet == testSetAll {
				assert.Equal(t, allMetricsCount, ms.Len())
			}
			validatedMetrics := make(map[string]bool)
			for i := 0; i < ms.Len(); i++ {
				switch ms.At(i).Name() {
				case "default.metric":
					assert.False(t, validatedMetrics["default.metric"], "Found a duplicate in the metrics slice: default.metric")
					validatedMetrics["default.metric"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Monotonic cumulative sum int metric enabled by default.", ms.At(i).Description())
					assert.Equal(t, "s", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("string_attr")
					assert.True(t, ok)
					assert.EqualValues(t, "attr-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("state")
					assert.True(t, ok)
					assert.EqualValues(t, 1, attrVal.Int())
					attrVal, ok = dp.Attributes().Get("enum_attr")
					assert.True(t, ok)
					assert.Equal(t, "red", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("slice_attr")
					assert.True(t, ok)
					assert.EqualValues(t, []any{"one", "two"}, attrVal.Slice().AsRaw())
					attrVal, ok = dp.Attributes().Get("map_attr")
					assert.True(t, ok)
					assert.EqualValues(t, map[string]any{"onek": "onev", "twok": "twov"}, attrVal.Map().AsRaw())
				case "default.metric.to_be_removed":
					assert.False(t, validatedMetrics["default.metric.to_be_removed"], "Found a duplicate in the metrics slice: default.metric.to_be_removed")
					validatedMetrics["default.metric.to_be_removed"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "[DEPRECATED] Non-monotonic delta sum double metric enabled by default.", ms.At(i).Description())
					assert.Equal(t, "s", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityDelta, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeDouble, dp.ValueType())
					assert.Equal(t, float64(1), dp.DoubleValue())
				case "optional.metric":
					assert.False(t, validatedMetrics["optional.metric"], "Found a duplicate in the metrics slice: optional.metric")
					validatedMetrics["optional.metric"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "[DEPRECATED] Gauge double metric disabled by default.", ms.At(i).Description())
					assert.Equal(t, "1", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeDouble, dp.ValueType())
					assert.Equal(t, float64(1), dp.DoubleValue())
					attrVal, ok := dp.Attributes().Get("string_attr")
					assert.True(t, ok)
					assert.EqualValues(t, "attr-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("boolean_attr")
					assert.True(t, ok)
					assert.EqualValues(t, true, attrVal.Bool())
				}
			}
		})
	}
}

func loadConfig(t *testing.T, name string) MetricsBuilderConfig {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	sub, err := cm.Sub(name)
	require.NoError(t, err)
	cfg := DefaultMetricsBuilderConfig()
	require.NoError(t, component.UnmarshalConfig(sub, &cfg))
	return cfg
}
