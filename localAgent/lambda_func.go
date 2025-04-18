package localAgent

import (
	"context"
	"time"
)

// newLambda component initialization function of node 'InputToQuery' in graph 'localAgent'
func newLambda(ctx context.Context, input *UserMessage, opts ...any) (output string, err error) {
	return input.Query, nil
}

// newLambda2 component initialization function of node 'InputToHistory' in graph 'localAgent'
func newLambda2(ctx context.Context, input *UserMessage, opts ...any) (output map[string]any, err error) {
	return map[string]any{
		"content": input.Query,
		"history": input.History,
		"date":    time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
