package rpctrip

import (
	"testing"

	mysqlrepo "github.com/dingzijian9527-del/Travel-Assistant/pkg/repository/mysql"
)

func TestMySQLTripRepositorySatisfiesServiceInterface(t *testing.T) {
	var _ tripRepository = (*mysqlrepo.TripRepository)(nil)
}
