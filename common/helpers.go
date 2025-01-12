package common

import (
	"bytes"
	"math"
)

// AmountToLotSize converts an amount to a lot sized amount | 将金额转为数量
/**
计算手数：首先通过 amount / lot 来计算出原始金额可以分为多少个完整的 lot（手数）。
向下取整：使用 math.Floor 对手数进行向下取整，以确保不会超过原始金额所能允许的最大手数。
乘回手数：然后将取整后的手数乘以 lot，恢复成金额形式。
应用精度：使用 math.Pow10(precision) 将数值扩大 precision 位，再使用 math.Trunc 截断多余的尾数，最后除以 math.Pow10(precision) 还原到原来的规模，但此时已经按照指定的精度进行了舍入。
*/
func AmountToLotSize(lot float64, precision int, amount float64) float64 {
	return math.Trunc(math.Floor(amount/lot)*lot*math.Pow10(precision)) / math.Pow10(precision)
}

// ToJSONList convert v to json list if v is a map | 如果 v 是map，则将 v 转换为 json 列表
func ToJSONList(v []byte) []byte {
	if len(v) > 0 && v[0] == '{' {
		var b bytes.Buffer
		b.Write([]byte("["))
		b.Write(v)
		b.Write([]byte("]"))
		return b.Bytes()
	}
	return v
}
