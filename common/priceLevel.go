package common

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

// 重新定义标准包
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ErrPriceLevel error 类型
var ErrPriceLevel = errors.New("failed to parse price level")

// PriceLevel 是订单簿中买价和卖价的通用结构
type PriceLevel struct {
	Price    float64 // 价格
	Quantity float64 // 数量
}

// UnmarshalJSON 解释JSON
func (p *PriceLevel) UnmarshalJSON(data []byte) error {
	var items []string
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	// 必须同时返回价格 和 数量
	if len(items) != 2 {
		return ErrPriceLevel
	}

	price, err := strconv.ParseFloat(items[0], 64) // 转成float64
	if err != nil {
		return err
	}
	quantity, err := strconv.ParseFloat(items[1], 64) // 转成float64
	if err != nil {
		return err
	}

	p.Price = price
	p.Quantity = quantity
	return nil
}

func (p *PriceLevel) MarshalJSON() ([]byte, error) {
	items := [2]string{}
	items[0] = strconv.FormatFloat(p.Price, 'f', -1, 64)
	items[1] = strconv.FormatFloat(p.Quantity, 'f', -1, 64)

	return json.Marshal(items)
}
