package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/rs/zerolog"

	"github.com/cosmos/iavl"
	idbm "github.com/cosmos/iavl/db"
)

func main() {
	db, err := OpenDB("test.db")
	if err != nil {
		panic(err)
	}

	tree := iavl.NewMutableTree(idbm.NewWrapper(db), 5, false, log.NewLogger(os.Stdout, log.LevelOption(zerolog.InfoLevel)))
	_, err = tree.LoadVersion(int64(0))
	if err != nil {
		panic(err)
	}

	m := generateMap(1000)
	for k, v := range m {
		tree.Set([]byte(k), []byte(v))
	}

	// Possible fix: sort the key-value pairs before inserting them
	//
	// insertions := make([]KV, 0, len(m))
	// for k, v := range m {
	// 	insertions = append(insertions, KV{[]byte(k), []byte(v)})
	// }
	// sort.Slice(insertions, func(i, j int) bool {
	// 	return bytes.Compare(insertions[i].Key, insertions[j].Key) < 0
	// })
	// for _, kv := range insertions {
	// 	tree.Set(kv.Key, kv.Value)
	// }

	// Save the version
	saveVersion(tree)
}

type KV struct {
	Key   []byte
	Value []byte
}

func generateMap(size int) map[string]string {
	m := make(map[string]string)
	for i := 0; i < size; i++ {
		m[strconv.Itoa(i)] = strconv.Itoa(i)
	}
	return m
}

func saveVersion(tree *iavl.MutableTree) {
	hash, newV, err := tree.SaveVersion()
	if err != nil {
		panic(err)
	}
	fmt.Println("Hash", hex.EncodeToString(hash), "Version", newV)
}

func OpenDB(dir string) (dbm.DB, error) {
	switch {
	case strings.HasSuffix(dir, ".db"):
		dir = dir[:len(dir)-3]
	case strings.HasSuffix(dir, ".db/"):
		dir = dir[:len(dir)-4]
	default:
		return nil, fmt.Errorf("database directory must end with .db")
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	// TODO: doesn't work on windows!
	cut := strings.LastIndex(dir, "/")
	if cut == -1 {
		return nil, fmt.Errorf("cannot cut paths on %s", dir)
	}
	name := dir[cut+1:]
	db, err := dbm.NewGoLevelDB(name, dir[:cut], nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}
