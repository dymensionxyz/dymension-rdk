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

	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func TestBypassIBCFeeDecorator(t *testing.T) {
	whitelistedSigner := sdk.AccAddress("whitelisted")
	nonWhitelistedSigner := sdk.AccAddress("nonwhitelisted")
	operatorAddr := "cosmosvaloper1tnh2q55v8wyygtt9srz5safamzdengsn9dsd7z"
	consAddr := sdk.ConsAddress("consAddr")

	lifecycle0 := &clienttypes.MsgCreateClient{}
	lifecycle1 := &channeltypes.MsgChannelOpenInit{}
	normal0 := &channeltypes.MsgRecvPacket{}
	normal1 := &clienttypes.MsgUpdateClient{}
	nonIBCMsg := &banktypes.MsgSend{}

	testCases := []struct {
		name            string
		msgs            []sdk.Msg
		signer          sdk.AccAddress
		sequencerExists bool
		sequencerOper   string
		wl              []string
		wlError         error
		expectErr       bool
		expectNoFee     bool
		freeIBC         bool
	}{
		{
			name:            "Non-IBC message, no error",
			msgs:            []sdk.Msg{nonIBCMsg},
			signer:          nonWhitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     false,
		},
		{
			name:            "All IBC messages, signer whitelisted",
			msgs:            []sdk.Msg{lifecycle0, lifecycle1},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "All IBC messages, sequencer not found",
			msgs:            []sdk.Msg{lifecycle0},
			signer:          whitelistedSigner,
			sequencerExists: false,
			sequencerOper:   "",
			wl:              nil,
			expectErr:       true,
			expectNoFee:     false,
		},
		{
			name:            "Packet IBC message, signer whitelisted",
			msgs:            []sdk.Msg{lifecycle0, normal0},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "Packet IBC message, signer not whitelisted",
			msgs:            []sdk.Msg{normal0, normal1},
			signer:          nonWhitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     false,
		},
		{
			name:            "Whitelisted IBC and Packet IBC message, signer not whitelisted",
			msgs:            []sdk.Msg{lifecycle0, normal0},
			signer:          nonWhitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       true,
			expectNoFee:     false,
		},
		{
			name:            "All IBC messages, GetWhitelistedRelayers returns error",
			msgs:            []sdk.Msg{lifecycle0},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              nil,
			wlError:         fmt.Errorf("some error"),
			expectErr:       true,
			expectNoFee:     false,
		},
		{
			name:            "All IBC messages, signer not in whitelist",
			msgs:            []sdk.Msg{lifecycle0},
			signer:          nonWhitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       true,
			expectNoFee:     false,
		},
		{
			name:            "Mixed messages (IBC and non-IBC), not allowed",
			msgs:            []sdk.Msg{lifecycle0, nonIBCMsg},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       true,
			expectNoFee:     false,
		},
		{
			name:            "Mixed messages (IBC and non-IBC), signer not in whitelist",
			msgs:            []sdk.Msg{lifecycle0, nonIBCMsg},
			signer:          nonWhitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       true,
			expectNoFee:     false,
		},
		{
			name: "Nested scenario: multi-level. MsgExec containing a MsgSubmitProposal(group) that returns [ibcMsg2]",
			msgs: []sdk.Msg{
				&authz.MsgExec{
					Grantee: whitelistedSigner.String(),
					Msgs: []*cdctypes.Any{
						func() *cdctypes.Any {
							msg, _ := cdctypes.NewAnyWithValue(&group.MsgSubmitProposal{
								Proposers: []string{whitelistedSigner.String()},
								Metadata:  "==",
								Messages: []*cdctypes.Any{
									func() *cdctypes.Any {
										ibcMsg2V := *lifecycle1
										ibcMsg2V.Signer = whitelistedSigner.String()
										msg, _ := cdctypes.NewAnyWithValue(&ibcMsg2V)
										return msg
									}(),
								},
								Exec: 0,
							})
							return msg
						}(),
					},
				},
			},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true, // all final msgs are IBC and whitelisted
		},
		{
			name: "Nested scenario: not checked. MsgExec containing MsgSubmitProposal(gov) with ibcMsg2 but signer not whitelisted",
			msgs: []sdk.Msg{
				&authz.MsgExec{
					Grantee: nonWhitelistedSigner.String(),
					Msgs: []*cdctypes.Any{
						func() *cdctypes.Any {
							msg, _ := cdctypes.NewAnyWithValue(&govtypesv1.MsgSubmitProposal{
								Messages: []*cdctypes.Any{
									func() *cdctypes.Any {
										ibcMsg2V := *lifecycle1
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
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     false, // signer not whitelisted
		},
		{
			name: "Nested scenario: multi-level. MsgExec containing MsgSubmitProposal(group) with ibcMsg2 but signer not whitelisted",
			msgs: []sdk.Msg{
				&authz.MsgExec{
					Grantee: nonWhitelistedSigner.String(),
					Msgs: []*cdctypes.Any{
						func() *cdctypes.Any {
							msg, _ := cdctypes.NewAnyWithValue(&group.MsgSubmitProposal{
								Proposers: []string{whitelistedSigner.String()},
								Messages: []*cdctypes.Any{
									func() *cdctypes.Any {
										ibcMsg2V := *lifecycle1
										ibcMsg2V.Signer = nonWhitelistedSigner.String()
										msg, _ := cdctypes.NewAnyWithValue(&ibcMsg2V)
										return msg
									}(),
								},
								Metadata: "==",
							})
							return msg
						}(),
					},
				},
			},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       true,
			expectNoFee:     false, // signer not whitelisted
		},
		{
			name: "Nested scenario: exceed maxDepth",
			msgs: []sdk.Msg{
				wrapMsgInSubmitProposal(
					wrapMsgInSubmitProposal(
						wrapMsgInSubmitProposal(
							wrapMsgInSubmitProposal(
								wrapMsgInSubmitProposal(
									wrapMsgInSubmitProposal(lifecycle0)))))),
			},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       true, // exceeds maxDepth
			expectNoFee:     false,
		},
		{
			name:            "Free does not allow lifecycle",
			msgs:            []sdk.Msg{lifecycle0},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{},
			expectErr:       true, // whitelist still checked
			expectNoFee:     false,
			freeIBC:         true,
		},
		{
			name:            "No charge if not lifecycle msg and free",
			msgs:            []sdk.Msg{normal0},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{},
			expectErr:       false,
			expectNoFee:     true,
			freeIBC:         true,
		},
		{
			name:            "MsgGrant should bypass fees unconditionally",
			msgs:            []sdk.Msg{&authz.MsgGrant{}},
			signer:          nonWhitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "MsgGrant with non-whitelisted signer should still bypass fees",
			msgs:            []sdk.Msg{&authz.MsgGrant{}},
			signer:          nonWhitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "MsgGrant should bypass fees even when sequencer not found",
			msgs:            []sdk.Msg{&authz.MsgGrant{}},
			signer:          whitelistedSigner,
			sequencerExists: false,
			sequencerOper:   "",
			wl:              nil,
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "MsgGrant should bypass fees even when GetWhitelistedRelayers returns error",
			msgs:            []sdk.Msg{&authz.MsgGrant{}},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              nil,
			wlError:         fmt.Errorf("some error"),
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "Multiple MsgGrant messages should all bypass fees",
			msgs:            []sdk.Msg{&authz.MsgGrant{}, &authz.MsgGrant{}},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "MsgGrant mixed with IBC messages should fail",
			msgs:            []sdk.Msg{&authz.MsgGrant{}, lifecycle0},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       true,
			expectNoFee:     false,
		},
		{
			name:            "MsgGrant mixed with non-IBC messages should fail",
			msgs:            []sdk.Msg{&authz.MsgGrant{}, nonIBCMsg},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       true,
			expectNoFee:     false,
		},
		{
			name: "MsgGrant in nested MsgExec should bypass fees",
			msgs: []sdk.Msg{
				&authz.MsgExec{
					Grantee: whitelistedSigner.String(),
					Msgs: []*cdctypes.Any{
						func() *cdctypes.Any {
							msg, _ := cdctypes.NewAnyWithValue(&authz.MsgGrant{})
							return msg
						}(),
					},
				},
			},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name: "MsgGrant in nested MsgSubmitProposal should bypass fees",
			msgs: []sdk.Msg{
				&group.MsgSubmitProposal{
					Proposers: []string{whitelistedSigner.String()},
					Metadata:  "==",
					Messages: []*cdctypes.Any{
						func() *cdctypes.Any {
							msg, _ := cdctypes.NewAnyWithValue(&authz.MsgGrant{})
							return msg
						}(),
					},
					Exec: 0,
				},
			},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name: "MsgGrant in deeply nested structure should bypass fees",
			msgs: []sdk.Msg{
				wrapMsgInSubmitProposal(
					wrapMsgInSubmitProposal(
						&authz.MsgGrant{})),
			},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "MsgGrant should bypass fees regardless of freeIBC setting",
			msgs:            []sdk.Msg{&authz.MsgGrant{}},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
			freeIBC:         false,
		},
		{
			name:            "MsgGrant should bypass fees regardless of freeIBC setting (true)",
			msgs:            []sdk.Msg{&authz.MsgGrant{}},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
			freeIBC:         true,
		},
		{
			name:            "MsgGrantAllowance should bypass fees unconditionally",
			msgs:            []sdk.Msg{&feegrant.MsgGrantAllowance{}},
			signer:          nonWhitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "MsgGrantAllowance with non-whitelisted signer should still bypass fees",
			msgs:            []sdk.Msg{&feegrant.MsgGrantAllowance{}},
			signer:          nonWhitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "MsgGrantAllowance should bypass fees even when sequencer not found",
			msgs:            []sdk.Msg{&feegrant.MsgGrantAllowance{}},
			signer:          whitelistedSigner,
			sequencerExists: false,
			sequencerOper:   "",
			wl:              nil,
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "MsgGrantAllowance should bypass fees even when GetWhitelistedRelayers returns error",
			msgs:            []sdk.Msg{&feegrant.MsgGrantAllowance{}},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              nil,
			wlError:         fmt.Errorf("some error"),
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "Multiple MsgGrantAllowance messages should all bypass fees",
			msgs:            []sdk.Msg{&feegrant.MsgGrantAllowance{}, &feegrant.MsgGrantAllowance{}},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "MsgGrantAllowance mixed with IBC messages should fail",
			msgs:            []sdk.Msg{&feegrant.MsgGrantAllowance{}, lifecycle0},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       true,
			expectNoFee:     false,
		},
		{
			name:            "MsgGrantAllowance mixed with non-IBC messages should fail",
			msgs:            []sdk.Msg{&feegrant.MsgGrantAllowance{}, nonIBCMsg},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       true,
			expectNoFee:     false,
		},
		{
			name: "MsgGrantAllowance in nested MsgExec should bypass fees",
			msgs: []sdk.Msg{
				&authz.MsgExec{
					Grantee: whitelistedSigner.String(),
					Msgs: []*cdctypes.Any{
						func() *cdctypes.Any {
							msg, _ := cdctypes.NewAnyWithValue(&feegrant.MsgGrantAllowance{})
							return msg
						}(),
					},
				},
			},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name: "MsgGrantAllowance in nested MsgSubmitProposal should bypass fees",
			msgs: []sdk.Msg{
				&group.MsgSubmitProposal{
					Proposers: []string{whitelistedSigner.String()},
					Metadata:  "==",
					Messages: []*cdctypes.Any{
						func() *cdctypes.Any {
							msg, _ := cdctypes.NewAnyWithValue(&feegrant.MsgGrantAllowance{})
							return msg
						}(),
					},
					Exec: 0,
				},
			},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
		},
		{
			name:            "MsgGrantAllowance should bypass fees regardless of freeIBC setting",
			msgs:            []sdk.Msg{&feegrant.MsgGrantAllowance{}},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
			freeIBC:         false,
		},
		{
			name:            "Mixed free non-IBC messages (MsgGrant + MsgGrantAllowance) should bypass fees",
			msgs:            []sdk.Msg{&authz.MsgGrant{}, &feegrant.MsgGrantAllowance{}},
			signer:          whitelistedSigner,
			sequencerExists: true,
			sequencerOper:   operatorAddr,
			wl:              []string{whitelistedSigner.String()},
			expectErr:       false,
			expectNoFee:     true,
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
				whitelistedRelayers: tc.wl,
			}
			pk := &mockParamsK{tc.freeIBC}
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
				case *channeltypes.MsgChannelOpenInit:
					m.Signer = tc.signer.String()
				case *channeltypes.MsgRecvPacket:
					m.Signer = tc.signer.String()
				case *banktypes.MsgSend:
					m.FromAddress = tc.signer.String()
				}
			}
			tx := &mockTx{
				msgs: tc.msgs,
			}

			decor := BypassIBCFeeDecorator{
				dk:           dk,
				sk:           sk,
				pk:           pk,
				bypassedAnte: nextAnte,
			}

			_, err := decor.AnteHandle(ctx, tx, false, next)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else {
					if tc.expectNoFee {
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

type mockParamsK struct {
	ret bool
}

func (m mockParamsK) FreeIBC(ctx sdk.Context) bool {
	return m.ret
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
