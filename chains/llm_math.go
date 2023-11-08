package chains

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/sywhang/langchaingo/llms"
	"github.com/sywhang/langchaingo/memory"
	"github.com/sywhang/langchaingo/prompts"
	"github.com/sywhang/langchaingo/schema"
	"go.starlark.net/lib/math"
	"go.starlark.net/starlark"
)

const (
	quote          = "```"
	_llmMathPrompt = `Translate a math problem into a expression that can be evaluated as Starlark.
Use the output of running this code to answer the question.

---
Question: (Question with math problem.)
` + quote + `starlark
$(single line expression that solves the problem)
` + quote + `

---
Question: What is 37593 * 67?
` + quote + `starlark
37593 * 67
` + quote + `

---
Question: {{.question}}
`
)

// LLMMathChain is a chain used for evaluating math expressions.
type LLMMathChain struct {
	LLMChain *LLMChain
}

var _ Chain = LLMMathChain{}

func NewLLMMathChain(llm llms.LanguageModel) LLMMathChain {
	p := prompts.NewPromptTemplate(_llmMathPrompt, []string{"question"})
	c := NewLLMChain(llm, p)
	return LLMMathChain{
		LLMChain: c,
	}
}

// Call gets relevant documents from the retriever and gives them to the combine
// documents chain.
func (c LLMMathChain) Call(ctx context.Context, values map[string]any, options ...ChainCallOption) (map[string]any, error) { //nolint: lll
	question, ok := values["question"].(string)
	if !ok {
		return nil, fmt.Errorf("%w: %w", ErrInvalidInputValues, ErrInputValuesWrongType)
	}
	output, err := Call(ctx, c.LLMChain, map[string]any{
		"question": question,
	}, options...)
	if err != nil {
		return nil, err
	}
	output["answer"], err = c.processLLMResult(output["text"].(string))
	return output, err
}

func (c LLMMathChain) GetMemory() schema.Memory { //nolint:ireturn
	return memory.NewSimple()
}

func (c LLMMathChain) GetInputKeys() []string {
	return []string{"question"}
}

func (c LLMMathChain) GetOutputKeys() []string {
	return []string{"answer"}
}

var starlarkBlockRegex = regexp.MustCompile("(?s)```starlark(.*)```")

func (c LLMMathChain) processLLMResult(llmOutput string) (string, error) {
	llmOutput = strings.TrimSpace(llmOutput)
	textMatch := starlarkBlockRegex.FindStringSubmatch(llmOutput)
	if len(textMatch) > 0 {
		expression := textMatch[1]
		output, err := c.evaluateExpression(expression)
		if err != nil {
			return "", err
		}
		return output, nil
	}
	if strings.Contains(llmOutput, "Answer:") {
		return strings.TrimSpace(strings.Split(llmOutput, "Answer:")[1]), nil
	}
	return "", fmt.Errorf("unknown format from LLM: %s", llmOutput) //nolint:goerr113
}

func (c LLMMathChain) evaluateExpression(expression string) (string, error) {
	expression = strings.TrimSpace(expression)
	v, err := starlark.Eval(&starlark.Thread{Name: "main"}, "input", expression, math.Module.Members)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}
