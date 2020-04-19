package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// DBConn is the Database interface
type DBConn interface {
	Init(conf *DBInitConfig) error
	Get(keyChain string, value interface{}) error
	GetRaw(keyChain string) ([]byte, error)
	GetType(keyChain string) string
	Set(keyChain string, value interface{}) error
	Del(keyChain string) error
	CreateBucket(keyChain string) error
	CreateBucketIfNotExists(keyChain string) error
	ReCreateBucket(keyChain string) error
	DeleteBucket(keyChain string) error
	RecordKeys(bucketChain string) ([]string, error)
	BucketKeys(bucketChain string) ([]string, error)
	AllKeys(bucketChain string) ([]string, error)
	Tree(bucketChain string) (map[string]interface{}, error)
	Close() error
}

// DataBase is the Database configuration object
type DataBase struct {
	Type string `json:"type" yaml:"type" mapstructure:"type"`
	URL  string `json:"url"  yaml:"url"  mapstructure:"url"`
}

// Validate the DataBase configurations
func (db DataBase) Validate() error {
	return validation.ValidateStruct(&db,
		validation.Field(&db.Type, validation.Required, validation.In("memory", "bolt")),
		validation.Field(&db.URL, validation.Required),
	)
}

// DatabaseStructure has a list of Main-buckets
type DatabaseStructure struct {
	MainBuckets []string `json:"main_buckets" yaml:"main_buckets" mapstructure:"main_buckets"`
}

// DBInitConfig is the object we will initialise our database with
// It consists of the user defined DBConfigs
// and the our DBStructure
type DBInitConfig struct {
	DBConfig    *DataBase          `json:"database"     yaml:"database"     mapstructure:"database"`
	DBStructure *DatabaseStructure `json:"db_structure" yaml:"db_structure" mapstructure:"db_structure"`
}

// GetDBStructure returns the Bounce 'n Junk server specific database structure: configs, capacity, usage
func GetDBStructure() *DatabaseStructure {
	return &DatabaseStructure{
		MainBuckets: []string{
			"users",
			"matches",
			"levels",
		},
	}
}

// GetDBInitConfig generates the DBInitConfig object, from our "defaults" and user settings
func GetDBInitConfig(dbConf *DataBase) *DBInitConfig {
	return &DBInitConfig{
		DBConfig:    dbConf,
		DBStructure: GetDBStructure(),
	}
}
