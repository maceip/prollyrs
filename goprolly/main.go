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



var prollydb *tree.ProllyTree
var prollydbCid *cid.Cid

//export Mutate
func Mutate(ffiArgs *C.char) *C.char {
	if prollydb == nil {
		return ffiError(fmt.Errorf(
			"need to init",
		))
	}
	err := prollydb.Mutate()
	if err != nil {
		return ffiError(err)
	}
        return ffiOk(nil)
}

//export Initialize
func Initialize(ffiArgs *C.char) *C.char {

	fmt.Println("ZZ: init() was called")
	ctx := context.Background()
	blockStore := blockstore.NewBlockstore(datastore.NewMapDatastore())
	nodeStore, err := tree.NewBlockNodeStore(blockStore, &tree.StoreConfig{CacheSize: 1 << 10})
	chunkConfig := tree.DefaultChunkConfig()
	framework, err := tree.NewFramework(ctx, nodeStore, chunkConfig, nil)
	prollydb, prollydbCid, err := framework.BuildTree(ctx)
	panicIfErr(err)
	fmt.Println("ZZZ: ", prollydb, " ", prollydbCid)


/*
	cookies, err := ffiDeserialize[[]*http.Cookie](ffiCookies)
	if err != nil {
		return ffiError(err)
	}

	if scraper != nil {
		return ffiError(fmt.Errorf("already initialized in (Initialize was called twice)"))
	}

	scraper = twitterscraper.New()

	scraper.SetCookies(*cookies)

	// This is required for the scraper to know we are logged in
	if !scraper.IsLoggedIn() {
		return ffiError(fmt.Errorf("failed to initialize (cookies may be invalid)"))
	}
*/
	return ffiOk(nil)
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
	fmt.Println("ZZ: main() was called")
	ctx := context.Background()
	blockStore := blockstore.NewBlockstore(datastore.NewMapDatastore())
	nodeStore, err := tree.NewBlockNodeStore(blockStore, &tree.StoreConfig{CacheSize: 1 << 10})
	chunkConfig := tree.DefaultChunkConfig()
	framework, err := tree.NewFramework(ctx, nodeStore, chunkConfig, nil)
	prollydb, prollydbCid, err := framework.BuildTree(ctx)
	panicIfErr(err)
	fmt.Println("ZZZ: ", prollydb, " ", prollydbCid)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
