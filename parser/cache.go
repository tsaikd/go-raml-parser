package parser

import (
	"bytes"
	"crypto/sha512"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
)

func loadFromCache(
	filePath string,
	fileData []byte,
	cacheDirectory string,
) (
	saveFunc func(RootDocument),
	rootdoc RootDocument,
	err error,
) {
	type cacheStruct struct {
		Hash    [sha512.Size]byte
		RootDoc RootDocument
	}

	hashpath := sha512.Sum512([]byte(filePath))
	hashfile := sha512.Sum512(fileData)
	cachename := filepath.Join(
		cacheDirectory,
		fmt.Sprintf("%x", hashpath[0:16]),
	)

	saveFunc = func(saveRootDoc RootDocument) {
		if err = os.MkdirAll(cacheDirectory, 0700); err != nil {
			return
		}

		var saveFile *os.File
		if saveFile, err = os.Create(cachename); err != nil {
			return
		}
		defer saveFile.Close()

		enc := gob.NewEncoder(saveFile)
		if err = enc.Encode(&cacheStruct{
			Hash:    hashfile,
			RootDoc: saveRootDoc,
		}); err != nil {
			return
		}
	}

	var cachefile *os.File
	if cachefile, err = os.Open(cachename); err != nil {
		return saveFunc, rootdoc, ErrorCacheNotFound.New(err)
	}
	defer cachefile.Close()

	dec := gob.NewDecoder(cachefile)
	cached := cacheStruct{}
	if err = dec.Decode(&cached); err != nil {
		return saveFunc, rootdoc, ErrorCacheNotFound.New(err)
	}

	if bytes.Equal(hashfile[:], cached.Hash[:]) {
		return nil, cached.RootDoc, nil
	}

	return saveFunc, rootdoc, ErrorCacheNotFound.New(nil)
}
