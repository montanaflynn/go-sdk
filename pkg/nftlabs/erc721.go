package nftlabs

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/montanaflynn/go-sdk/internal/abi"
)

type erc721 interface {
	commonModule
}

type erc721Module struct {
	Client  *ethclient.Client
	Address string
	Options *SdkOptions
	module  *abi.ERC721

	privateKey    *ecdsa.PrivateKey
	signerAddress common.Address
}

func newErc721SdkModule(client *ethclient.Client, address string, opt *SdkOptions) (*erc721Module, error) {
	module, err := abi.NewERC721(common.HexToAddress(address), client)
	if err != nil {
		// TODO: return better error
		return nil, err
	}

	return &erc721Module{
		Client:  client,
		Address: address,
		Options: opt,
		module:  module,
	}, nil
}

func (sdk *erc721Module) SetPrivateKey(privateKey string) error {
	if pKey, publicAddress, err := processPrivateKey(privateKey); err != nil {
		return &NoSignerError{typeName: "erc721", Err: err}
	} else {
		sdk.privateKey = pKey
		sdk.signerAddress = publicAddress
	}
	return nil
}
func (sdk *erc721Module) getSigner() func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
	return func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
		ctx := context.Background()
		chainId, _ := sdk.Client.ChainID(ctx)
		return types.SignTx(transaction, types.NewEIP155Signer(chainId), sdk.privateKey)
	}
}
