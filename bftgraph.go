package bft

/*
#include <bft/bft.h>
#cgo LDFLAGS: -lbft -lJudy -ljemalloc
*/
import "C"
import (
	"runtime"
	"sync"
	"unsafe"
)

var (
	Alloc int
	Freed int
	lock  sync.Mutex
)

////////////////
/// BFTGraph ///
////////////////

type BFTGraph struct {
	graph *C.BFT
	//sync.Mutex // not used yet
}

func NewBFTGraph(path string) *BFTGraph {
	cstrPath := C.CString(path)
	defer C.free(unsafe.Pointer(cstrPath))

	g := &BFTGraph{C.load_BFT(cstrPath)}
	runtime.SetFinalizer(g, (*BFTGraph).Free)
	Alloc++

	return g
}

func (g *BFTGraph) Free() {
	C.free_cdbg(g.graph)
	Freed++
}

func (g *BFTGraph) SetMarking() {
	C.set_marking(g.graph)
}

func (g *BFTGraph) GetKmer(kmer string) *BFTKmer {
	cstrKmer := C.CString(kmer)
	defer C.free(unsafe.Pointer(cstrKmer))

	k := &BFTKmer{C.get_kmer(cstrKmer, g.graph), g, 1} // Magic number is because an array of 1 is returned
	runtime.SetFinalizer(k, (*BFTKmer).Free)
	Alloc++

	return k
}
