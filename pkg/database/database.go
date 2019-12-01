package database

import (
	"github.com/jackc/pgx"
)

var Connection *pgx.ConnPool

func Init() {
	Connection, _ = pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host: "localhost",
				Port: 5432,
				//Database: "postgres",
				Database: "forum",
				User:     "forum",
				Password: "forum",
			},
			MaxConnections: 50,
		})
}
