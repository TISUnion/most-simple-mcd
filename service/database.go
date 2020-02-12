package service

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/utils"
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
	"path/filepath"
	"time"
)

var (
	dataBases = make(map[string]_interface.Database)
)

type DataBase struct {
	badgerDb *badger.DB
	name     string
}

func (d *DataBase) InitCallBack() {
}

func (d *DataBase) ChangeConfCallBack() {
}

func (d *DataBase) DestructCallBack() {
	if d.badgerDb != nil {
		_ = d.badgerDb.Close()
	}
}

func (d *DataBase) Get(k string) string {
	var result string
	if d.badgerDb == nil {
		return result
	}
	_ = d.badgerDb.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		if err != nil {
			return err
		}
		resultByte, err := item.ValueCopy(nil)
		if len(resultByte) > 0 {
			result = string(resultByte)
		}
		if err != nil {
			return err
		}
		return nil
	})
	return result
}

func (d *DataBase) Set(k string, v string) {
	if d.badgerDb == nil {
		return
	}
	_ = d.badgerDb.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(k), []byte(v))
		return err
	})
}

func (d *DataBase) SetWiteTTL(k string, v string, t time.Duration) {
	if d.badgerDb == nil {
		return
	}
	_ = d.badgerDb.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(k), []byte(v)).WithTTL(t)
		err := txn.SetEntry(e)
		return err
	})
}

func GetDataBaseInstance(dbname string) _interface.Database {
	name := string(dbname)
	if db, ok := dataBases[name]; ok {
		return db
	}
	db := newDataBase(name)
	if db == nil {
		return nil
	}
	// 注册回调
	RegisterCallBack(db)
	dataBases[name] = db
	return db
}

func GetFromDefault(key string) string{
	db := GetDataBaseInstance(constant.DEFAULT_DATABASE_NAME)
	return db.Get(key)
}

func SetFromDefault(key string, value string) {
	db := GetDataBaseInstance(constant.DEFAULT_DATABASE_NAME)
	db.Set(key, value)
}

func newDataBase(name string) _interface.Database {

	if badgerDb, err := badger.Open(newDataBaseOptions(name)); err != nil {
		fmt.Println(err)
		return nil
	} else {
		return &DataBase{badgerDb, name}
	}

}

func newDataBaseOptions(name string) badger.Options {
	path := filepath.Join(GetConfVal(constant.WORKSPACE), "db-resource", name)
	_ = utils.CreatDir(path)
	return badger.Options{
		Dir:                     path,
		ValueDir:                path,
		LevelOneSize:            256 << 20,
		LevelSizeMultiplier:     10,
		TableLoadingMode:        options.MemoryMap,
		ValueLogLoadingMode:     options.MemoryMap,
		MaxLevels:               7,
		MaxTableSize:            64 << 20,
		NumCompactors:           2,
		NumLevelZeroTables:      5,
		NumLevelZeroTablesStall: 10,
		NumMemtables:            5,
		SyncWrites:              false,
		NumVersionsToKeep:       1,
		CompactL0OnClose:        true,
		ValueLogFileSize:        1<<30 - 1,

		ValueLogMaxEntries: 1000000,
		ValueThreshold:     32,
		Truncate:           false,
		Logger:             GetLogContainerInstance().AddLog("db-"+name, constant.LOG_DEBUG),
		LogRotatesToFlush:  2,
	}
}
