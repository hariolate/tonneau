package config

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net"
)

type Parsed struct {
	Config *Config

	DBConn      *gorm.DB
	RedisClient *redis.Client
	Listener    net.Listener

	Context context.Context
}

func (p *Parsed) setupListener() (err error) {
	addr, port := p.Config.ServeOn.Addr, p.Config.ServeOn.Port
	if addr == "" {
		if port == "" {
			p.Listener, err = net.Listen("tcp", "")
			return
		}
		p.Listener, err = net.Listen("tcp", ":"+port)
		return
	}
	p.Listener, err = net.Listen("tcp", addr+":"+port)
	return
}

func (p *Parsed) setupDB() error {
	var DBOpener func(dsn string) gorm.Dialector

	switch p.Config.Storage.DBType {
	case "mysql":
		DBOpener = mysql.Open
	case "postgres":
		DBOpener = postgres.Open
	case "sqlite":
		DBOpener = sqlite.Open
	default:
		return errors.New("unsupported db type")
	}

	db, err := gorm.Open(DBOpener(p.Config.Storage.DBDsn), &gorm.Config{})
	if err != nil {
		return err
	}
	p.DBConn = db
	return nil
}

func (p *Parsed) setupRedis() error {
	opt, err := redis.ParseURL(p.Config.Storage.RedisUrl)
	if err != nil {
		return err
	}
	p.RedisClient = redis.NewClient(opt)
	return nil
}

func (p *Parsed) setup() error {
	if err := p.setupListener(); err != nil {
		return err
	}
	log.Println("serving on", p.Listener.Addr())
	if err := p.setupDB(); err != nil {
		return err
	}
	if err := p.setupRedis(); err != nil {
		return err
	}
	return nil
}

func (c *Config) Parse() (*Parsed, error) {
	parsed := &Parsed{Config: c}
	if err := parsed.setup(); err != nil {
		return nil, err
	}
	return parsed, nil
}

func (c *Config) MustParse() *Parsed {
	parsed, err := c.Parse()
	if err != nil {
		log.Panicln(err)
	}
	return parsed
}
