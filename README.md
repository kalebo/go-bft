# go-bft

This package provides a go wrapper around @GuillaumeHolley's [BloomFilterTrie](https://github.com/GuillaumeHolley/BloomFilterTrie).
It is a work in progress and I make no guarantees about it's correctness or utility. This is doubly true in regards to using goroutines with this library.

You will need to have the jemalloc, Judy, and bloomfiltertrie libraries already installed in order to compile. The final story for integrating is still being figured out, but you can get all these packages by using conda. E.g., `conda create -n bft bloomfiltertrie jellyfish` and then entering that cond env. 

