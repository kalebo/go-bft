# go-bft

This package provides a go wrapper around @GuillaumeHolley's [BloomFilterTrie](https://github.com/GuillaumeHolley/BloomFilterTrie).
It is a work in progress and I make no guarantees about it's correctness or utility. This is doubly true in regards to using goroutines with this library.

The final story for how to integrate go's build tools is still being figured out.  Currently the static libraries and header files for jemalloc, Judy, and bloomfiltertrie are included in the repo but this is obviously not ideal for version control. The benefit of this, however, is that you can just run `go get github.com/kalebo/go-bft` to install locally and then include the library in your own project. Your resulting binary would then only depend on having the standard shared library dependencies for a go binary e.g., libpthread, libc, and libm. 

The typical go way would be to have the c libraries installed by hand and then modify the CGO LDFLAGS and CFLAGS to point to where you installed the include files and the libraries. An easy way to do this would be to use conda after after adding bioconda channels e.g., `conda create -n bft bloomfiltertrie jellyfish`. However, the current state of the bloomfiltertrie means that a static library is not generated and thus any resulting go binary generated using this (go-bft) library will also need libbft to be installed locally. 
