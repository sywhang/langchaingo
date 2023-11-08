//nolint:forbidigo
package callbacks

import (
	"context"

	"github.com/sywhang/langchaingo/llms"
	"github.com/sywhang/langchaingo/schema"
)

type SimpleHandler struct{}

var _ Handler = SimpleHandler{}

func (SimpleHandler) HandleText(context.Context, string)                            {}
func (SimpleHandler) HandleLLMStart(context.Context, []string)                      {}
func (SimpleHandler) HandleLLMEnd(context.Context, llms.LLMResult)                  {}
func (SimpleHandler) HandleChainStart(context.Context, map[string]any)              {}
func (SimpleHandler) HandleChainEnd(context.Context, map[string]any)                {}
func (SimpleHandler) HandleToolStart(context.Context, string)                       {}
func (SimpleHandler) HandleToolEnd(context.Context, string)                         {}
func (SimpleHandler) HandleAgentAction(context.Context, schema.AgentAction)         {}
func (SimpleHandler) HandleRetrieverStart(context.Context, string)                  {}
func (SimpleHandler) HandleRetrieverEnd(context.Context, string, []schema.Document) {}
func (SimpleHandler) HandleStreamingFunc(context.Context, []byte)                   {}
