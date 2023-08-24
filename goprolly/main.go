package main

import (
	"C"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	tree "github.com/kenlabs/go-ipld-prolly-trees/pkg/tree"
)

var dbtree *tree.ProllyTree
var rootcid cid.Cid

//export Mutate
func Mutate(ffiArgs *C.char) *C.char {
	if dbtree == nil {
		return ffiError(fmt.Errorf(
			"need to init",
		))
	}
	err := dbtree.Mutate()
	if err != nil {
		return ffiError(err)
	}
	return ffiOk(nil)
}

//export Initialize
func Initialize(ffiArgs *C.char) *C.char {

	ctx := context.Background()
	blockStore := blockstore.NewBlockstore(datastore.NewMapDatastore())
	nodeStore, err := tree.NewBlockNodeStore(blockStore, &tree.StoreConfig{CacheSize: 1 << 10})
	chunkConfig := tree.DefaultChunkConfig()
	framework, err := tree.NewFramework(ctx, nodeStore, chunkConfig, nil)
	dbtree, rootcid, err = framework.BuildTree(ctx)
	panicIfErr(err)
	treeConfig := dbtree.TreeConfig()
	linkPrefix := treeConfig.CidPrefix()

	fmt.Println("go-ipld-prolly-trees init: ", linkPrefix.MhLength, " root: ", rootcid)

	return ffiOk(map[string]interface{}{
		"root_cid": rootcid.String(),
	})
}

func ffiOk(obj interface{}) *C.char {
	return allocCJsonString(map[string]interface{}{
		"Ok": obj,
	})
}

func ffiError(err error) *C.char {
	return allocCJsonString(map[string]interface{}{
		"Err": err.Error(),
	})
}

func allocCJsonString(obj interface{}) *C.char {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return C.CString(fmt.Sprintf("failed to serialize to JSON: %s", err.Error()))
	}
	return C.CString(string(bytes))
}

// Based on https://dev.to/goncalorodrigues/using-go-generics-for-cleaner-code-4em1
func ffiDeserialize[T any](jsonChars *C.char) (*T, error) {
	jsonString := C.GoString(jsonChars)
	out := new(T)
	if err := json.Unmarshal([]byte(jsonString), out); err != nil {
		return nil, fmt.Errorf("invalid input: %s\nInput:\n%s", err.Error(), jsonString)
	}
	return out, nil
}
func main() {
}
func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
