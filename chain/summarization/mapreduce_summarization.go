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
	"github.com/wejick/gochain/textsplitter"
)

const (
	promptSummarizeMapReduce = `Write a concise summary of the following:
"{{.text}}""
CONCISE SUMMARY:
	`
)

type MapReduceSummarizationChain struct {
	mapReduceCombineDocument *combine_document.MapReduceCombineDocument
	callbackManager          *callback.Manager
}

var _ chain.BaseChain = &MapReduceSummarizationChain{}

// NewMapReduceSummarizationChain create new map reduce summarization chain instance
// put empty "" string to use default prompt
// put 0 to use default maxToken
func NewMapReduceSummarizationChain(llmChain *llm_chain.LLMChain, callbackManager *callback.Manager, mapPromptString string, reducePromptString string,
	promptTemplateKey string,
	splitter textsplitter.TextSplitter, maxToken int, verbose bool) (m *MapReduceSummarizationChain, err error) {

	if verbose {
		callbackManager.RegisterCallback(chain.CallbackChainEnd, callback.VerboseCallback)
	}

	var promptTemplateMap, promptTemplateReduce *prompt.PromptTemplate

	if mapPromptString == "" {
		promptTemplateMap, err = prompt.NewPromptTemplate("map", promptSummarizeMapReduce)
		if err != nil {
			return
		}
		promptTemplateKey = "text"
	}

	if reducePromptString == "" {
		promptTemplateReduce, err = prompt.NewPromptTemplate("map", promptSummarizeMapReduce)
		if err != nil {
			return
		}
	}

	if maxToken == 0 {
		maxToken = 1000
	}

	mapReduceCombineDocument := combine_document.NewMapReduceCombineDocument(promptTemplateMap,
		promptTemplateReduce, promptTemplateKey, llmChain, splitter, maxToken, callbackManager, verbose)
	m = &MapReduceSummarizationChain{
		mapReduceCombineDocument: mapReduceCombineDocument,
		callbackManager:          callbackManager,
	}

	return
}

// Run expect input["input"] as input, and put the result to output["output"]
func (M *MapReduceSummarizationChain) Run(ctx context.Context, input map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if _, ok := input["input"]; !ok {
		return output, errors.New("input[\"input\"] is not specified")
	}

	//trigger callback chain start
	M.callbackManager.TriggerEvent(ctx, chain.CallbackChainStart, callback.CallbackData{
		FunctionName: "MapReduceSummarizationChain.Run",
		Input:        input,
		Output:       output,
	})

	output, err = M.mapReduceCombineDocument.Run(ctx, input, options...)

	// trigger callback chain end
	M.callbackManager.TriggerEvent(ctx, chain.CallbackChainEnd, callback.CallbackData{
		FunctionName: "MapReduceSummarizationChain.Run",
		Input:        input,
		Output:       output,
	})

	return
}

// SimpleRun will run the input prompt string againts llmchain
func (M *MapReduceSummarizationChain) SimpleRun(ctx context.Context, input string, options ...func(*model.Option)) (output string, err error) {
	//trigger callback chain start
	M.callbackManager.TriggerEvent(ctx, chain.CallbackChainStart, callback.CallbackData{
		FunctionName: "MapReduceSummarizationChain.SimpleRun",
		Input:        map[string]string{"input": input},
		Output:       map[string]string{"output": output},
	})

	output, err = M.mapReduceCombineDocument.SimpleRun(ctx, input, options...)

	// trigger callback chain end
	M.callbackManager.TriggerEvent(ctx, chain.CallbackChainEnd, callback.CallbackData{
		FunctionName: "MapReduceSummarizationChain.SimpleRun",
		Input:        map[string]string{"input": input},
		Output:       map[string]string{"output": output},
	})
	return
}
