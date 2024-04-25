/*
Package staking defines a "wrapper" module around the Cosmos SDK's native
x/staking module. In other words, it provides the exact same functionality as
the native module in that it simply embeds the native module.
However, it overrides `EndBlock` method.
Specifically, these method perform no-ops and return no validator set updates, as validator sets
are tracked by the Sequencers module
*/
package staking
