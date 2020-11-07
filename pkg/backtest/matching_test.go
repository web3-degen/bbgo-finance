package backtest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/c9s/bbgo/pkg/fixedpoint"
	"github.com/c9s/bbgo/pkg/types"
)

func newLimitOrder(symbol string, side types.SideType, price, quantity float64) types.SubmitOrder {
	return types.SubmitOrder{
		Symbol:      symbol,
		Side:        side,
		Type:        types.OrderTypeLimit,
		Quantity:    quantity,
		Price:       price,
		TimeInForce: "GTC",
	}
}

func TestSimplePriceMatching(t *testing.T) {
	engine := &SimplePriceMatching{
		CurrentTime: time.Now(),
		OrderID:     1,
	}

	for i := 0; i < 5; i++ {
		_, _, err := engine.PlaceOrder(newLimitOrder("BTCUSDT", types.SideTypeBuy, 8000.0-float64(i), 1.0))
		assert.NoError(t, err)
	}
	assert.Len(t, engine.bidOrders, 5)
	assert.Len(t, engine.askOrders, 0)

	for i := 0; i < 5; i++ {
		_, _, err := engine.PlaceOrder(newLimitOrder("BTCUSDT", types.SideTypeSell, 9000.0+float64(i), 1.0))
		assert.NoError(t, err)
	}
	assert.Len(t, engine.bidOrders, 5)
	assert.Len(t, engine.askOrders, 5)

	closedOrders, trades := engine.SellToPrice(fixedpoint.NewFromFloat(8100.0))
	assert.Len(t, closedOrders, 0)
	assert.Len(t, trades, 0)

	closedOrders, trades = engine.SellToPrice(fixedpoint.NewFromFloat(8000.0))
	assert.Len(t, closedOrders, 1)
	assert.Len(t, trades, 1)
	for _, o := range closedOrders {
		assert.Equal(t, types.SideTypeBuy, o.Side)
	}

	closedOrders, trades = engine.SellToPrice(fixedpoint.NewFromFloat(7000.0))
	assert.Len(t, closedOrders, 4)
	assert.Len(t, trades, 4)

	closedOrders, trades = engine.BuyToPrice(fixedpoint.NewFromFloat(8900.0))
	assert.Len(t, closedOrders, 0)
	assert.Len(t, trades, 0)

	closedOrders, trades = engine.BuyToPrice(fixedpoint.NewFromFloat(9000.0))
	assert.Len(t, closedOrders, 1)
	assert.Len(t, trades, 1)
	for _, o := range closedOrders {
		assert.Equal(t, types.SideTypeSell, o.Side)
	}

	for _, trade := range trades {
		assert.Equal(t, types.SideTypeSell, trade.Side)
	}

	closedOrders, trades = engine.BuyToPrice(fixedpoint.NewFromFloat(9500.0))
	assert.Len(t, closedOrders, 4)
	assert.Len(t, trades, 4)
}
