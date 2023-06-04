#! /bin/bash

tmp=$(mktemp)


set_EVM_params() {
  jq '.consensus_params["block"]["max_gas"] = "40000000"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq '.app_state["feemarket"]["params"]["no_base_fee"] = true' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
}

set_denom() {
  denom=$1
  jq --arg denom $denom '.app_state.mint.params.mint_denom = $denom' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq --arg denom $denom '.app_state.staking.params.bond_denom = $denom' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq --arg denom $denom '.app_state.crisis.constant_fee.denom = $denom' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq --arg denom $denom '.app_state.evm.params.evm_denom = $denom' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq --arg denom $denom '.app_state.gov.deposit_params.min_deposit[0].denom = $denom' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
}

set_minting_params() {
  echo "setting minting params"
  jq '.app_state.mint.params.genesis_epoch_provisions = "1000000"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  # jq '(.app_state.epochs.epochs[] | select(.identifier=="mint") .duration)="60s"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
}


set_distribution_params() {
  echo "setting distribution params"
  jq '.app_state.distribution.params.base_proposer_reward = "0.8"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq '.app_state.distribution.params.community_tax = "0.00002"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
}


set_gov_params() {
  echo "setting gov params"
  # jq '.app_state.gov.tally_params.quorum = ""' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  # jq '.app_state.gov.tally_params.threshold = ""' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  # jq '.app_state.gov.tally_params.veto_threshold = ""' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq '.app_state.gov.voting_params.voting_period = "300s"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
}

set_staking_params() {
  echo "setting staking params"
  jq '.app_state.staking.params.unbonding_time = "3628800s"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
}
