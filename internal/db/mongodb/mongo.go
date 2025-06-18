package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type MongoDBClient struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoDBClient(config MongoDBConfig) (*MongoDBClient, error) {
	log.Printf("[INFO] Connecting to MongoDB at %s:%d", config.Host, config.Port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", config.User, config.Password, config.Host, config.Port)
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("[ERROR] MongoDB connection ping failed: %v", err)
		return nil, err
	}

	log.Printf("[INFO] Successfully connected to MongoDB")

	db := client.Database(config.Database)

	mongoClient := &MongoDBClient{
		client: client,
		db:     db,
	}

	if err := mongoClient.initCollections(); err != nil {
		log.Printf("[ERROR] Failed to initialize collections: %v", err)
		return nil, err
	}

	return mongoClient, nil
}

func (c *MongoDBClient) initCollections() error {
	log.Println("[DEBUG] Initializing MongoDB collections")
	
	// 创建邮件集合
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collections, err := c.db.ListCollectionNames(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to list collections: %v", err)
	}

	// 检查邮件集合是否存在
	if !contains(collections, "emails") {
		log.Println("[INFO] Creating emails collection")
		if err := c.db.CreateCollection(ctx, "emails"); err != nil {
			return fmt.Errorf("failed to create emails collection: %v", err)
		}
	}

	log.Println("[INFO] MongoDB collections initialized successfully")
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (c *MongoDBClient) GetDB() *mongo.Database {
	return c.db
}

func (c *MongoDBClient) Close() error {
	log.Println("[INFO] Closing MongoDB connection")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.client.Disconnect(ctx)
}