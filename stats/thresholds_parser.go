package stats

import (
	"fmt"

	c "go.k6.io/k6/pkg/combinators"
)

func ParseAggregationMethod() c.Parser {
	parser := c.Expect(c.Alternative(
		ParseCounter(),
		ParseGauge(),
		ParseRate(),
		ParseTrend(),
		ParsePercentile(),
	), "aggregation method")

	return func(input []rune) c.Result {
		res := parser(input)
		if res.Err != nil {
			return res
		}

		return c.Success(res.Payload, res.Remaining)
	}
}

func ParseOperator() c.Parser {
	parser := c.Expect(c.Alternative(
		c.Tag(">="),
		c.Tag("<="),
		c.Tag(">"),
		c.Tag("<"),
		c.Tag("==="),
		c.Tag("=="),
		c.Tag("!="),
	))

	return func(input []rune) c.Result {
		res := parser(input)
		if res.Err != nil {
			return res
		}

		return c.Success(res.Payload.(string), res.Remaining)
	}
}

func ParsePercentile() c.Parser {
	parser := c.Expect(c.Sequence(
		c.Tag("p("),
		c.Float(),
		c.Char(')'),
	))

	return func(input []rune) c.Result {
		res := parser(input)
		if res.Err != nil {
			return res
		}

		parsed := res.Payload.([]interface{})
		percentile := fmt.Sprintf("%s%g%s", parsed[0].(string), parsed[1].(float64), parsed[2].(string))

		return c.Success(percentile, res.Remaining)
	}
}

func ParseTrend() c.Parser {
	parser := c.Expect(c.Alternative(
		c.Tag("mean"),
		c.Tag("min"),
		c.Tag("max"),
		c.Tag("avg"),
		ParsePercentile(), // FIXME: not happy to return mixed types here, especially in Go
	))

	return func(input []rune) c.Result {
		res := parser(input)
		if res.Err != nil {
			return res
		}

		return c.Success(res.Payload, res.Remaining)
	}
}

func ParseRate() c.Parser {
	parser := c.Expect(c.Tag("rate"))

	return func(input []rune) c.Result {
		res := parser(input)
		if res.Err != nil {
			return res
		}

		return c.Success(res.Payload, res.Remaining)
	}
}

func ParseGauge() c.Parser {
	parser := c.Expect(c.Alternative(
		c.Tag("last"),
		c.Tag("min"),
		c.Tag("max"),
		c.Tag("value"),
	))

	return func(input []rune) c.Result {
		res := parser(input)
		if res.Err != nil {
			return res
		}

		return c.Success(res.Payload, res.Remaining)
	}
}

func ParseCounter() c.Parser {
	parser := c.Expect(c.Alternative(
		c.Tag("count"),
		c.Tag("sum"),
		c.Tag("rate"),
	))

	return func(input []rune) c.Result {
		res := parser(input)
		if res.Err != nil {
			return res
		}

		return c.Success(res.Payload, res.Remaining)
	}
}

func ParseAssertion(aggregation c.Parser) c.Parser {
	parser := c.Expect(c.Sequence(
		aggregation,
		ParseOperator(),
		c.Float(),
	))

	return func(input []rune) c.Result {
		res := parser(input)
		if res.Err != nil {
			return res
		}

		return c.Success(res.Payload, res.Remaining)
	}
}
