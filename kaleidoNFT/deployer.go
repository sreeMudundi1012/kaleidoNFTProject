package main
	 
import (
  "github.com/ethereum/go-ethereum/common"
  "github.com/pkg/errors"
)

type Contract struct {
  Address common.Address
  Instance *Main
}

func (c *Client) DeployContract() (Contract, error) {
  contractAddress, _, instance, err := DeployMain(c.Auth, c.Rpc)
  if err != nil {
    return Contract{}, errors.Wrap(err, "failed to deploy the contract")
  }

  return Contract{contractAddress, instance}, nil
}