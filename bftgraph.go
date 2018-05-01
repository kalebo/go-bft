package bft

/*
#include <./include/bft.h>
#cgo CFLAGS: -std=c99
#cgo LDFLAGS:-L./libs/ -lbft -lJudy -ljemalloc -lm
*/
import "C"
import (
	"runtime"
	"sync"
	"unsafe"
	"fmt"
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
	fmt.Println("Freeing Graph!")
	C.free_cdbg(g.graph)
	Freed++
}

func (g *BFTGraph) SetMarking() {
	C.set_marking(g.graph)
}

func (g *BFTGraph) GetKmer(kmer string) *BFTKmer {
	k := NewBFTKmer(kmer, g)

	return k
}
