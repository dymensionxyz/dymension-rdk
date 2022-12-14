# rollapp
**rollapp** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

The Rollapp chain is already wired with CosmWasm module and using Dymension settlement.
Scaffolded with igniteCLI@v0.22.2


## Get started
### Build
either using 
`go install` from `cmd/rollappd/`
or
`ignite chain build`


### Init
set configuration params at `scripts/shared.sh`

run
`sh scripts/init_rollapp.sh`


This will initilize the rollapp with single initial staked account

### Register rollapp on settlement
run
`sh scripts/register_rollapp.sh`

validate using 
`dymd q rollapp list-rollapp`


### Register sequencer for rollapp on settlement
run
`sh scripts/register_sequencer.sh`

validate using 
`dymd q sequencer list-sequencer`


### Run rollapp
run
`sh scripts/run_rollapp.sh`



## Using ignite (WIP. not tested)
```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Ignite CLI has scaffolded a Vue.js-based web app in the `vue` directory. Run the following commands to install dependencies and start the app:

```
cd vue
npm install
npm run serve
```

The frontend app is built using the `@starport/vue` and `@starport/vuex` packages. For details, see the [monorepo for Ignite front-end development](https://github.com/ignite/web).

## Release
To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/dymensionxyz/rollapp@latest! | sudo bash
```
`dymensionxyz/rollapp` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/allinbits/starport-installer).

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)
