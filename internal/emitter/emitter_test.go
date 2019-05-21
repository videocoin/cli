package emitter

import (
	"math/big"
	"testing"
)

func TestConvertWeiToVDC(t *testing.T) {
	tables := []struct {
		wei string
		vdc string
	}{
		{"86544738311000000000", "86.54473831"},
		{"415311000000000", "0.000415311"},
		{"0", "0"},
	}

	for _, table := range tables {
		n := new(big.Int)
		n, ok := n.SetString(table.wei, 10)
		if !ok {
			t.Fail()
		}
		vdc, err := convertWeiToVDC(n)
		if err != nil || vdc.String() != table.vdc {
			t.Errorf("Convertion of %s wei is incorrect, got: %s, want: %s.", table.wei, vdc.String(), table.vdc)
		}
	}
}
