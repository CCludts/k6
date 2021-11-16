package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperator(t *testing.T) {
	// Arrange
	parser := ParseOperator()

	// Act
	greaterOrEqualThan := parser([]rune(">="))
	lowerOrEqualThan := parser([]rune("<="))
	greaterThan := parser([]rune(">"))
	lowerThan := parser([]rune("<"))
	strictlyEqual := parser([]rune("==="))
	looselyEqual := parser([]rune("=="))
	notEqual := parser([]rune("!="))

	// Assert
	assert.Equal(t, ">=", greaterOrEqualThan.Payload)
	assert.Equal(t, "", string(greaterOrEqualThan.Remaining))
	assert.Equal(t, "<=", lowerOrEqualThan.Payload)
	assert.Equal(t, "", string(lowerOrEqualThan.Remaining))
	assert.Equal(t, ">", greaterThan.Payload)
	assert.Equal(t, "", string(greaterThan.Remaining))
	assert.Equal(t, "<", lowerThan.Payload)
	assert.Equal(t, "", string(lowerThan.Remaining))
	assert.Equal(t, "===", strictlyEqual.Payload)
	assert.Equal(t, "", string(strictlyEqual.Remaining))
	assert.Equal(t, "==", looselyEqual.Payload)
	assert.Equal(t, "", string(looselyEqual.Remaining))
	assert.Equal(t, "!=", notEqual.Payload)
	assert.Equal(t, "", string(notEqual.Remaining))
}

func TestPercentile_ValueWithoutDecimal(t *testing.T) {
	// Arrange
	parser := ParsePercentile()

	// Act
	result := parser([]rune("p(99)"))

	// Assert
	assert.Equal(t, string("p(99)"), result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func TestPercentile_ValueWithDecimal(t *testing.T) {
	// Arrange
	parser := ParsePercentile()

	// Act
	result := parser([]rune("p(99.9)"))

	// Assert
	assert.Equal(t, string("p(99.9)"), result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func TestTrend(t *testing.T) {
	// Arrange
	parser := ParseTrend()

	// Act
	mean := parser([]rune("mean"))
	min := parser([]rune("min"))
	max := parser([]rune("max"))
	avg := parser([]rune("avg"))
	p99 := parser([]rune("p(99.9)"))

	// Assert
	assert.Equal(t, "mean", mean.Payload)
	assert.Equal(t, "", string(mean.Remaining))
	assert.Equal(t, "min", min.Payload)
	assert.Equal(t, "", string(min.Remaining))
	assert.Equal(t, "max", max.Payload)
	assert.Equal(t, "", string(max.Remaining))
	assert.Equal(t, "avg", avg.Payload)
	assert.Equal(t, "", string(avg.Remaining))
	assert.Equal(t, float64(99.9), p99.Payload)
	assert.Equal(t, "", string(p99.Remaining))
}

func TestRate(t *testing.T) {
	// Arrange
	parser := ParseRate()

	// Act
	result := parser([]rune("rate"))

	// Assert
	assert.Equal(t, "rate", result.Payload)
	assert.Equal(t, "", string(result.Remaining))
}

func TestGauge(t *testing.T) {
	// Arrange
	parser := ParseGauge()

	// Act
	last := parser([]rune("last"))
	min := parser([]rune("min"))
	max := parser([]rune("max"))
	value := parser([]rune("value"))

	// Assert
	assert.Equal(t, "last", last.Payload)
	assert.Equal(t, "", string(last.Remaining))
	assert.Equal(t, "min", min.Payload)
	assert.Equal(t, "", string(min.Remaining))
	assert.Equal(t, "max", max.Payload)
	assert.Equal(t, "", string(max.Remaining))
	assert.Equal(t, "value", value.Payload)
	assert.Equal(t, "", string(value.Remaining))
}

func TestCounter(t *testing.T) {
	// Arrange
	parser := ParseCounter()

	// Act
	count := parser([]rune("count"))
	sum := parser([]rune("sum"))
	rate := parser([]rune("rate"))

	// Assert
	assert.Equal(t, "count", count.Payload)
	assert.Equal(t, "", string(count.Remaining))
	assert.Equal(t, "sum", sum.Payload)
	assert.Equal(t, "", string(sum.Remaining))
	assert.Equal(t, "rate", rate.Payload)
	assert.Equal(t, "", string(rate.Remaining))
}

func TestAssertion(t *testing.T) {
	// Arrange
	trendAssertionParser := ParseAssertion(ParseTrend())
	rateAssertionParser := ParseAssertion(ParseRate())
	// FIXME: see next FIXME
	// gaugeAssertionParser := ParseAssertion(ParseGauge())
	counterAssertionParser := ParseAssertion(ParseCounter())

	// Act
	trendAssertionResult := trendAssertionParser([]rune("p(99.9)<300"))
	rateAssertionResult := rateAssertionParser([]rune("rate>0.95"))
	// FIXME: see next FIXME
	// gaugeAssertionResult := gaugeAssertionParser([]rune("value<4000"))
	counterAssertionResult := counterAssertionParser([]rune("count<100"))

	// Assert
	assert.Equal(t, []interface{}{float64(99.9), "<", float64(300)}, trendAssertionResult.Payload)
	assert.Equal(t, "", string(trendAssertionResult.Remaining))
	assert.Equal(t, []interface{}{"rate", ">", float64(0.95)}, rateAssertionResult.Payload)
	assert.Equal(t, "", string(rateAssertionResult.Remaining))
	// FIXME: this test returns a []int32 Remaining for some reason.
	// assert.Equal(t, []interface{}{"value", "<", float64(4000)}, gaugeAssertionResult)
	// assert.Equal(t, "", string(gaugeAssertionResult.Remaining))
	assert.Equal(t, []interface{}{"count", "<", float64(100)}, counterAssertionResult.Payload)
}
