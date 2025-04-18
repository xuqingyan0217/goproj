package localAgent

import (
	"context"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildLocalAgent(ctx context.Context) (r compose.Runnable[*UserMessage, *schema.Message], err error) {
	const (
		InputToQuery   = "InputToQuery"
		RedisRetriever = "RedisRetriever"
		ChatTemplate   = "ChatTemplate"
		ReactAgent     = "ReactAgent"
		InputToHistory = "InputToHistory"
	)
	g := compose.NewGraph[*UserMessage, *schema.Message]()
	_ = g.AddLambdaNode(InputToQuery, compose.InvokableLambdaWithOption(newLambda), compose.WithNodeName("UserMessageToQuery"))
	redisRetrieverKeyOfRetriever, err := newRetriever(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(RedisRetriever, redisRetrieverKeyOfRetriever, compose.WithOutputKey("documents"))
	chatTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate, chatTemplateKeyOfChatTemplate)
	lambda2KeyOfLambda, err := newLambda1(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(ReactAgent, lambda2KeyOfLambda, compose.WithNodeName("ReAct Agent"))
	_ = g.AddLambdaNode(InputToHistory, compose.InvokableLambdaWithOption(newLambda2), compose.WithNodeName("UserMessageToVariables"))
	_ = g.AddEdge(compose.START, InputToQuery)
	_ = g.AddEdge(compose.START, InputToHistory)
	_ = g.AddEdge(ReactAgent, compose.END)
	_ = g.AddEdge(InputToQuery, RedisRetriever)
	_ = g.AddEdge(RedisRetriever, ChatTemplate)
	_ = g.AddEdge(InputToHistory, ChatTemplate)
	_ = g.AddEdge(ChatTemplate, ReactAgent)
	r, err = g.Compile(ctx, compose.WithGraphName("localAgent"), compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
