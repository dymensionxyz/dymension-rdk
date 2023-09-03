/*
Package dist defines a "wrapper" module around the Cosmos SDK's native
x/distribution module. In other words, it provides the exact same functionality as
the native module in that it simply embeds the native module. However, it
overrides the AllocateTokens function so the the allocation of rewards is done
towards the Sequencer, the stakers, and the community pool.
*/
package dist
