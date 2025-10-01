package main

import (
	matp_contract "aztec/matp"
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatal(err)
		return
	}

	atpAddress := "0x5af96494ee0aa3921e84fcad3b38233a07257c57"
	matp, err := matp_contract.NewMatp(common.HexToAddress(atpAddress), client)
	if err != nil {
		log.Fatal(err)
	}

	callOpts := &bind.CallOpts{Context: context.Background()}

	claimable, err := matp.GetClaimable(callOpts)
	if err != nil {
		log.Fatal(err)
		return
	}
	allocation, err := matp.GetAllocation(callOpts)
	if err != nil {
		log.Fatal(err)
		return
	}
	claimed, err := matp.GetClaimed(callOpts)
	if err != nil {
		log.Fatal(err)
		return
	}

	token, err := matp.GetToken(callOpts)
	if err != nil {
		log.Fatal(err)
		return
	}

	beneficiary, err := matp.GetBeneficiary(callOpts)
	if err != nil {
		log.Fatal(err)
		return
	}

	isRevoked, err := matp.GetIsRevoked(callOpts)
	if err != nil {
		log.Fatal(err)
		return
	}

	locked := new(big.Int).Sub(allocation, claimable)

	fmt.Println("Token:", token.String())
	fmt.Println("IsRevoked:", isRevoked)
	fmt.Println("Beneficiary:", beneficiary)
	fmt.Println("Claimable:", WeiToETH(claimable))
	fmt.Println("Claimed:", WeiToETH(claimed))
	// fmt.Println("Staked:", staked)
	fmt.Println("Locked:", WeiToETH(locked))
}

// Converts wei to ETH as float64
func WeiToETH(amount *big.Int) float64 {
	fAmount := new(big.Float).SetInt(amount)
	ethValue := new(big.Float).Quo(fAmount, big.NewFloat(1e18))
	result, _ := ethValue.Float64()
	return result
}
