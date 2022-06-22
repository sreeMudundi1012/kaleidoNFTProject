package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Contract) BurnToken(client Client, tokenID *big.Int) (*types.Transaction, error) {
	tx, err := c.Instance.BurnNFT(client.Auth, tokenID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create signed burn transaction")
	}

	return tx, nil
}
