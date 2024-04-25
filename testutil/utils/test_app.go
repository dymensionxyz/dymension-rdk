package utils

import (
	"encoding/json"
	"testing"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"

	app "github.com/dymensionxyz/dymension-rdk/testutil/app"
	seqtypes "github.com/dymensionxyz/dymension-rdk/x/sequencers/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	// unnamed import of statik for swagger UI support
	_ "github.com/cosmos/cosmos-sdk/client/docs/statik"
)

var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   -1,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

var (
	ProposerPK       = simapp.CreateTestPubKeys(1)[0]
	ProposerConsAddr = sdk.ConsAddress(ProposerPK.Address())

	OperatorPK = secp256k1.GenPrivKey().PubKey()
)

func setup(withGenesis bool, invCheckPeriod uint) (*app.App, map[string]json.RawMessage) {
	db := dbm.NewMemDB()

	encCdc := app.MakeEncodingConfig()
	testApp := app.NewRollapp(
		log.NewNopLogger(), db, nil, true, map[int64]bool{}, "/tmp", invCheckPeriod, encCdc, EmptyAppOptions{},
	)
	if withGenesis {
		return testApp, app.NewDefaultGenesisState(encCdc.Codec)
	}
	return testApp, map[string]json.RawMessage{}
}

// Setup initializes a new Rollapp. A Nop logger is set in Rollapp.
func Setup(t *testing.T, isCheckTx bool) *app.App {
	t.Helper()

	pk, err := cryptocodec.ToTmProtoPublicKey(ProposerPK)
	require.NoError(t, err)

	app, genesisState := setup(true, 5)

	// setup for sequencer
	seqGenesis := seqtypes.GenesisState{
		Params:                 seqtypes.DefaultParams(),
		GenesisOperatorAddress: sdk.ValAddress(OperatorPK.Address()).String(),
	}
	genesisState[seqtypes.ModuleName] = app.AppCodec().MustMarshalJSON(&seqGenesis)

	// for now bank genesis won't be set here, funding accounts should be called with fund utils.FundModuleAccount

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)
	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			Time:            time.Time{},
			ChainId:         "test_100-1",
			ConsensusParams: DefaultConsensusParams,
			Validators: []abci.ValidatorUpdate{
				{PubKey: pk, Power: 1},
			},
			AppStateBytes: stateBytes,
			InitialHeight: 0,
		},
	)

	return app
}

// TODO: tech debt - this is almost the same as in github.com/cosmos/ibc-go/v6/testing/app.go
// but unlike the other one, this one adds the sequencer to the genesis state on InitChain
func SetupWithGenesisValSet(t *testing.T, chainID, rollAppDenom string, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances []banktypes.Balance) *app.App {
	t.Helper()

	app, genesisState := setup(true, 5)

	seqGenesis := seqtypes.GenesisState{
		Params:                 seqtypes.DefaultParams(),
		GenesisOperatorAddress: sdk.ValAddress(OperatorPK.Address()).String(),
	}
	genesisState[seqtypes.ModuleName] = app.AppCodec().MustMarshalJSON(&seqGenesis)

	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction)
	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.ZeroInt(),
		}

		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))
	}

	// set validators and delegations
	var stakingGenesis stakingtypes.GenesisState
	app.AppCodec().MustUnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesis)

	bondDenom := stakingGenesis.Params.BondDenom

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(bondDenom, bondAmt.Mul(sdk.NewInt(int64(len(valSet.Validators)))))},
	})

	genModuleAmount, ok := sdk.NewIntFromString("100000000000000000000")
	require.True(t, ok)

	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(types.ModuleName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(rollAppDenom, genModuleAmount)},
	})

	pk, err := cryptocodec.ToTmProtoPublicKey(ProposerPK)
	require.NoError(t, err)

	// set validators and delegations
	stakingGenesis = *stakingtypes.NewGenesisState(stakingGenesis.Params, validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(&stakingGenesis)

	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, sdk.NewCoins(), []banktypes.Metadata{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			Time:            time.Time{},
			ChainId:         chainID,
			ConsensusParams: DefaultConsensusParams,
			Validators: []abci.ValidatorUpdate{
				{PubKey: pk, Power: 1},
			},
			AppStateBytes: stateBytes,
			InitialHeight: 0,
		},
	)

	return app
}
