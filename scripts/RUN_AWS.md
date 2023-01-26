
# Run env on AWS
## Exports
```
export SERVER1=ec2-18-202-218-150.eu-west-1.compute.amazonaws.com
export SERVER2=ec2-34-255-11-227.eu-west-1.compute.amazonaws.com
export SERVER3=ec2-34-249-17-158.eu-west-1.compute.amazonaws.com
```

## Prepare servers
```
sudo yum update
sudo yum install git go jq make
mkdir dymension && cd dymension
export PATH=$PATH:$HOME/go/bin
```

## Server 1 (hub validator 1)
```
git clone https://github.com/dymensionxyz/dymension.git
cd dymension
make install

sh scripts/setup_local.sh
```

copy genesis from server 1 to a local directory
```
scp ec2-user@$SERVER1:/home/ec2-user/.dymension/config/genesis.json .
```

write the node-id of server1
```
validator1_node_id=dymd tendermint show-node-id
```

## Server 2 (hub validator 2)
```
git clone https://github.com/dymensionxyz/dymension.git
cd dymension
make install

sh scripts/setup_local.sh
```

copy genesis from local directory to server2
```
scp genesis.json ec2-user@$SERVER2:/home/ec2-user/.dymension/config/genesis.json
```

Edit peers into config.toml on server 2:
```
validator1_node_id@$SERVER1:36656
```

check the address of this account 
```
dymd keys show -a local-user --keyring-backend test
```
and fund it on server1:
```
    dymd tx bank send $(dymd keys show -a local-user --keyring-backend test) XXXXX 10000000000udym --keyring-backend test
```

Create a validator on node 2 as well
```
dymd tx staking create-validator \
  --amount 1000000udym \
  --commission-max-change-rate "0.1" \
  --commission-max-rate "0.20" \
  --commission-rate "0.1" \
  --min-self-delegation "1" \
  --details "validators write bios too" \
  --pubkey=$(dymd tendermint show-validator) \
  --moniker "2ndmoniker" \
  --chain-id "local-testnet" \
  --gas-prices 0.025udym \
  --from local-user \
  --keyring-backend test
```



## Server3 (rollapp sequencer)
```
git clone https://github.com/dymensionxyz/dymension.git
cd dymension
make install
cd ..

git clone https://github.com/dymensionxyz/relayer.git
cd relayer
make install
cd ..

git clone https://github.com/dymensionxyz/dymension-rdk.git
cd dymension-rdk
go install cmd/rollappd 
```

export the following:
```
export SERVER1=ec2-18-202-218-150.eu-west-1.compute.amazonaws.com
export SETTLEMENT_RPC="$SERVER1:36657"

export SETTLEMENT_RPC_FOR_RELAYER=$SETTLEMENT_RPC
export ROLLAPP_RPC_FOR_RELAYER="127.0.0.1:26667"
```

Run the rollapp installation readme


To run the relayer:
```
sh scripts/setup_ibc.sh
sh scripts/run_relayer.sh
```