package repository

import (
	"database/sql"
	"dns-proxy/pkg/domain/proxy"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/dns/dnsmessage"
)

type repository struct {
	db *sql.DB
}

func NewRepository(path string) (proxy.Repository, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &repository{db}, nil
}

func (r *repository) SaveQuery(dnsm *dnsmessage.Message, action string) error {
	for _, q := range dnsm.Questions {
		domain := q.Name.String()
		queryType := q.Type.String()
		class := q.Class.String()
		statement, err := r.db.Prepare("INSERT INTO resolved(domain, type, class, action, date) values(?,?,?,?,?)")
		if err != nil {
			return fmt.Errorf("could not prepare statement: %v", err)
		}

		_, err = statement.Exec(domain, queryType, class, action, time.Now())
		if err != nil {
			return fmt.Errorf("could not insert register: %v", err)
		}
	}
	return nil
}
