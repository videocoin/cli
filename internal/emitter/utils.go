package emitter

import (
	"math/big"
)

func convertWeiToVDC(wei *big.Int) (*big.Float, error) {
	var factor, exp = big.NewInt(18), big.NewInt(10)
	exp = exp.Exp(exp, factor, nil)

	fwei := new(big.Float).SetInt(wei)

	return new(big.Float).Quo(fwei, new(big.Float).SetInt(exp)), nil
}
