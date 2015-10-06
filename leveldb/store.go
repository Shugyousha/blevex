//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package leveldb

import (
	"fmt"
	"sync"

	"github.com/blevesearch/bleve/index/store"
	"github.com/blevesearch/bleve/registry"
	"github.com/jmhodges/levigo"
)

const Name = "leveldb"

type Store struct {
	path string
	opts *levigo.Options
	db   *levigo.DB
	mo   store.MergeOperator

	mergeMutex sync.Mutex
}

func New(mo store.MergeOperator, config map[string]interface{}) (store.KVStore, error) {
	path, ok := config["path"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify path")
	}

	rv := Store{
		path: path,
		opts: levigo.NewOptions(),
		mo:   mo,
	}

	_, err := applyConfig(rv.opts, config)
	if err != nil {
		return nil, err
	}

	rv.db, err = levigo.Open(rv.path, rv.opts)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// func (ldbs *Store) get(key []byte) ([]byte, error) {
// 	options := defaultReadOptions()
// 	b, err := ldbs.db.Get(options, key)
// 	options.Close()
// 	return b, err
// }

// func (ldbs *Store) getWithSnapshot(key []byte, snapshot *levigo.Snapshot) ([]byte, error) {
// 	options := defaultReadOptions()
// 	options.SetSnapshot(snapshot)
// 	b, err := ldbs.db.Get(options, key)
// 	options.Close()
// 	return b, err
// }

// func (ldbs *Store) set(key, val []byte) error {
// 	ldbs.writer.Lock()
// 	defer ldbs.writer.Unlock()
// 	return ldbs.setlocked(key, val)
// }

// func (ldbs *Store) setlocked(key, val []byte) error {
// 	options := defaultWriteOptions()
// 	err := ldbs.db.Put(options, key, val)
// 	options.Close()
// 	return err
// }

// func (ldbs *Store) delete(key []byte) error {
// 	ldbs.writer.Lock()
// 	defer ldbs.writer.Unlock()
// 	return ldbs.deletelocked(key)
// }

// func (ldbs *Store) deletelocked(key []byte) error {
// 	options := defaultWriteOptions()
// 	err := ldbs.db.Delete(options, key)
// 	options.Close()
// 	return err
// }

func (s *Store) Close() error {
	s.db.Close()
	s.opts.Close()
	return nil
}

// func (ldbs *Store) iterator(key []byte) store.KVIterator {
// 	rv := newIterator(ldbs)
// 	rv.Seek(key)
// 	return rv
// }

func (s *Store) Reader() (store.KVReader, error) {
	return &Reader{
		store:    s,
		snapshot: s.db.NewSnapshot(),
	}, nil
}

func (s *Store) Writer() (store.KVWriter, error) {
	return &Writer{
		store: s,
	}, nil
}

func init() {
	registry.RegisterKVStore(Name, New)
}
