package database

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/donbattery/bnj/model"
)

// BoltDB is an implementation of the DB interface
// it use the Bolt to maintain a single database file
// and store our buckets there
type BoltDB struct {
	db *bolt.DB
}

// New returns a new BoltDB instance
func New() *BoltDB {
	return &BoltDB{}
}

// Init sets up the database,
// it tries to load the database file passed as config, if it not exists
// it tries to create the file. Then it creates the the root bucket
// and the configs, capacity and usage main-buckets if they are not already exists.
func (bdb *BoltDB) Init(conf *model.DBInitConfig) error {
	if conf == nil {
		return errors.New("nil configuration provided")
	}
	db, err := bolt.Open(conf.DBConfig.URL, 0600, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to open db file")
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("root"))
		if err != nil {
			return errors.Wrap(err, "Cannot create root bucket")
		}
		for _, mainBucket := range conf.DBStructure.MainBuckets {
			_, err = root.CreateBucketIfNotExists([]byte(mainBucket))
			if err != nil {
				return errors.Wrapf(err, "Cannot create Main-bucket: %s", mainBucket)
			}
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "Failed to set up buckets")
	}
	bdb.db = db
	return nil
}

// parseKeyChain is a helper function to turn a keyChain (e.g.: configs.remote.master)
// into BoltDB bucket names and db-key
func parseKeyChain(keyChain string) ([][]byte, []byte, error) {
	var buckets [][]byte
	var key []byte
	// Separate the keyChain elements by .
	chain := strings.Split(keyChain, ".")
	// If the chain has no elements return error
	if len(chain) < 1 {
		return buckets, key, errors.Errorf("Invalid keyChain %s", keyChain)
	}
	// Let key be the last element of the chain
	key = []byte(chain[len(chain)-1])
	// Range over the elements except the last and add them as buckets
	for _, bucketName := range chain[:len(chain)-1] {
		buckets = append(buckets, []byte(bucketName))
	}
	return buckets, key, nil
}

// Get a key from the databas by the keyChain e.g.: configs.remote.master.macstadium-atl-pub-1
// where configs is a Main-bucket, remote and master are a subuckets, and is the actual db-key
// is macstadium-atl-pub-1. The underlying BoltDB byte array will be unmarshalled into the
// supplyed value interface
func (bdb *BoltDB) Get(keyChain string, value interface{}) error {
	// Parse the keyChain into buckets and key elements
	buckets, key, err := parseKeyChain(keyChain)
	if err != nil {
		return errors.Wrap(err, "Failed to parse keyChain")
	}
	// Open the database transaction
	if err := bdb.db.View(func(tx *bolt.Tx) error {
		// Get the root bucket
		root := tx.Bucket([]byte("root"))
		// This will be our first targetBucket
		targetBucket := root
		// Range over Sub-buckets and change target accordingly
		for _, bucket := range buckets {
			targetBucket = targetBucket.Bucket(bucket)
			if targetBucket == nil {
				return errors.Errorf("Cannot get: %s Sub-bucket does not exists: %s", keyChain, bucket)
			}
		}
		// Get the value by the last element of the keyChain
		valueBytes := targetBucket.Get(key)
		if err := json.Unmarshal(valueBytes, &value); err != nil {
			return errors.Wrapf(err, "Failed to decode db value into go object: %s", valueBytes)
		}
		return nil
	}); err != nil {
		return errors.Wrapf(err, "Failed to get: %s", keyChain)
	}
	return nil
}

// Get the DBType by the keyChain: Undefined, Bucket or Key
func (bdb *BoltDB) GetType(keyChain string) string {
	// Empty string represents the root, and it is a bucket
	if keyChain == "" {
		return "Bucket"
	}

	dbType := "Undefined"

	buckets, key, err := parseKeyChain(keyChain)
	if err != nil {
		return dbType
	}

	// we do not return any error from the View,
	// just let the dbType remain undefined
	_ = bdb.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		targetBucket := root
		for _, bucket := range buckets {
			targetBucket = targetBucket.Bucket(bucket)
			if targetBucket == nil {
				return nil
			}
		}
		if subBucket := targetBucket.Bucket(key); subBucket != nil {
			dbType = "Bucket"
		} else {
			if val := targetBucket.Get(key); val != nil {
				dbType = "Key"
			}
		}
		return nil
	})

	return dbType
}

func (bdb *BoltDB) GetRaw(keyChain string) ([]byte, error) {
	var out []byte
	buckets, key, err := parseKeyChain(keyChain)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse keyChain")
	}
	if err := bdb.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		targetBucket := root
		for _, bucket := range buckets {
			targetBucket = targetBucket.Bucket(bucket)
			if targetBucket == nil {
				return errors.Errorf("Cannot get: %s Sub-bucket does not exists: %s", keyChain, bucket)
			}
		}
		out = targetBucket.Get(key)
		return nil
	}); err != nil {
		return nil, errors.Wrapf(err, "Failed to get RAW value of: %s", keyChain)
	}
	return out, nil
}

// Set a value by the keyChain
func (bdb *BoltDB) Set(keyChain string, value interface{}) error {
	buckets, key, err := parseKeyChain(keyChain)
	if err != nil {
		return errors.Wrap(err, "Failed to parse keyChain")
	}
	return bdb.db.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		targetBucket := root
		for _, bucket := range buckets {
			targetBucket = targetBucket.Bucket(bucket)
			if targetBucket == nil {
				return errors.Errorf("Cannot set: %s Sub-bucket does not exists: %s", keyChain, bucket)
			}
		}
		valueBytes, err := json.Marshal(value)
		if err != nil {
			return errors.Wrapf(err, "Failed to encode go object %+v into database value", value)
		}
		if err := targetBucket.Put(key, valueBytes); err != nil {
			return errors.Wrapf(err, "Cannot set value %s to %s", valueBytes, keyChain)
		}
		return nil
	})
}

// Delete a value by the keyChain
func (bdb *BoltDB) Del(keyChain string) error {
	buckets, key, err := parseKeyChain(keyChain)
	if err != nil {
		return errors.Wrap(err, "Failed to parse keyChain")
	}
	return bdb.db.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		targetBucket := root
		for _, bucket := range buckets {
			targetBucket = targetBucket.Bucket(bucket)
			if targetBucket == nil {
				return errors.Errorf("Cannot delete: %s Sub-bucket does not exists: %s", keyChain, bucket)
			}
		}
		if err := targetBucket.Delete(key); err != nil {
			return errors.Wrap(err, "Cannot delete value")
		}
		return nil
	})
}

func (bdb *BoltDB) CreateBucket(keyChain string) error {
	buckets, key, err := parseKeyChain(keyChain)
	if err != nil {
		return errors.Wrap(err, "Failed to parse keyChain")
	}
	return bdb.db.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		targetBucket := root
		for _, bucket := range buckets {
			targetBucket = targetBucket.Bucket(bucket)
			if targetBucket == nil {
				return errors.Errorf("Cannot create new bucket: %s Sub-bucket does not exists: %s", keyChain, bucket)
			}
		}
		if _, err := targetBucket.CreateBucket(key); err != nil {
			return errors.Wrap(err, "Cannot create new bucket")
		}
		return nil
	})
}

func (bdb *BoltDB) ReCreateBucket(keyChain string) error {
	buckets, key, err := parseKeyChain(keyChain)
	if err != nil {
		return errors.Wrap(err, "Failed to parse keyChain")
	}
	return bdb.db.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		targetBucket := root
		for _, bucket := range buckets { // find the parent of the bucket
			targetBucket = targetBucket.Bucket(bucket)
			if targetBucket == nil {
				return errors.Errorf("Cannot create new bucket: %s Sub-bucket does not exists: %s", keyChain, bucket)
			}
		}
		if oldBucket := targetBucket.Bucket(key); oldBucket != nil { // get the old bucket
			if err := targetBucket.DeleteBucket(key); err != nil { // if it exists delete it
				return errors.Wrap(err, "Cannot delete bucket")
			}
		}
		if _, err := targetBucket.CreateBucket(key); err != nil { // create the bucket
			return errors.Wrap(err, "Cannot create new bucket")
		}
		return nil
	})
}

func (bdb *BoltDB) CreateBucketIfNotExists(keyChain string) error {
	buckets, key, err := parseKeyChain(keyChain)
	if err != nil {
		return errors.Wrap(err, "Failed to parse keyChain")
	}
	return bdb.db.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		targetBucket := root
		for _, bucket := range buckets {
			targetBucket = targetBucket.Bucket(bucket)
			if targetBucket == nil {
				return errors.Errorf("Cannot create new bucket: %s Sub-bucket does not exists: %s", keyChain, bucket)
			}
		}
		_, err := targetBucket.CreateBucketIfNotExists(key)
		if err != nil {
			return errors.Wrap(err, "Cannot create bucket (if not exists)")
		}
		return nil
	})
}

func (bdb *BoltDB) DeleteBucket(keyChain string) error {
	buckets, key, err := parseKeyChain(keyChain)
	if err != nil {
		return errors.Wrap(err, "Failed to parse keyChain")
	}
	return bdb.db.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		targetBucket := root
		for _, bucket := range buckets {
			targetBucket = targetBucket.Bucket(bucket)
			if targetBucket == nil {
				return errors.Errorf("Cannot create delete bucket: %s Sub-bucket does not exists: %s", keyChain, bucket)
			}
		}
		if err := targetBucket.DeleteBucket(key); err != nil {
			return errors.Wrap(err, "Cannot delete bucket")
		}
		return nil
	})
}

// parseBucketChain is a helper function to turn a bucketChain (e.g.: configs.remote.master)
// into a BoltDB bucket name
func parseBucketChain(bucketChain string) []string {
	buckets := strings.Split(bucketChain, ".")
	if len(buckets) == 1 && buckets[0] == "" { // if the chain has only one element and it is an empty string
		buckets = buckets[:0] // turn buckets an empty slice, essencially targeting the root bucket
	}
	return buckets
}

func (bdb *BoltDB) RecordKeys(bucketChain string) ([]string, error) {
	res := []string{}
	buckets := parseBucketChain(bucketChain)
	if err := bdb.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		targetBucket := root
		for _, bucket := range buckets {
			targetBucket = targetBucket.Bucket([]byte(bucket))
			if targetBucket == nil {
				return errors.Errorf("Cannot get: %s Sub-bucket does not exists: %s", bucketChain, bucket)
			}
		}
		cursor := targetBucket.Cursor()
		for key, val := cursor.First(); key != nil; key, val = cursor.Next() {
			if val != nil {
				res = append(res, string(key))
			}
		}
		return nil
	}); err != nil {
		return res, errors.Wrapf(err, "Failed to get: %s", bucketChain)
	}
	return res, nil
}

func (bdb *BoltDB) BucketKeys(bucketChain string) ([]string, error) {
	res := []string{}
	buckets := parseBucketChain(bucketChain)
	if err := bdb.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		targetBucket := root
		for _, bucket := range buckets {
			targetBucket = targetBucket.Bucket([]byte(bucket))
			if targetBucket == nil {
				return errors.Errorf("Cannot get: %s Sub-bucket does not exists: %s", bucketChain, bucket)
			}
		}
		cursor := targetBucket.Cursor()
		for key, val := cursor.First(); key != nil; key, val = cursor.Next() {
			if val == nil {
				res = append(res, string(key))
			}
		}
		return nil
	}); err != nil {
		return res, errors.Wrapf(err, "Failed to get: %s", bucketChain)
	}
	return res, nil
}

func (bdb *BoltDB) AllKeys(bucketChain string) ([]string, error) {
	recordKeys, err := bdb.RecordKeys(bucketChain)
	if err != nil {
		return nil, err
	}

	bucketKeys, err := bdb.BucketKeys(bucketChain)
	if err != nil {
		return nil, err
	}

	return append(recordKeys, bucketKeys...), nil
}

// Tree build the tree view of the database from an antry bucket (if it is empty, the root bucket will be the entry point)
func (bdb *BoltDB) Tree(bucketChain string) (map[string]interface{}, error) {
	var res = make(map[string]interface{})
	buckets := parseBucketChain(bucketChain)
	if err := bdb.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		targetBucket := root
		for _, bucket := range buckets {
			targetBucket = targetBucket.Bucket([]byte(bucket))
			if targetBucket == nil {
				return errors.Errorf("Cannot get: %s Sub-bucket does not exists: %s", bucketChain, bucket)
			}
		}
		cursor := targetBucket.Cursor()
		for key, val := cursor.First(); key != nil; key, val = cursor.Next() {
			if val == nil {
				subChain := append(buckets, string(key))
				subKeys, err := bdb.Tree(strings.Join(subChain, "."))
				if err != nil {
					return errors.Wrap(err, "Failed to traverse tree")
				}
				res[string(key)] = subKeys
			} else {
				res[string(key)] = fmt.Sprintf("%d Byte", len(val))
			}
		}
		return nil
	}); err != nil {
		return res, errors.Wrapf(err, "Failed to get: %s", bucketChain)
	}
	return res, nil
}

// Close waits all DB transaction to be finished and closes the database
func (bdb *BoltDB) Close() error {
	return bdb.db.Close()
}
