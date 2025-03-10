#  (2025-03-10)


### Bug Fixes

* added keyring-backend flag for gentx_seq ([#239](https://github.com/dymensionxyz/dymension-rdk/issues/239)) ([7bf5d44](https://github.com/dymensionxyz/dymension-rdk/commit/7bf5d44d353a0eb5e6823542c44dd0cc94bed921))
* **allocation:** do not destroy fee rewards if allocating to proposer fails ([#358](https://github.com/dymensionxyz/dymension-rdk/issues/358)) ([bfd0943](https://github.com/dymensionxyz/dymension-rdk/commit/bfd09432d61eb2c33e29621a1ab868661aeda50c))
* **ante:** skip whitelist check for IBC packet messages ([#636](https://github.com/dymensionxyz/dymension-rdk/issues/636)) ([8695a42](https://github.com/dymensionxyz/dymension-rdk/commit/8695a42df50671051da6b5080e727807c93a30e1))
* Bump golang.org/x/net from 0.3.0 to 0.7.0 ([#150](https://github.com/dymensionxyz/dymension-rdk/issues/150)) ([1963021](https://github.com/dymensionxyz/dymension-rdk/commit/1963021c7eced52dafb580057716a32dcc53b327))
* changed Int64 to string when logging ([#220](https://github.com/dymensionxyz/dymension-rdk/issues/220)) ([c6cfd1a](https://github.com/dymensionxyz/dymension-rdk/commit/c6cfd1a5f6502e7d24d036709dae53823a5afbe4))
* changed the ibc_go replace directive to the release/v5.1.x_dymint branch ([#185](https://github.com/dymensionxyz/dymension-rdk/issues/185)) ([bb3ef71](https://github.com/dymensionxyz/dymension-rdk/commit/bb3ef71a4211e5904a8fcce55080c36d1d5f118c))
* ci lint failure  ([#165](https://github.com/dymensionxyz/dymension-rdk/issues/165)) ([1b92970](https://github.com/dymensionxyz/dymension-rdk/commit/1b9297044204c7b16868008942a9be82b3ca0813))
* clarified whitelisted relayer logs ([#595](https://github.com/dymensionxyz/dymension-rdk/issues/595)) ([36d0a4d](https://github.com/dymensionxyz/dymension-rdk/commit/36d0a4d5d2381cc75bbfa4060d5d78355b83fda7))
* **code standards:** removes blocked addresses from distr keeper (unused) ([#476](https://github.com/dymensionxyz/dymension-rdk/issues/476)) ([843d029](https://github.com/dymensionxyz/dymension-rdk/commit/843d029d48e8b0fd9806edb1a29236e8a67840a3))
* **code standards:** use https://github.com/dymensionxyz/sdk-utils ([#464](https://github.com/dymensionxyz/dymension-rdk/issues/464)) ([58f9faf](https://github.com/dymensionxyz/dymension-rdk/commit/58f9fafe43fb6519347c3b469f7320695091edba))
* correct Any.TypeURL in consensus msg response ([#596](https://github.com/dymensionxyz/dymension-rdk/issues/596)) ([15957e8](https://github.com/dymensionxyz/dymension-rdk/commit/15957e8e3a1cc1ceed1a2d70e685192bd279aa02))
* **denommetadata:** fixed ibc voucher denom query ([#378](https://github.com/dymensionxyz/dymension-rdk/issues/378)) ([bf14e28](https://github.com/dymensionxyz/dymension-rdk/commit/bf14e28256c444930195b3a6df760d3626c205fa))
* **deps:** bump dymint v.1.3.0-rc04 ([#678](https://github.com/dymensionxyz/dymension-rdk/issues/678)) ([fdf0ed7](https://github.com/dymensionxyz/dymension-rdk/commit/fdf0ed7fa2a707df47aa0b50a286a899df2557df))
* **doc:** adds a doc.go for hub genesis ([#478](https://github.com/dymensionxyz/dymension-rdk/issues/478)) ([23b96de](https://github.com/dymensionxyz/dymension-rdk/commit/23b96dea3641bd6b69631b4284ecde7a771400e7))
* **drs:** fetch only relevant param on upgrade to avoid panic on new params ([#628](https://github.com/dymensionxyz/dymension-rdk/issues/628)) ([62fded9](https://github.com/dymensionxyz/dymension-rdk/commit/62fded99c5dead75a9b0df67c3b22b9090c1bb08))
* dymint blockmanager fix for rollback command ([#575](https://github.com/dymensionxyz/dymension-rdk/issues/575)) ([b043136](https://github.com/dymensionxyz/dymension-rdk/commit/b0431362bfc16fa9b0b7f522e063a58de3e38bfc))
* dymint reset-state command not working ([#309](https://github.com/dymensionxyz/dymension-rdk/issues/309)) ([4f5a707](https://github.com/dymensionxyz/dymension-rdk/commit/4f5a70734530f4da1c8980cdeb9b59d5a2148577))
* **epochs:** counting epoch from first block rather than genesis time ([#639](https://github.com/dymensionxyz/dymension-rdk/issues/639)) ([bcb6034](https://github.com/dymensionxyz/dymension-rdk/commit/bcb6034a7c14c38335c9715562d56e8dccc06f06))
* evm module to use sequencers module instead of staking ([#223](https://github.com/dymensionxyz/dymension-rdk/issues/223)) ([315aa83](https://github.com/dymensionxyz/dymension-rdk/commit/315aa83afd3daf381608d59af088f405ad05b3ad))
* fixed bug when modifying staking params in init script ([#176](https://github.com/dymensionxyz/dymension-rdk/issues/176)) ([96b35de](https://github.com/dymensionxyz/dymension-rdk/commit/96b35deb02bd0b09d8c787a7e0856e93eef598a3))
* fixed epochs proto API path ([#236](https://github.com/dymensionxyz/dymension-rdk/issues/236)) ([0013f9f](https://github.com/dymensionxyz/dymension-rdk/commit/0013f9f12bb52c0e2a2547c3508e242764083ec7))
* fixed ERC20 support ([#140](https://github.com/dymensionxyz/dymension-rdk/issues/140)) ([f4aae00](https://github.com/dymensionxyz/dymension-rdk/commit/f4aae00b8f32d1502d6f5b7fe1d1c1389a7f8b68))
* fixed overwritten signal handler in `cmd/start` ([#195](https://github.com/dymensionxyz/dymension-rdk/issues/195)) ([72bd37f](https://github.com/dymensionxyz/dymension-rdk/commit/72bd37fe31f6229a7f59087b0b264c194e0a02a8))
* **gasless:** removed custom logic as there's native implementation ([#564](https://github.com/dymensionxyz/dymension-rdk/issues/564)) ([82c4d5f](https://github.com/dymensionxyz/dymension-rdk/commit/82c4d5f8c09365b20b4378c0cc459b414fd306e8))
* **gasless:** usage identifier validation and mapping to gas tank IDs ([#608](https://github.com/dymensionxyz/dymension-rdk/issues/608)) ([ee8f0b4](https://github.com/dymensionxyz/dymension-rdk/commit/ee8f0b450c58116e78165526e7a556976a90e35a))
* generated proto ([#616](https://github.com/dymensionxyz/dymension-rdk/issues/616)) ([44f850f](https://github.com/dymensionxyz/dymension-rdk/commit/44f850f67b9e65e5b5598c4ccaa622a33bf0b456))
* **genesis bridge:** better UX for chains without genesis transfers ([#458](https://github.com/dymensionxyz/dymension-rdk/issues/458)) ([fd8806d](https://github.com/dymensionxyz/dymension-rdk/commit/fd8806da3d31e6dffbc96ccd5c8b0f9caf83c54d))
* genesis validation ([#618](https://github.com/dymensionxyz/dymension-rdk/issues/618)) ([b7e7afd](https://github.com/dymensionxyz/dymension-rdk/commit/b7e7afdd4c34cce46c2494f2ba979050dca3962b))
* **genesis_bridge:** added genesisInfo query ([#573](https://github.com/dymensionxyz/dymension-rdk/issues/573)) ([1899432](https://github.com/dymensionxyz/dymension-rdk/commit/189943271174cb853d5cdf61efef97abd827b812))
* **hub genesis:** check if oustanding seq num exists before decrementing unacked transfers cnt ([#477](https://github.com/dymensionxyz/dymension-rdk/issues/477)) ([21faaaf](https://github.com/dymensionxyz/dymension-rdk/commit/21faaaf4ce2516d78d2298e2433d7d947a7b4269))
* **hubgenesis:** check for duplicates in genesis info accounts ([#598](https://github.com/dymensionxyz/dymension-rdk/issues/598)) ([c740fcc](https://github.com/dymensionxyz/dymension-rdk/commit/c740fccae4202a4d3b9f26e6c19a9f5e99f347fe))
* **hub:** simplify registered denom on hub state keeping ([#555](https://github.com/dymensionxyz/dymension-rdk/issues/555)) ([3fe31b2](https://github.com/dymensionxyz/dymension-rdk/commit/3fe31b2db4b2b6f4a8474b4cb2218641b1d6e4ba))
* **hub:** track hub's registered denoms in a separate store ([#579](https://github.com/dymensionxyz/dymension-rdk/issues/579)) ([a1ddc11](https://github.com/dymensionxyz/dymension-rdk/commit/a1ddc11f9928dd3b2320020115c73e6b9d315180))
* **ibc:** allow whitelist MsgSendTransfer for channel lifecycle ([#699](https://github.com/dymensionxyz/dymension-rdk/issues/699)) ([9d27fbd](https://github.com/dymensionxyz/dymension-rdk/commit/9d27fbd7f0fcf741cdbea33e7647fd24cc85f219))
* **ibc:** wrap IBC error acknowledgement with an error event ([#542](https://github.com/dymensionxyz/dymension-rdk/issues/542)) ([a28124a](https://github.com/dymensionxyz/dymension-rdk/commit/a28124a685302b569442e1d5ce360e93df0e91f5))
* **mint:** added missing cli query for current inflation ([#366](https://github.com/dymensionxyz/dymension-rdk/issues/366)) ([6ecadd8](https://github.com/dymensionxyz/dymension-rdk/commit/6ecadd818429e1870da7e6f6e044bc87f900ec8b))
* only allow whitelisted relayer account to send IBC relayer messages ([#614](https://github.com/dymensionxyz/dymension-rdk/issues/614)) ([5cc3d81](https://github.com/dymensionxyz/dymension-rdk/commit/5cc3d819013c7af5abbdd8cbc8cd4f70df60b674))
* panic in inspect-state command ([#597](https://github.com/dymensionxyz/dymension-rdk/issues/597)) ([d3ad5c8](https://github.com/dymensionxyz/dymension-rdk/commit/d3ad5c8c25a40ed2abc50f218d5fad1932b16236))
* **proto:** clean up unnecessary deps and parts of makefile ([#441](https://github.com/dymensionxyz/dymension-rdk/issues/441)) ([8fdf7ef](https://github.com/dymensionxyz/dymension-rdk/commit/8fdf7ef2605047cd9cd41f797e6dc421fb2c9900))
* **proto:** use buf, like we do on dymension repo ([#432](https://github.com/dymensionxyz/dymension-rdk/issues/432)) ([f2d2ad1](https://github.com/dymensionxyz/dymension-rdk/commit/f2d2ad1ff599aa0f4c54350918e43b418d77fff1))
* RDK upgrade handler fix from v0.45 to v0.46 ([#192](https://github.com/dymensionxyz/dymension-rdk/issues/192)) ([d632283](https://github.com/dymensionxyz/dymension-rdk/commit/d63228382d1df08706d07ef668cb1a06cf404333))
* **rollappparams:** change DRS version from commit to int ([#585](https://github.com/dymensionxyz/dymension-rdk/issues/585)) ([88f668f](https://github.com/dymensionxyz/dymension-rdk/commit/88f668faab661c484bb6c7125603b732800c70a9))
* **rollappparams:** changed `version` and `da` to be overridable during build  ([#534](https://github.com/dymensionxyz/dymension-rdk/issues/534)) ([7afd5f7](https://github.com/dymensionxyz/dymension-rdk/commit/7afd5f77560e33dfeee3d3c5e20a8fc07e4f7f65))
* **rollappparams:** prevent setting drs version through gov prop params change ([#635](https://github.com/dymensionxyz/dymension-rdk/issues/635)) ([b3df828](https://github.com/dymensionxyz/dymension-rdk/commit/b3df828e7611db51673cd265ff1d6c34b0842b22))
* **rollappparams:** verify denom for gov proposal gas price change ([#648](https://github.com/dymensionxyz/dymension-rdk/issues/648)) ([516b6b6](https://github.com/dymensionxyz/dymension-rdk/commit/516b6b6d3cc7af264c354c134e3b751a5a392cf4))
* **rollback:** add header hash to UpdateStateFromApp ([#582](https://github.com/dymensionxyz/dymension-rdk/issues/582)) ([5b22225](https://github.com/dymensionxyz/dymension-rdk/commit/5b22225071b4841ba70e4c3558639942e909f3a4))
* sequencer module returns two objects in initchain ([#237](https://github.com/dymensionxyz/dymension-rdk/issues/237)) ([76c5329](https://github.com/dymensionxyz/dymension-rdk/commit/76c5329f3b3801c5ce3ab61a8dfc1a5a8319bc9c))
* **sequencer:** added HistoricalEntries param validation ([#369](https://github.com/dymensionxyz/dymension-rdk/issues/369)) ([8e4b891](https://github.com/dymensionxyz/dymension-rdk/commit/8e4b891e68304912b839c2b3ef2a733d91426ce2))
* **sequencer:** doc improvement tweaks from previous merged PR feedback ([#526](https://github.com/dymensionxyz/dymension-rdk/issues/526)) ([2c6d706](https://github.com/dymensionxyz/dymension-rdk/commit/2c6d706d8e7704ffe277d8b2ccd5f500a5300f08))
* **sequencers:** sequencers genesis operator address not exported ([#367](https://github.com/dymensionxyz/dymension-rdk/issues/367)) ([2a1f308](https://github.com/dymensionxyz/dymension-rdk/commit/2a1f308c31141143c20fadee37578c5be0df4f87))
* **tests:** update ibcKeeper use SequencersKeeper instead of stakingKeeper ([#418](https://github.com/dymensionxyz/dymension-rdk/issues/418)) ([312bd67](https://github.com/dymensionxyz/dymension-rdk/commit/312bd6760dab305ee8f183454b2aab52b1288743))
* **timeupgrade:** validating time in software upgrade proposal ([#652](https://github.com/dymensionxyz/dymension-rdk/issues/652)) ([68cd3e1](https://github.com/dymensionxyz/dymension-rdk/commit/68cd3e1f07b70ff8112ef74299538b2f63227981))
* **tokenfactory:** fixed tokenfactory denom-metadata override ([#634](https://github.com/dymensionxyz/dymension-rdk/issues/634)) ([14119cc](https://github.com/dymensionxyz/dymension-rdk/commit/14119cc73f17fd9a57b8ff954c03b36b10c65d62))
* **tokenfactory:** remove denom validation upon change admin ([#675](https://github.com/dymensionxyz/dymension-rdk/issues/675)) ([848d94f](https://github.com/dymensionxyz/dymension-rdk/commit/848d94f2003559db78a881fb5db79f4819f82e60))
* **tokenfactory:** sets missing fields in the denommetadata ([#651](https://github.com/dymensionxyz/dymension-rdk/issues/651)) ([35ed372](https://github.com/dymensionxyz/dymension-rdk/commit/35ed372cb049264519998e51fcd41b28dde989f9))
* **transfer genesis:** do not allow ibc denoms ([#466](https://github.com/dymensionxyz/dymension-rdk/issues/466)) ([74f2127](https://github.com/dymensionxyz/dymension-rdk/commit/74f21279f16d3a5608ee442406b439024b18f474))
* update dymint and cometbft ([#590](https://github.com/dymensionxyz/dymension-rdk/issues/590)) ([ef99c12](https://github.com/dymensionxyz/dymension-rdk/commit/ef99c12f13b49cdf843ce2d30da101f0b76402b9))
* **upgrade:** drs upgrade plan renamed ([#633](https://github.com/dymensionxyz/dymension-rdk/issues/633)) ([68ba9bc](https://github.com/dymensionxyz/dymension-rdk/commit/68ba9bcc36c0dcf7ac4d360d4b5c3a01d86cf994))


### Features

* add automatic erc20 conversion ([#154](https://github.com/dymensionxyz/dymension-rdk/issues/154)) ([d719490](https://github.com/dymensionxyz/dymension-rdk/commit/d71949027affe46ef5d4bec520fe33a55ab4558d))
* add channel 0 check for tokenless ([#642](https://github.com/dymensionxyz/dymension-rdk/issues/642)) ([77cdd63](https://github.com/dymensionxyz/dymension-rdk/commit/77cdd63f165303924d209dcb214fdaec1489b9fb))
* Add genesis import/export for timeupgrade ([#660](https://github.com/dymensionxyz/dymension-rdk/issues/660)) ([2976834](https://github.com/dymensionxyz/dymension-rdk/commit/2976834aa4571fc11d9c8c7b2803bf82e0874cd4))
* add GenesisChecksum to genesis info ([#587](https://github.com/dymensionxyz/dymension-rdk/issues/587)) ([89acf9e](https://github.com/dymensionxyz/dymension-rdk/commit/89acf9ec41b946014b2cde9e7eb63efe79ec18e2))
* add get ibc denom base on denom trace query ([#376](https://github.com/dymensionxyz/dymension-rdk/issues/376)) ([807717b](https://github.com/dymensionxyz/dymension-rdk/commit/807717b3219609a68e14d4f11100c539db001008))
* add readme ([#199](https://github.com/dymensionxyz/dymension-rdk/issues/199)) ([935d9ee](https://github.com/dymensionxyz/dymension-rdk/commit/935d9eec25f194cbe086261ac5f7aa77d39ca7c6))
* add tokenfactory module ([#612](https://github.com/dymensionxyz/dymension-rdk/issues/612)) ([0862741](https://github.com/dymensionxyz/dymension-rdk/commit/08627417ea7d304e93ca6392d3f194d92e790fb4))
* adding_inspect_command ([#289](https://github.com/dymensionxyz/dymension-rdk/issues/289)) ([a58be89](https://github.com/dymensionxyz/dymension-rdk/commit/a58be893b8aa113b302aeb308036858b80f7958c))
* change ibc script to support dymension relayer ([#156](https://github.com/dymensionxyz/dymension-rdk/issues/156)) ([01c204c](https://github.com/dymensionxyz/dymension-rdk/commit/01c204c87fd0c28edbe577bd0255e231eb28333c))
* change minting to be epoch based ([#198](https://github.com/dymensionxyz/dymension-rdk/issues/198)) ([404bb2b](https://github.com/dymensionxyz/dymension-rdk/commit/404bb2b468b34c607757324e23fa52572ad879f6))
* changed log rotate max size to be configurable vs hardcoded ([#107](https://github.com/dymensionxyz/dymension-rdk/issues/107)) ([1c40dfd](https://github.com/dymensionxyz/dymension-rdk/commit/1c40dfd5ba464ac3df3cea291fb6a1bbb242e765))
* **ci:** Add changelog workflow ([#393](https://github.com/dymensionxyz/dymension-rdk/issues/393)) ([568a98f](https://github.com/dymensionxyz/dymension-rdk/commit/568a98f0ad1b4644bff91b29a278566d017d98dc))
* clean up scripts folder ([#149](https://github.com/dymensionxyz/dymension-rdk/issues/149)) ([de8b201](https://github.com/dymensionxyz/dymension-rdk/commit/de8b201a6f5f157207d5518cc94608df123a40f2))
* **cli:** rollback command added to move state back to specific height ([#437](https://github.com/dymensionxyz/dymension-rdk/issues/437)) ([bccd5b9](https://github.com/dymensionxyz/dymension-rdk/commit/bccd5b90b6bf61edc0dc2fbb507c45fa1ecd9762))
* common consensus msgs processing ([#570](https://github.com/dymensionxyz/dymension-rdk/issues/570)) ([234d438](https://github.com/dymensionxyz/dymension-rdk/commit/234d438a56c792c7e8975cfa965065c66d16e3e1))
* consensus messages ([#563](https://github.com/dymensionxyz/dymension-rdk/issues/563)) ([baa5f01](https://github.com/dymensionxyz/dymension-rdk/commit/baa5f01465dc47177c1e73fff3808de36f6410d3))
* **da:** enable changing da in rollappparams ([#697](https://github.com/dymensionxyz/dymension-rdk/issues/697)) ([f4ba9f5](https://github.com/dymensionxyz/dymension-rdk/commit/f4ba9f5bd096fc20132bb9c1bf8d4a7f4696ef72))
* **denommetadata:** Add support for providing denom-trace and creation of multiple denom-metadatas  ([#384](https://github.com/dymensionxyz/dymension-rdk/issues/384)) ([c720640](https://github.com/dymensionxyz/dymension-rdk/commit/c720640d7d392d58fb191852cd8f4a5d85afc04e))
* **denommetadata:** inject denom metadata to ibc transfers from RA to Hub ([#455](https://github.com/dymensionxyz/dymension-rdk/issues/455)) ([c98e3e6](https://github.com/dymensionxyz/dymension-rdk/commit/c98e3e67c156b8cdc68dc1b1842c9c4d4c1e6033))
* **deonmetadata:** enable testing for denom metadata creation with erc20 hook ([#343](https://github.com/dymensionxyz/dymension-rdk/issues/343)) ([8e3dc66](https://github.com/dymensionxyz/dymension-rdk/commit/8e3dc66c8c6d6894cbf2d5d12112772069ea8774))
* **dividends:** module implementation ([#686](https://github.com/dymensionxyz/dymension-rdk/issues/686)) ([6583a5a](https://github.com/dymensionxyz/dymension-rdk/commit/6583a5a31946d787a6bf26ce4fc4c1b6010e48c8))
* Enhanced log config and enable log level override on a module basis ([#109](https://github.com/dymensionxyz/dymension-rdk/issues/109)) ([7292826](https://github.com/dymensionxyz/dymension-rdk/commit/729282680fe93335ff7c40d5e5c6586936222aaa))
* **gasless:** add gasless module  ([#427](https://github.com/dymensionxyz/dymension-rdk/issues/427)) ([7b38c5a](https://github.com/dymensionxyz/dymension-rdk/commit/7b38c5a132513cbf3079a1b45b5187149c994606))
* **genesis bridge:** genesis transfers ([#449](https://github.com/dymensionxyz/dymension-rdk/issues/449)) ([b3343e1](https://github.com/dymensionxyz/dymension-rdk/commit/b3343e1bece799a4f8e21e6fc87d4e328b460715))
* **genesis_bridge:** revised genesis bridge impl ([#565](https://github.com/dymensionxyz/dymension-rdk/issues/565)) ([623826e](https://github.com/dymensionxyz/dymension-rdk/commit/623826e0daa4d365863a64f205e8ce59405e9e88))
* **hub-genesis:** added a query to generate genesis bridge data ([#609](https://github.com/dymensionxyz/dymension-rdk/issues/609)) ([ffec6e8](https://github.com/dymensionxyz/dymension-rdk/commit/ffec6e81a1ca9b485bac0bcda53de1d68492ba7f))
* **hub-genesis:** json encoding GenesisBridgeData ([#576](https://github.com/dymensionxyz/dymension-rdk/issues/576)) ([b250db7](https://github.com/dymensionxyz/dymension-rdk/commit/b250db7af580e891036d05496958a9f08712e4c3))
* **hub-genesis:** whitelisted address trigger genesis transfer ([#625](https://github.com/dymensionxyz/dymension-rdk/issues/625)) ([3d910ef](https://github.com/dymensionxyz/dymension-rdk/commit/3d910ef66fc5566138784184a13dafacd10529ab))
* **hubgenesis:** support non native token rollapps ([#638](https://github.com/dymensionxyz/dymension-rdk/issues/638)) ([ab95cfd](https://github.com/dymensionxyz/dymension-rdk/commit/ab95cfdbec35d395d320ea3e2ccd5afde23e0045))
* **ibc transfer:** Register IBC denom on transfer ([#433](https://github.com/dymensionxyz/dymension-rdk/issues/433)) ([6aaf3c7](https://github.com/dymensionxyz/dymension-rdk/commit/6aaf3c76c70a024f4fc57640278bf83b29d34124))
* make update client, recv, ack, timeout free without whitelist ([#694](https://github.com/dymensionxyz/dymension-rdk/issues/694)) ([684aaf2](https://github.com/dymensionxyz/dymension-rdk/commit/684aaf2625cfe1b53492857e764b7d33c52b62a7))
* migration for mainnet rollapps (nim + mande) to upgrade to 3D ([#667](https://github.com/dymensionxyz/dymension-rdk/issues/667)) ([5198f72](https://github.com/dymensionxyz/dymension-rdk/commit/5198f720fca43221d930953d4c1a9082bbde22a9))
* mint add query cmd for current inflation ([#328](https://github.com/dymensionxyz/dymension-rdk/issues/328)) ([1c49688](https://github.com/dymensionxyz/dymension-rdk/commit/1c49688dc88de8cb70a358ba80eced1a9062543c))
* move dymint related commands to dymint repo ([#206](https://github.com/dymensionxyz/dymension-rdk/issues/206)) ([961dc98](https://github.com/dymensionxyz/dymension-rdk/commit/961dc9886589c58dea880125da2574021cab2fb0))
* new rollapp consensus params module ([#495](https://github.com/dymensionxyz/dymension-rdk/issues/495)) ([287e367](https://github.com/dymensionxyz/dymension-rdk/commit/287e3677e7f180ca9e7278846b9a49a4b7751da3))
* passing max_log_size to the logger ([#111](https://github.com/dymensionxyz/dymension-rdk/issues/111)) ([ef2d364](https://github.com/dymensionxyz/dymension-rdk/commit/ef2d36476034ced14f2460719f0c90dca523f791))
* permissioned creation of token metadata ([#327](https://github.com/dymensionxyz/dymension-rdk/issues/327)) ([ffe8e06](https://github.com/dymensionxyz/dymension-rdk/commit/ffe8e06500e248b4e45097700c3b296131971084))
* rollapp genesis token locking upon channel creation ([#333](https://github.com/dymensionxyz/dymension-rdk/issues/333)) ([4bbdcda](https://github.com/dymensionxyz/dymension-rdk/commit/4bbdcdadb10c44d8458b48ac1c552fa9fcea3321))
* **rollapp params:** add ability to query rollapp params ([#521](https://github.com/dymensionxyz/dymension-rdk/issues/521)) ([d117bfe](https://github.com/dymensionxyz/dymension-rdk/commit/d117bfe31ea3ecdc9170f4e076f0e9ab171327bc))
* **rollapp params:** commit param renamed to version ([#520](https://github.com/dymensionxyz/dymension-rdk/issues/520)) ([78bf4ad](https://github.com/dymensionxyz/dymension-rdk/commit/78bf4ad4eccfc496bc7e6c103e53bc3dec3ecbc9))
* **rollapp params:** remove maxgas from rollappparams ([#524](https://github.com/dymensionxyz/dymension-rdk/issues/524)) ([4b256cd](https://github.com/dymensionxyz/dymension-rdk/commit/4b256cd9fd36d214e768a5dbf668d542c013c77f))
* **rollapp:** Add pruning command ([#484](https://github.com/dymensionxyz/dymension-rdk/issues/484)) ([f87d7fc](https://github.com/dymensionxyz/dymension-rdk/commit/f87d7fc1389a0eec0d5793a3e6897b28749a58c0))
* **rollapparams:** remove block params from rollapp params ([#545](https://github.com/dymensionxyz/dymension-rdk/issues/545)) ([79c7d11](https://github.com/dymensionxyz/dymension-rdk/commit/79c7d1157304b60b2ce5ff70cfd4de45188b5380))
* **rollappparam:** block da change ([#688](https://github.com/dymensionxyz/dymension-rdk/issues/688)) ([2e22575](https://github.com/dymensionxyz/dymension-rdk/commit/2e22575c4e5cfb05e7ae0f6fa5a3d4b32faab172))
* **rollappparams:** added param for min gas price and respective ante handlers ([#611](https://github.com/dymensionxyz/dymension-rdk/issues/611)) ([30ef719](https://github.com/dymensionxyz/dymension-rdk/commit/30ef7199888d33da4d48ead777caa69a04a574fe))
* **rollback:** enable rollback support ([#656](https://github.com/dymensionxyz/dymension-rdk/issues/656)) ([7fe0f88](https://github.com/dymensionxyz/dymension-rdk/commit/7fe0f886c5027cfb93b42fe1978c2c622eae76b0))
* **sequencers:** add queries for reward address and whitelisted relayers ([#602](https://github.com/dymensionxyz/dymension-rdk/issues/602)) ([def6322](https://github.com/dymensionxyz/dymension-rdk/commit/def6322e4345a8af5d8295a59c1ea3f1b4c15ac8))
* **sequencers:** added MsgUpsertSequencer ([#568](https://github.com/dymensionxyz/dymension-rdk/issues/568)) ([69c1b7c](https://github.com/dymensionxyz/dymension-rdk/commit/69c1b7cf38c67ac838a7543f924f5db8055ffe93))
* **sequencer:** set/update reward addr ([#510](https://github.com/dymensionxyz/dymension-rdk/issues/510)) ([28628bc](https://github.com/dymensionxyz/dymension-rdk/commit/28628bcb329ac7b4cdd6d19f1c43f26d03410f84))
* **sequencers:** support upsert sequencer consensus msg ([#572](https://github.com/dymensionxyz/dymension-rdk/issues/572)) ([5450517](https://github.com/dymensionxyz/dymension-rdk/commit/545051749bc0fcd9f8d8f27c2bb4d582e960f1e2))
* **server:** add helper to set default pruning settings ([#488](https://github.com/dymensionxyz/dymension-rdk/issues/488)) ([005f2bb](https://github.com/dymensionxyz/dymension-rdk/commit/005f2bb6dd18f40e97046f592dd7cdbd94485024))
* simple way to add sequencer on genesis ([#183](https://github.com/dymensionxyz/dymension-rdk/issues/183)) ([53f56f9](https://github.com/dymensionxyz/dymension-rdk/commit/53f56f9d9aef67c613b94b88f78ff107f900f5d5))
* support getting active sequencers from dymint ([#147](https://github.com/dymensionxyz/dymension-rdk/issues/147)) ([b3c8f82](https://github.com/dymensionxyz/dymension-rdk/commit/b3c8f8262c48670eb8cc48e743b09d2804f70ffd))
* support reading dymint configuration from file ([#187](https://github.com/dymensionxyz/dymension-rdk/issues/187)) ([5feb350](https://github.com/dymensionxyz/dymension-rdk/commit/5feb3504c4c4820222b5c6ce18a84d029c0f0154))
* support stakers on genesis ([#260](https://github.com/dymensionxyz/dymension-rdk/issues/260)) ([4632268](https://github.com/dymensionxyz/dymension-rdk/commit/463226890c589dee8d1b41fda246e44130a11e4f))
* time based upgrades ([#543](https://github.com/dymensionxyz/dymension-rdk/issues/543)) ([c620e0b](https://github.com/dymensionxyz/dymension-rdk/commit/c620e0b915d0329dea508cd25735dc4918c6b3f0))
* update rdk to use `cometbft` ([#209](https://github.com/dymensionxyz/dymension-rdk/issues/209)) ([f4ae301](https://github.com/dymensionxyz/dymension-rdk/commit/f4ae30189c635859bd509e0b077b309403d7034a))
* upgrade RDK to use cosmos-sdk v0.46.10 ([#180](https://github.com/dymensionxyz/dymension-rdk/issues/180)) ([86713e4](https://github.com/dymensionxyz/dymension-rdk/commit/86713e4448559e2f1761e0bc75d5c9e62de8c4ce))
* **upgrade:** time-based upgrade refactor ([#631](https://github.com/dymensionxyz/dymension-rdk/issues/631)) ([6627c3d](https://github.com/dymensionxyz/dymension-rdk/commit/6627c3ddeb4b2e24dfd6e69ce6463bc38dd0be52))
* **utils:** dymint p2p status command added ([#420](https://github.com/dymensionxyz/dymension-rdk/issues/420)) ([00991ea](https://github.com/dymensionxyz/dymension-rdk/commit/00991eaefe96cad39ce272ef70060b7dd48ff5aa))
* validate genesis command ([#620](https://github.com/dymensionxyz/dymension-rdk/issues/620)) ([6e0296b](https://github.com/dymensionxyz/dymension-rdk/commit/6e0296bf3b2c456ae525d1bdc2dd84b7630d8cef))


### Reverts

* Revert "removed vue and ts-client" ([10bff86](https://github.com/dymensionxyz/dymension-rdk/commit/10bff865d6b3a676c574423f81805ed3929efcf7))



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



