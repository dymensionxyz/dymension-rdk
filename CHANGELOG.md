# [](https://github.com/dymensionxyz/dymension-rdk/compare/v1.5.0-beta...v) (2024-05-24)


### Features

* Add  collect-gentxs command to pickup custom gendoc ([#388](https://github.com/dymensionxyz/dymension-rdk/issues/388)) ([f6997e0](https://github.com/dymensionxyz/dymension-rdk/commit/f6997e0d1a5c38ac5637f6e615845bfb78842976))
* **denommetadata:** Add support for providing denom-trace and creation of multiple denom-metadatas  ([#384](https://github.com/dymensionxyz/dymension-rdk/issues/384)) ([c720640](https://github.com/dymensionxyz/dymension-rdk/commit/c720640d7d392d58fb191852cd8f4a5d85afc04e))



# [1.5.0-beta](https://github.com/dymensionxyz/dymension-rdk/compare/v1.4.0-beta...v1.5.0-beta) (2024-04-14)



# [1.4.0-beta](https://github.com/dymensionxyz/dymension-rdk/compare/v1.3.0-beta...v1.4.0-beta) (2024-04-13)


### Bug Fixes

* **denommetadata:** fixed ibc voucher denom query ([#378](https://github.com/dymensionxyz/dymension-rdk/issues/378)) ([bf14e28](https://github.com/dymensionxyz/dymension-rdk/commit/bf14e28256c444930195b3a6df760d3626c205fa))



# [1.3.0-beta](https://github.com/dymensionxyz/dymension-rdk/compare/v1.2.0-beta...v1.3.0-beta) (2024-04-11)


### Bug Fixes

* **allocation:** do not destroy fee rewards if allocating to proposer fails ([#358](https://github.com/dymensionxyz/dymension-rdk/issues/358)) ([bfd0943](https://github.com/dymensionxyz/dymension-rdk/commit/bfd09432d61eb2c33e29621a1ab868661aeda50c))
* **mint:** added missing cli query for current inflation ([#366](https://github.com/dymensionxyz/dymension-rdk/issues/366)) ([6ecadd8](https://github.com/dymensionxyz/dymension-rdk/commit/6ecadd818429e1870da7e6f6e044bc87f900ec8b))
* **sequencer:** added HistoricalEntries param validation ([#369](https://github.com/dymensionxyz/dymension-rdk/issues/369)) ([8e4b891](https://github.com/dymensionxyz/dymension-rdk/commit/8e4b891e68304912b839c2b3ef2a733d91426ce2))
* **sequencers:** sequencers genesis operator address not exported ([#367](https://github.com/dymensionxyz/dymension-rdk/issues/367)) ([2a1f308](https://github.com/dymensionxyz/dymension-rdk/commit/2a1f308c31141143c20fadee37578c5be0df4f87))


### Features

* add get ibc denom base on denom trace query ([#376](https://github.com/dymensionxyz/dymension-rdk/issues/376)) ([807717b](https://github.com/dymensionxyz/dymension-rdk/commit/807717b3219609a68e14d4f11100c539db001008))



# [1.2.0-beta](https://github.com/dymensionxyz/dymension-rdk/compare/v1.1.0-beta...v1.2.0-beta) (2024-03-26)


### Features

* **deonmetadata:** enable testing for denom metadata creation with erc20 hook ([#343](https://github.com/dymensionxyz/dymension-rdk/issues/343)) ([8e3dc66](https://github.com/dymensionxyz/dymension-rdk/commit/8e3dc66c8c6d6894cbf2d5d12112772069ea8774))



# [1.1.0-beta](https://github.com/dymensionxyz/dymension-rdk/compare/v1.0.0-beta...v1.1.0-beta) (2024-03-22)


### Bug Fixes

* dymint reset-state command not working ([#309](https://github.com/dymensionxyz/dymension-rdk/issues/309)) ([4f5a707](https://github.com/dymensionxyz/dymension-rdk/commit/4f5a70734530f4da1c8980cdeb9b59d5a2148577))


### Features

* adding_inspect_command ([#289](https://github.com/dymensionxyz/dymension-rdk/issues/289)) ([a58be89](https://github.com/dymensionxyz/dymension-rdk/commit/a58be893b8aa113b302aeb308036858b80f7958c))
* mint add query cmd for current inflation ([#328](https://github.com/dymensionxyz/dymension-rdk/issues/328)) ([1c49688](https://github.com/dymensionxyz/dymension-rdk/commit/1c49688dc88de8cb70a358ba80eced1a9062543c))
* permissioned creation of token metadata ([#327](https://github.com/dymensionxyz/dymension-rdk/issues/327)) ([ffe8e06](https://github.com/dymensionxyz/dymension-rdk/commit/ffe8e06500e248b4e45097700c3b296131971084))
* rollapp genesis token locking upon channel creation ([#333](https://github.com/dymensionxyz/dymension-rdk/issues/333)) ([4bbdcda](https://github.com/dymensionxyz/dymension-rdk/commit/4bbdcdadb10c44d8458b48ac1c552fa9fcea3321))



# [1.0.0-beta](https://github.com/dymensionxyz/dymension-rdk/compare/v0.4.0-rc1...v1.0.0-beta) (2023-10-19)


### Bug Fixes

* added keyring-backend flag for gentx_seq ([#239](https://github.com/dymensionxyz/dymension-rdk/issues/239)) ([7bf5d44](https://github.com/dymensionxyz/dymension-rdk/commit/7bf5d44d353a0eb5e6823542c44dd0cc94bed921))
* evm module to use sequencers module instead of staking ([#223](https://github.com/dymensionxyz/dymension-rdk/issues/223)) ([315aa83](https://github.com/dymensionxyz/dymension-rdk/commit/315aa83afd3daf381608d59af088f405ad05b3ad))
* fixed epochs proto API path ([#236](https://github.com/dymensionxyz/dymension-rdk/issues/236)) ([0013f9f](https://github.com/dymensionxyz/dymension-rdk/commit/0013f9f12bb52c0e2a2547c3508e242764083ec7))
* sequencer module returns two objects in initchain ([#237](https://github.com/dymensionxyz/dymension-rdk/issues/237)) ([76c5329](https://github.com/dymensionxyz/dymension-rdk/commit/76c5329f3b3801c5ce3ab61a8dfc1a5a8319bc9c))


### Features

* support stakers on genesis ([#260](https://github.com/dymensionxyz/dymension-rdk/issues/260)) ([4632268](https://github.com/dymensionxyz/dymension-rdk/commit/463226890c589dee8d1b41fda246e44130a11e4f))



# [0.4.0-rc1](https://github.com/dymensionxyz/dymension-rdk/compare/v0.2.0-beta...v0.4.0-rc1) (2023-07-02)


### Bug Fixes

* Bump golang.org/x/net from 0.3.0 to 0.7.0 ([#150](https://github.com/dymensionxyz/dymension-rdk/issues/150)) ([1963021](https://github.com/dymensionxyz/dymension-rdk/commit/1963021c7eced52dafb580057716a32dcc53b327))
* changed Int64 to string when logging ([#220](https://github.com/dymensionxyz/dymension-rdk/issues/220)) ([c6cfd1a](https://github.com/dymensionxyz/dymension-rdk/commit/c6cfd1a5f6502e7d24d036709dae53823a5afbe4))
* changed the ibc_go replace directive to the release/v5.1.x_dymint branch ([#185](https://github.com/dymensionxyz/dymension-rdk/issues/185)) ([bb3ef71](https://github.com/dymensionxyz/dymension-rdk/commit/bb3ef71a4211e5904a8fcce55080c36d1d5f118c))
* ci lint failure  ([#165](https://github.com/dymensionxyz/dymension-rdk/issues/165)) ([1b92970](https://github.com/dymensionxyz/dymension-rdk/commit/1b9297044204c7b16868008942a9be82b3ca0813))
* fixed bug when modifying staking params in init script ([#176](https://github.com/dymensionxyz/dymension-rdk/issues/176)) ([96b35de](https://github.com/dymensionxyz/dymension-rdk/commit/96b35deb02bd0b09d8c787a7e0856e93eef598a3))
* fixed ERC20 support ([#140](https://github.com/dymensionxyz/dymension-rdk/issues/140)) ([f4aae00](https://github.com/dymensionxyz/dymension-rdk/commit/f4aae00b8f32d1502d6f5b7fe1d1c1389a7f8b68))
* fixed overwritten signal handler in `cmd/start` ([#195](https://github.com/dymensionxyz/dymension-rdk/issues/195)) ([72bd37f](https://github.com/dymensionxyz/dymension-rdk/commit/72bd37fe31f6229a7f59087b0b264c194e0a02a8))
* RDK upgrade handler fix from v0.45 to v0.46 ([#192](https://github.com/dymensionxyz/dymension-rdk/issues/192)) ([d632283](https://github.com/dymensionxyz/dymension-rdk/commit/d63228382d1df08706d07ef668cb1a06cf404333))


### Features

* add automatic erc20 conversion ([#154](https://github.com/dymensionxyz/dymension-rdk/issues/154)) ([d719490](https://github.com/dymensionxyz/dymension-rdk/commit/d71949027affe46ef5d4bec520fe33a55ab4558d))
* add readme ([#199](https://github.com/dymensionxyz/dymension-rdk/issues/199)) ([935d9ee](https://github.com/dymensionxyz/dymension-rdk/commit/935d9eec25f194cbe086261ac5f7aa77d39ca7c6))
* change ibc script to support dymension relayer ([#156](https://github.com/dymensionxyz/dymension-rdk/issues/156)) ([01c204c](https://github.com/dymensionxyz/dymension-rdk/commit/01c204c87fd0c28edbe577bd0255e231eb28333c))
* change minting to be epoch based ([#198](https://github.com/dymensionxyz/dymension-rdk/issues/198)) ([404bb2b](https://github.com/dymensionxyz/dymension-rdk/commit/404bb2b468b34c607757324e23fa52572ad879f6))
* changed log rotate max size to be configurable vs hardcoded ([#107](https://github.com/dymensionxyz/dymension-rdk/issues/107)) ([1c40dfd](https://github.com/dymensionxyz/dymension-rdk/commit/1c40dfd5ba464ac3df3cea291fb6a1bbb242e765))
* clean up scripts folder ([#149](https://github.com/dymensionxyz/dymension-rdk/issues/149)) ([de8b201](https://github.com/dymensionxyz/dymension-rdk/commit/de8b201a6f5f157207d5518cc94608df123a40f2))
* Enhanced log config and enable log level override on a module basis ([#109](https://github.com/dymensionxyz/dymension-rdk/issues/109)) ([7292826](https://github.com/dymensionxyz/dymension-rdk/commit/729282680fe93335ff7c40d5e5c6586936222aaa))
* move dymint related commands to dymint repo ([#206](https://github.com/dymensionxyz/dymension-rdk/issues/206)) ([961dc98](https://github.com/dymensionxyz/dymension-rdk/commit/961dc9886589c58dea880125da2574021cab2fb0))
* passing max_log_size to the logger ([#111](https://github.com/dymensionxyz/dymension-rdk/issues/111)) ([ef2d364](https://github.com/dymensionxyz/dymension-rdk/commit/ef2d36476034ced14f2460719f0c90dca523f791))
* simple way to add sequencer on genesis ([#183](https://github.com/dymensionxyz/dymension-rdk/issues/183)) ([53f56f9](https://github.com/dymensionxyz/dymension-rdk/commit/53f56f9d9aef67c613b94b88f78ff107f900f5d5))
* support getting active sequencers from dymint ([#147](https://github.com/dymensionxyz/dymension-rdk/issues/147)) ([b3c8f82](https://github.com/dymensionxyz/dymension-rdk/commit/b3c8f8262c48670eb8cc48e743b09d2804f70ffd))
* support reading dymint configuration from file ([#187](https://github.com/dymensionxyz/dymension-rdk/issues/187)) ([5feb350](https://github.com/dymensionxyz/dymension-rdk/commit/5feb3504c4c4820222b5c6ce18a84d029c0f0154))
* update rdk to use `cometbft` ([#209](https://github.com/dymensionxyz/dymension-rdk/issues/209)) ([f4ae301](https://github.com/dymensionxyz/dymension-rdk/commit/f4ae30189c635859bd509e0b077b309403d7034a))
* upgrade RDK to use cosmos-sdk v0.46.10 ([#180](https://github.com/dymensionxyz/dymension-rdk/issues/180)) ([86713e4](https://github.com/dymensionxyz/dymension-rdk/commit/86713e4448559e2f1761e0bc75d5c9e62de8c4ce))



# [0.2.0-beta](https://github.com/dymensionxyz/dymension-rdk/compare/10bff865d6b3a676c574423f81805ed3929efcf7...v0.2.0-beta) (2023-02-15)


### Reverts

* Revert "removed vue and ts-client" ([10bff86](https://github.com/dymensionxyz/dymension-rdk/commit/10bff865d6b3a676c574423f81805ed3929efcf7))



