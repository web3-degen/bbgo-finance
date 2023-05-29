package indicator

import (
	"github.com/c9s/bbgo/pkg/datatype/floats"
	"github.com/c9s/bbgo/pkg/types"
)

/*
NEW INDICATOR DESIGN:

klines := kLines(marketDataStream)
closePrices := closePrices(klines)
macd := MACD(klines, {Fast: 12, Slow: 10})

equals to:

klines := KLines(marketDataStream)
closePrices := ClosePrice(klines)
fastEMA := EMA(closePrices, 7)
slowEMA := EMA(closePrices, 25)
macd := Subtract(fastEMA, slowEMA)
signal := EMA(macd, 16)
histogram := Subtract(macd, signal)
*/

type Float64Source interface {
	types.Series
	OnUpdate(f func(v float64))
}

type Float64Subscription interface {
	types.Series
	AddSubscriber(f func(v float64))
}

//go:generate callbackgen -type EWMAStream
type EWMAStream struct {
	Float64Updater
	types.SeriesBase

	slice floats.Slice

	window     int
	multiplier float64
}

func EWMA2(source Float64Source, window int) *EWMAStream {
	s := &EWMAStream{
		window:     window,
		multiplier: 2.0 / float64(1+window),
	}

	s.SeriesBase.Series = s.slice

	if sub, ok := source.(Float64Subscription); ok {
		sub.AddSubscriber(s.calculateAndPush)
	} else {
		source.OnUpdate(s.calculateAndPush)
	}

	return s
}

func (s *EWMAStream) calculateAndPush(v float64) {
	v2 := s.calculate(v)
	s.slice.Push(v2)
	s.EmitUpdate(v2)
}

func (s *EWMAStream) calculate(v float64) float64 {
	last := s.slice.Last()
	m := s.multiplier
	return (1.0-m)*last + m*v
}
