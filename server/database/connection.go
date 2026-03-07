package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

var DB *gorm.DB

func (c *Connection) loadFromEnv() {
	c.Host = os.Getenv("PG_HOST")
	c.Port = os.Getenv("PG_PORT")
	c.User = os.Getenv("PG_USER")
	c.Password = os.Getenv("PG_PASSWORD")
	c.Database = os.Getenv("PG_DATABASE")
}

func (c *Connection) validateConfig() error {
	if c.Host == "" || c.Port == "" || c.User == "" || c.Database == "" {
		return errors.New("missing required PG env vars: PG_HOST, PG_PORT, PG_USER, PG_DATABASE")
	}
	return nil
}

func (c *Connection) URL() (string, error) {
	c.loadFromEnv()
	if err := c.validateConfig(); err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	), nil
}

func (c *Connection) InitGorm(ctx context.Context) (*gorm.DB, error) {
	connectionURL, err := c.URL()
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(postgres.Open(connectionURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func Init(db *gorm.DB) {
	DB = db
}
