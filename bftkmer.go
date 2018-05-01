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

////////////////////
/// BFTKmerArray ///
////////////////////

type BFTKmerArray struct {
	arrayPtr  *C.BFT_kmer // pointing to the array start
	count int // count of characters in the associated array
	graph *BFTGraph // holding on to reference to prevent preemptive GC of graph
}

func NewBFTKmerArray(kmers *C.BFT_kmer, count int, graph *BFTGraph) *BFTKmerArray {
	a := &BFTKmerArray{kmers, count, graph}
	runtime.SetFinalizer(a, (*BFTKmerArray).Free)
	Alloc++
	return a
}

func (a *BFTKmerArray) Free() {
	fmt.Print("freeing: kmer array", " of count ", a.count)
	C.free_BFT_kmer(a.arrayPtr, C.int(a.count))
	fmt.Println("...Done")
	Freed++
}

func (a *BFTKmerArray) RegisterFinalization(kmer *BFTKmer) {
	fmt.Println("Child was finalized:", C.GoString(kmer.kmers.kmer))
}

///////////////
/// BFTKmer ///
///////////////

type BFTKmer struct {
	kmers  *C.BFT_kmer
	graph  *BFTGraph
	containingArray *BFTKmerArray
}

func NewBFTKmer(kmer string, graph *BFTGraph) *BFTKmer {
	cstrKmer := C.CString(kmer)
	defer C.free(unsafe.Pointer(cstrKmer))

	kmerPtr := C.get_kmer(cstrKmer, graph.graph)
	kmerArr := NewBFTKmerArray(kmerPtr, 1, graph)

	k := &BFTKmer{kmerPtr, graph, kmerArr}
	runtime.SetFinalizer(k, (*BFTKmer).Free)
	Alloc++

	return k
}

func (k *BFTKmer) Free() {
	k.containingArray.RegisterFinalization(k)
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
	kmerArr := NewBFTKmerArray(arrayPtr, 4, k.graph)

	// Makes a go slice backed by a c array without copy
	kmerSlice := (*[1 << 30]C.BFT_kmer)(unsafe.Pointer(arrayPtr))[:4:4] // length of four

	result := make([]*BFTKmer, 0)
	for i := 0; i < 4; i++ {
		kmer := &BFTKmer{&kmerSlice[i], k.graph, kmerArr}
		runtime.SetFinalizer(kmer, (*BFTKmer).Free)
		Alloc++

		if kmer.Exists() {
			result = append(result, kmer)
		}
	}

	return result
}
