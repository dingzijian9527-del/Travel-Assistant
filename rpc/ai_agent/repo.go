package rpcaiagent

import "context"

type aiAgentRepository interface {
	AppendMessages(ctx context.Context, userID int64, messages ...*chatMessageModel)
	History(ctx context.Context, userID int64) []chatMessageModel
}
