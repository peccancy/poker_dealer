package mongo

import (
	"context"
	"github.com/globalsign/mgo"
	"time"

	"github.com/peccancy/chassi/mongo"
	"github.com/peccancy/poker_dealer/logs"
)

const (
	GameCollection = "game"
)

type (
	Config struct {
		Host       string
		Port       int
		Database   string
		User       string
		Password   string
		SSLEnabled bool
		Timeout    time.Duration
		Debug      bool
	}

	Repo struct {
		db *mongo.Mongo
	}
)

// Check mongo connection
func (r *Repo) HealthCheck() error {
	return r.db.Connection.Ping()
}

// Close mongo connection
func (r *Repo) Close() {
	if r.db.Connection != nil {
		r.db.Connection.Close()
	}
}

// Create mongo connection by configurations
func Connect(ctx context.Context, config *Config) (*mongo.Mongo, error) {
	var (
		conn *mongo.Mongo
		err  error
	)

	if config.SSLEnabled {
		conn, err = mongo.ConnectClientSecure(
			config.Host,
			config.Port,
			config.Database,
			"root",
			config.User,
			config.Password,
			"",
			config.Timeout,
		)
	} else {
		conn, err = mongo.ConnectClient(
			config.Host,
			config.Port,
			config.Database,
			"root",
			config.Timeout,
		)
	}

	if err != nil {
		return nil, err
	}

	if config.Debug {
		mgo.SetDebug(true)
		mgo.SetLogger(logs.MongoDB)
	}

	return conn, nil
}

// Create repository
func New(ctx context.Context, config *Config) (*Repo, error) {
	conn, err := Connect(ctx, config)
	if err != nil {
		return nil, err
	}

	err = createIndexes(conn)
	if err != nil {
		return nil, err
	}

	return &Repo{
		db: conn,
	}, nil
}

func createIndexes(conn *mongo.Mongo) error {
	collection := conn.Database.C(GameCollection)

	index := mgo.Index{
		Name:       "id",
		Key:        []string{"id"},
		Unique:     true,
		Background: true,
	}
	err := collection.EnsureIndex(index)
	if err != nil {
		return err
	}

	index = mgo.Index{
		Name:       "account_linked_entities",
		Key:        []string{"account_id", "linked_entities.key"},
		Background: true,
	}
	err = collection.EnsureIndex(index)
	if err != nil {
		return err
	}

	return nil
}
