package main

import (

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Contract) MintToken(client Client, tokenURI string) (*types.Transaction, error) {
	tx, err := c.Instance.MintNFT(client.Auth, tokenURI)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create signed mint transaction")
	}

	return tx, nil
}
