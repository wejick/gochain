package summarization

import (
	"context"
	"errors"

	"github.com/wejick/gochain/callback"
	"github.com/wejick/gochain/chain"
	"github.com/wejick/gochain/chain/combine_document"
	"github.com/wejick/gochain/chain/llm_chain"
	"github.com/wejick/gochain/model"
	"github.com/wejick/gochain/prompt"
)

const (
	promptSummarizeStuff = `Write a concise summary of the following:
"{{.text}}"
CONCISE SUMMARY:`
)

type StuffSummarizationChain struct {
	stuffCombineDocument *combine_document.StuffCombineDocument
	callbackManager      *callback.Manager
}

var _ chain.BaseChain = &StuffSummarizationChain{}

func NewStuffSummarizationChain(llm_chain *llm_chain.LLMChain, callbackManager *callback.Manager,
	promptTemplateString string, promptTemplateKey string, verbose bool) (s *StuffSummarizationChain, err error) {

	if verbose {
		callbackManager.RegisterCallback(chain.CallbackChainEnd, callback.VerboseCallback)
	}

	var promptTemplate *prompt.PromptTemplate

	if promptTemplateString == "" {
		promptTemplate, err = prompt.NewPromptTemplate("stuff", promptSummarizeStuff)
		if err != nil {
			return
		}
		promptTemplateKey = "text"
	}

	stuffCombineDocument := combine_document.NewStuffCombineDocument(callbackManager, promptTemplate, promptTemplateKey, llm_chain, verbose)
	s = &StuffSummarizationChain{
		stuffCombineDocument: stuffCombineDocument,
		callbackManager:      callbackManager,
	}

	return
}

// Run all entries in input map will be treated as document to be combined
// output will be output["output"]
func (S *StuffSummarizationChain) Run(ctx context.Context, input map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if _, ok := input["input"]; !ok {
		return output, errors.New("input[\"input\"] is not specified")
	}
	//trigger callback chain start
	S.callbackManager.TriggerEvent(ctx, chain.CallbackChainStart, callback.CallbackData{
		FunctionName: "StuffSummarizationChain.Run",
		Input:        input,
		Output:       output,
	})

	output, err = S.stuffCombineDocument.Run(ctx, input, options...)

	// trigger callback chain end
	S.callbackManager.TriggerEvent(ctx, chain.CallbackChainEnd, callback.CallbackData{
		FunctionName: "StuffSummarizationChain.Run",
		Input:        input,
		Output:       output,
	})
	return
}

// SimpleRun will run the input prompt string againts llmchain
func (S *StuffSummarizationChain) SimpleRun(ctx context.Context, input string, options ...func(*model.Option)) (output string, err error) {
	//trigger callback chain start
	S.callbackManager.TriggerEvent(ctx, chain.CallbackChainStart, callback.CallbackData{
		FunctionName: "StuffSummarizationChain.SimpleRun",
		Input:        map[string]string{"input": input},
		Output:       map[string]string{"output": output},
	})
	output, err = S.stuffCombineDocument.SimpleRun(ctx, input, options...)
	// trigger callback chain end
	S.callbackManager.TriggerEvent(ctx, chain.CallbackChainEnd, callback.CallbackData{
		FunctionName: "StuffSummarizationChain.SimpleRun",
		Input:        map[string]string{"input": input},
		Output:       map[string]string{"output": output},
	})
	return
}
