package rpcaiagent

import (
	"testing"

	mysqlrepo "github.com/dingzijian9527-del/Travel-Assistant/pkg/repository/mysql"
)

func TestMySQLAIAgentMessageRepositorySatisfiesServiceInterface(t *testing.T) {
	var _ aiAgentRepository = (*mysqlrepo.AIAgentMessageRepository)(nil)
}
