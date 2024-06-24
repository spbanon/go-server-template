package dbhandler

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDBConnection(host string, port string, username string, password string, dbname string, readOnly bool) (*pgxpool.Pool, error) {
	conninfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable pool_max_conn_idle_time=360s", username, password, host, port, dbname)
	conn, err := pgxpool.Connect(context.Background(), conninfo)
	if err != nil {
		return nil, err
	}
	return conn, err
}
