package bft

/*
#include <./include/bft.h>
#cgo CFLAGS: -std=c99
#cgo LDFLAGS:-L./libs/ -lbft -lJudy -ljemalloc -lm
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

///////////////
/// BFTKmer ///
///////////////

type BFTKmer struct {
	kmers  *C.BFT_kmer
	graph  *BFTGraph
	number int // number of kmers in array
	// Update assumption: number should always be 1; each element in 4-array will be collected separately
	//sync.Mutex // not used yet...
}

func NewBFTKmer(kmer string, graph *BFTGraph) *BFTKmer {
	cstrKmer := C.CString(kmer)
	defer C.free(unsafe.Pointer(cstrKmer))

	k := &BFTKmer{C.get_kmer(cstrKmer, graph.graph), graph, 1}
	runtime.SetFinalizer(k, (*BFTKmer).Free)
	Alloc++

	return k
}

func (k *BFTKmer) Free() {
	if k != nil && k.kmers != nil {
		fmt.Print("freeing:", C.GoString(k.kmers.kmer))
		//C.free_BFT_kmer(k.kmers, C.int(k.number)) // Why does this cause an inconsistent SIGSEGV but only freeing the content doesn't?
		C.free_BFT_kmer_content(k.kmers, C.int(k.number))
		fmt.Println("...Done")
	}
	Freed++
}

func (k *BFTKmer) String() string {
	return C.GoString(k.kmers.kmer)
}

func (k *BFTKmer) Exists() bool {
	return (C.is_kmer_in_cdbg(k.kmers)) != C.bool(false) // why can I query the graph without supplying any reference to the graph, but I can't in GetSuccessors?!
	// Also, can we talk about how ugly working with C99 bool types is with CGO. Am I missing something?
}

func (k *BFTKmer) GetSuccessors() []*BFTKmer {
	var arrayPtr *C.BFT_kmer = C.get_successors(k.kmers, k.graph.graph)

	// Makes a go slice backed by a c array without copy
	kmerSlice := (*[1 << 30]C.BFT_kmer)(unsafe.Pointer(arrayPtr))[:4:4] // length of four

	result := make([]*BFTKmer, 0)
	for i := 0; i < 4; i++ {
		kmer := &BFTKmer{&kmerSlice[i], k.graph, 1}
		runtime.SetFinalizer(kmer, (*BFTKmer).Free)
		Alloc++

		if kmer.Exists() {
			result = append(result, kmer)
		}
	}

	return result
}
