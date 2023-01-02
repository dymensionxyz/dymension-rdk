<div align="center">
  <h1> Dymension Rollapp </h1>
</div>

![banner](https://user-images.githubusercontent.com/109034310/204804891-bdc0f7bc-4b17-4b4a-99ff-25153d3887ee.jpg)


<!-- <style>
img[src*="#thumbnail"] {
     display: block;
  margin-left: auto;
  margin-right: auto;
}
</style>  -->

[![license](https://img.shields.io/github/license/cosmos/cosmos-sdk.svg#thumbnail)](https://github.com/dymensionxyz/rdk/blob/main/LICENSE)


Dymension RDK, which stands for *RollApp Development Kit* is based on the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) framework, modified and enhanced for building autonomous RollApps (app-specfic-rollups) on top of the [Dymension Hub](https://github.com/dymensionxyz/dymension). 

The RDK provides the following capabilites ***on top*** of the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) framework: 
* The RDK is coupled with the [Dymint](https://github.com/dymensionxyz/dymint) client to form RollApp's blazing speed consensus and networking layer, while ***the Dymension Hub is securing the rollapp***
* Custom modules that converts a cosmos based PoS (proof-of-stake) chain to a rollapp
* wasm and EVM support (based on CosmWasm and Ethermint)

### Learn more
For more information about Dymension RollApps please visit the [documentation center](https://docs.dymension.xyz/)

To learn how the Cosmos SDK works from a high-level perspective, see the Cosmos SDK [High-Level Intro](https://docs.cosmos.network/main/intro/overview.html).

If you want to get started quickly and learn how to build on top of Cosmos SDK, visit [Cosmos SDK Tutorials](https://tutorials.cosmos.network). You can also fork the tutorial's repository to get started building your own Cosmos SDK application.

---

This repository hosts `rollappd`, the first implementation of a dymension rollapp.


**Note**: Requires [Go 1.18](https://go.dev/)

**Note**: This code was intiially scaffolded with igniteCLI@v0.22.2

For critical security issues & disclosure, see [SECURITY.md](SECURITY.md).


## Quick guide
Get started with [building RollApps](https://docs.dymension.xyz/developers/getting-started/intro) 

## Installing / Getting started
```
make install
```
or
```shell
cd cmd/rollappd/
go install
```

This will build the ```rollappd``` binary


### Initial configuration
set custom configuration params at `scripts/shared.sh`


>sh scripts/init_rollapp.sh

This will initilize the rollapp with single initial staked account

### Register rollapp on settlement

>sh scripts/register_rollapp.sh

validate using 
>dymd q rollapp list-rollapp


### Register sequencer for rollapp on settlement

>sh scripts/register_sequencer.sh

validate using 
>dymd q sequencer list-sequencer


### Run rollapp

>sh scripts/run_rollapp.sh


## Developers guide
TODO








<!-- 
# Future features (WIP)

## Fully support ignite 
### using ignite as developing framework
Rollapp customization should allow the usage of `ignite` for scaffolding custom modules
### using ignite to run rollapp
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
- [Developer Chat](https://discord.gg/ignite) -->
