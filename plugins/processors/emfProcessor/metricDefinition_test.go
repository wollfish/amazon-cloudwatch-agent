package emfProcessor

import (
	"reflect"
	"testing"

	"github.com/aws/amazon-cloudwatch-agent/internal/structuredlogscommon"
	"github.com/stretchr/testify/assert"
)

func buildTestMetricDeclaration() *metricDeclaration {
	md := metricDeclaration{
		SourceLabels:    []string{"tagA", "tagB"},
		LabelMatcher:    "^v1;v2$",
		MetricSelectors: []string{"metric_a", "metric_b"},
		Dimensions:      [][]string{{"tagB", "tagA"}, {"tagA"}, {"tag2", "tag1"}},
	}
	md.init()
	return &md
}

func Test_metricDeclaration_init(t *testing.T) {
	md := buildTestMetricDeclaration()
	assert.Equal(t, ";", md.LabelSeparator)
	assert.Equal(t, [][]string{{"tagA", "tagB"}, {"tagA"}, {"tag1", "tag2"}}, md.Dimensions)
}

func Test_getConcatenatedLabels_complete(t *testing.T) {
	md := buildTestMetricDeclaration()

	metricTags := map[string]string{"tagA": "valueA", "tagB": "valueB"}
	result := md.getConcatenatedLabels(metricTags)

	assert.Equal(t, "valueA;valueB", result)
}

func Test_getConcatenatedLabels_incomplete(t *testing.T) {
	md := buildTestMetricDeclaration()

	metricTags := map[string]string{"tagC": "valueA", "tagD": "valueB"}
	result := md.getConcatenatedLabels(metricTags)

	assert.Equal(t, ";", result)
}

func Test_process_match(t *testing.T) {
	md := buildTestMetricDeclaration()

	metricTags := map[string]string{"tagA": "v1", "tagB": "v2"}
	metricFields := map[string]interface{}{"metric_a": "valueA", "metric_c": 10.0}
	result := md.process(metricTags, metricFields, "ContainerInsights/Prometheus")
	assert.Equal(t, "ContainerInsights/Prometheus", result.Namespace)
	assert.True(t, reflect.DeepEqual([][]string{{"tagA", "tagB"}, {"tagA"}, {"tag1", "tag2"}}, result.DimensionSets))
	assert.True(t, reflect.DeepEqual([]structuredlogscommon.MetricAttr{structuredlogscommon.MetricAttr{Name: "metric_a"}},
		result.Metrics))
}

func Test_process_mismatch(t *testing.T) {
	md := buildTestMetricDeclaration()

	metricTags := map[string]string{"tagA": "v1", "tagC": "v3"}
	metricFields := map[string]interface{}{"metric_a": "valueA", "metric_b": 10.0, "metric_c": 10.0}
	result := md.process(metricTags, metricFields, "ContainerInsights/Prometheus")

	assert.True(t, reflect.ValueOf(result).IsNil())
}
