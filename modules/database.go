package modules

import (
	"encoding/json"
	"github.com/TISUnion/most-simple-mcd/constant"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/utils"
	"github.com/boltdb/bolt"
	"path/filepath"
	"time"
)

type DataBase struct {
	boltDb *bolt.DB
	name   string
}

var dbMap = make(map[string]_interface.Database)

func (b *DataBase) close() {
	_ = b.boltDb.Close()
	delete(dbMap, b.name)
}

func (b *DataBase) ChangeConfCallBack() {
}

func (b *DataBase) DestructCallBack() {
	b.close()
}

func (b *DataBase) InitCallBack() {}

func (b *DataBase) Get(key string) (val string) {
	_ = b.boltDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.name))
		v := b.Get([]byte(key))
		valObj := &models.ValueWithTTL{}
		_ = json.Unmarshal(v, valObj)
		if valObj.NanoExpire >= time.Now().UnixNano() || valObj.NanoExpire == constant.PERMANENT_VAL_TTL {
			val = valObj.Value
		}
		return nil
	})
	return
}

func (b *DataBase) Set(key string, val string) {
	_ = b.boltDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.name))
		err := b.Put([]byte(key), setValJson(val, constant.PERMANENT_VAL_TTL))
		return err
	})
}

func (b *DataBase) SetWiteTTL(key string, val string, ttl time.Duration) {
	_ = b.boltDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.name))
		err := b.Put([]byte(key), setValJson(val, ttl))
		return err
	})
}

func GetDatabase(name string) (db _interface.Database, err error) {
	if datab, ok := dbMap[name]; ok {
		db = datab
		return
	}
	path := filepath.Join(GetConfVal(constant.WORKSPACE), "db-resource", name, "datum.db")
	_, _ = utils.CreateFile(path)
	boltdb, err := bolt.Open(path, 0664, nil)
	if err != nil {
		return
	}
	db = &DataBase{
		boltDb: boltdb,
		name:   name,
	}
	dbMap[name] = db
	return
}

func setValJson(val string, ttl time.Duration) []byte {
	t := time.Now().Add(ttl).UnixNano()
	obj := &models.ValueWithTTL{
		Value:      val,
		NanoExpire: t,
	}
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return []byte{}
	}
	return jsonStr
}

func GetFromDatabase(key string) string {
	db, err := GetDatabase(constant.DEFAULT_DATABASE_NAME)
	if err != nil {
		return ""
	}
	return db.Get(key)
}
func SetFromDatabase(key string, value string) {
	db, err := GetDatabase(constant.DEFAULT_DATABASE_NAME)
	if err != nil {
		return
	}
	db.Set(key, value)
}
func SetWiteTTLFromDatabase(key string, value string, d time.Duration) {
	db, err := GetDatabase(constant.DEFAULT_DATABASE_NAME)
	if err != nil {
		return
	}
	db.SetWiteTTL(key, value, d)
}