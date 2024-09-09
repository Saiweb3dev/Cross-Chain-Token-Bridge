package services

import (
	"backend/config"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func ProcessContractEvent(eventData interface{}) error {
	connect,err := config.GetEthereumConnection()
	if err != nil {
		return err
	}

	contract,err := NewMyContract(common.HexToAddress("contract_address"),connect)
	if err != nil {
		return err
	}

	contractFilterer, err := contract.FilterMyEvent(&bind.FilterOpts{})
	if err != nil {
		return err
	}

	for event := range contractFilterer.Logs {

	}
	return nil
}

type MyContract struct{
	//Contract ABI
}

func NewMyContract(address common.Address, connect *ethclient.Client) (*MyContract, error) {
	// Create a new contract instance
}