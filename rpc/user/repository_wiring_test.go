package rpcuser

import (
	"testing"

	mysqlrepo "github.com/dingzijian9527-del/Travel-Assistant/pkg/repository/mysql"
)

func TestMySQLUserRepositorySatisfiesServiceInterface(t *testing.T) {
	var _ userRepository = (*mysqlrepo.UserRepository)(nil)
}
