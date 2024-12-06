package ante

import (
	"fmt"
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/group"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func TestBypassIBCFeeDecorator(t *testing.T) {
	whitelistedSigner := sdk.AccAddress("whitelisted")
	nonWhitelistedSigner := sdk.AccAddress("nonwhitelisted")
	operatorAddr := "cosmosvaloper1tnh2q55v8wyygtt9srz5safamzdengsn9dsd7z"
	consAddr := sdk.ConsAddress("consAddr")

	// IBC relayer msg
	ibcMsg1 := &clienttypes.MsgCreateClient{}
	// IBC relayer msg
	ibcMsg2 := &channeltypes.MsgAcknowledgement{}
	// Non-IBC msg
	nonIBCMsg := &banktypes.MsgSend{}

	testCases := []struct {
		name             string
		msgs             []sdk.Msg
		signer           sdk.AccAddress
		sequencerExists  bool
		sequencerOper    string
		wlRelayers       []string
		wlError          error
		expectedErr      bool
		expectedIBCNoFee bool
	}{
		{
			name:             "Non-IBC message, no error",
			msgs:             []sdk.Msg{nonIBCMsg},
			signer:           nonWhitelistedSigner,
			sequencerExists:  true,
			sequencerOper:    operatorAddr,
			wlRelayers:       []string{whitelistedSigner.String()},
			expectedErr:      false,
			expectedIBCNoFee: false,
		},
		{
			name:             "All IBC messages, signer whitelisted",
			msgs:             []sdk.Msg{ibcMsg1, ibcMsg2},
			signer:           whitelistedSigner,
			sequencerExists:  true,
			sequencerOper:    operatorAddr,
			wlRelayers:       []string{whitelistedSigner.String()},
			expectedErr:      false,
			expectedIBCNoFee: true,
		},
		{
			name:             "All IBC messages, sequencer not found",
			msgs:             []sdk.Msg{ibcMsg1},
			signer:           whitelistedSigner,
			sequencerExists:  false,
			sequencerOper:    "",
			wlRelayers:       nil,
			expectedErr:      true,
			expectedIBCNoFee: false,
		},
		{
			name:             "All IBC messages, GetWhitelistedRelayers returns error",
			msgs:             []sdk.Msg{ibcMsg1},
			signer:           whitelistedSigner,
			sequencerExists:  true,
			sequencerOper:    operatorAddr,
			wlRelayers:       nil,
			wlError:          fmt.Errorf("some error"),
			expectedErr:      true,
			expectedIBCNoFee: false,
		},
		{
			name:             "All IBC messages, signer not in whitelist",
			msgs:             []sdk.Msg{ibcMsg1},
			signer:           nonWhitelistedSigner,
			sequencerExists:  true,
			sequencerOper:    operatorAddr,
			wlRelayers:       []string{whitelistedSigner.String()},
			expectedErr:      true,
			expectedIBCNoFee: false,
		},
		{
			name:             "Mixed messages (IBC and non-IBC), not allowed",
			msgs:             []sdk.Msg{ibcMsg1, nonIBCMsg},
			signer:           whitelistedSigner,
			sequencerExists:  true,
			sequencerOper:    operatorAddr,
			wlRelayers:       []string{whitelistedSigner.String()},
			expectedErr:      true,
			expectedIBCNoFee: false,
		},
		{
			name:             "Mixed messages (IBC and non-IBC), signer not in whitelist",
			msgs:             []sdk.Msg{ibcMsg1, nonIBCMsg},
			signer:           nonWhitelistedSigner,
			sequencerExists:  true,
			sequencerOper:    operatorAddr,
			wlRelayers:       []string{whitelistedSigner.String()},
			expectedErr:      true,
			expectedIBCNoFee: false,
		},
		{
			name: "Nested scenario: multi-level. MsgExec containing a MsgSubmitProposal(gov) that returns [ibcMsg2]",
			msgs: []sdk.Msg{
				&authz.MsgExec{
					Grantee: whitelistedSigner.String(),
					Msgs: []*cdctypes.Any{
						func() *cdctypes.Any {
							msg, _ := cdctypes.NewAnyWithValue(&govtypesv1.MsgSubmitProposal{
								Messages: []*cdctypes.Any{
									func() *cdctypes.Any {
										ibcMsg2V := *ibcMsg2
										ibcMsg2V.Signer = whitelistedSigner.String()
										msg, _ := cdctypes.NewAnyWithValue(&ibcMsg2V)
										return msg
									}(),
								},
								InitialDeposit: sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(100))),
								Proposer:       whitelistedSigner.String(),
								Metadata:       "==",
							})
							return msg
						}(),
					},
				},
			},
			signer:           whitelistedSigner,
			sequencerExists:  true,
			sequencerOper:    operatorAddr,
			wlRelayers:       []string{whitelistedSigner.String()},
			expectedErr:      false,
			expectedIBCNoFee: true, // all final msgs are IBC and whitelisted
		},
		{
			name: "Nested scenario: multi-level. MsgExec containing MsgSubmitProposal(gov) with ibcMsg2 but signer not whitelisted",
			msgs: []sdk.Msg{
				&authz.MsgExec{
					Grantee: nonWhitelistedSigner.String(),
					Msgs: []*cdctypes.Any{
						func() *cdctypes.Any {
							msg, _ := cdctypes.NewAnyWithValue(&govtypesv1.MsgSubmitProposal{
								Messages: []*cdctypes.Any{
									func() *cdctypes.Any {
										ibcMsg2V := *ibcMsg2
										ibcMsg2V.Signer = nonWhitelistedSigner.String()
										msg, _ := cdctypes.NewAnyWithValue(&ibcMsg2V)
										return msg
									}(),
								},
								InitialDeposit: sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(100))),
								Proposer:       whitelistedSigner.String(),
								Metadata:       "==",
							})
							return msg
						}(),
					},
				},
			},
			signer:           whitelistedSigner,
			sequencerExists:  true,
			sequencerOper:    operatorAddr,
			wlRelayers:       []string{whitelistedSigner.String()},
			expectedErr:      true,
			expectedIBCNoFee: false, // signer not whitelisted
		},
		{
			name: "Nested scenario: exceed maxDepth",
			msgs: []sdk.Msg{
				wrapMsgInSubmitProposal(
					wrapMsgInSubmitProposal(
						wrapMsgInSubmitProposal(
							wrapMsgInSubmitProposal(
								wrapMsgInSubmitProposal(
									wrapMsgInSubmitProposal(ibcMsg1)))))),
			},
			signer:           whitelistedSigner,
			sequencerExists:  true,
			sequencerOper:    operatorAddr,
			wlRelayers:       []string{whitelistedSigner.String()},
			expectedErr:      true, // exceeds maxDepth
			expectedIBCNoFee: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := sdk.Context{}.WithLogger(log.NewNopLogger())
			dk := &mockDistrKeeper{consAddr: consAddr}
			sk := &mockSequencerKeeper{
				sequencerFound:      tc.sequencerExists,
				operatorAddr:        tc.sequencerOper,
				getWhitelistedError: tc.wlError,
				whitelistedRelayers: tc.wlRelayers,
			}
			nextAnte := &mockNextAnte{}

			nextCalled := false
			next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
				nextCalled = true
				return ctx, nil
			}

			for _, msg := range tc.msgs {
				switch m := msg.(type) {
				case *clienttypes.MsgCreateClient:
					m.Signer = tc.signer.String()
				case *channeltypes.MsgAcknowledgement:
					m.Signer = tc.signer.String()
				case *banktypes.MsgSend:
					m.FromAddress = tc.signer.String()
				}
			}
			tx := &mockTx{
				msgs: tc.msgs,
			}

			decor := BypassIBCFeeDecorator{
				dk:       dk,
				sk:       sk,
				nextAnte: nextAnte,
			}

			_, err := decor.AnteHandle(ctx, tx, false, next)
			if tc.expectedErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else {
					if tc.expectedIBCNoFee {
						// we expect next to be called directly (no fee)
						if !nextCalled {
							t.Errorf("expected next handler to be called, but it wasn't")
						}
					} else {
						// if no ibc no fee scenario, we fall back to normal fee ante
						// means we call n.nextAnte.AnteHandle
						if !nextAnte.called {
							t.Errorf("expected fallback to normal fee ante")
						}
					}
				}
			}
		})
	}
}

type mockSequencerKeeper struct {
	sequencerFound      bool
	operatorAddr        string
	getWhitelistedError error
	whitelistedRelayers []string
}

func (sk *mockSequencerKeeper) GetSequencerByConsAddr(sdk.Context, sdk.ConsAddress) (stakingtypes.Validator, bool) {
	if !sk.sequencerFound {
		return stakingtypes.Validator{}, false
	}
	return stakingtypes.Validator{
		OperatorAddress: sk.operatorAddr,
	}, true
}

func (sk *mockSequencerKeeper) GetWhitelistedRelayers(sdk.Context, sdk.ValAddress) (types.WhitelistedRelayers, error) {
	if sk.getWhitelistedError != nil {
		return types.WhitelistedRelayers{}, sk.getWhitelistedError
	}
	return types.WhitelistedRelayers{Relayers: sk.whitelistedRelayers}, nil
}

type mockDistrKeeper struct {
	consAddr sdk.ConsAddress
}

func (m mockDistrKeeper) GetPreviousProposerConsAddr(sdk.Context) sdk.ConsAddress {
	return m.consAddr
}

type mockNextAnte struct {
	called bool
}

func (m *mockNextAnte) AnteHandle(ctx sdk.Context, _ sdk.Tx, _ bool, _ sdk.AnteHandler) (sdk.Context, error) {
	m.called = true
	return ctx, nil
}

type mockTx struct {
	msgs []sdk.Msg
}

func (m mockTx) GetMsgs() []sdk.Msg   { return m.msgs }
func (m mockTx) ValidateBasic() error { return nil }

func wrapMsgInSubmitProposal(inMsg sdk.Msg) sdk.Msg {
	return &group.MsgSubmitProposal{
		Messages: []*cdctypes.Any{
			func() *cdctypes.Any {
				msg, _ := cdctypes.NewAnyWithValue(inMsg)
				return msg
			}(),
		},
	}
}
