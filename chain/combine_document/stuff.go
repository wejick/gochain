package combine_document

import (
	"context"
	"errors"

	"github.com/wejick/gochain/callback"
	"github.com/wejick/gochain/chain"
	"github.com/wejick/gochain/chain/llm_chain"
	"github.com/wejick/gochain/model"
	"github.com/wejick/gochain/prompt"
)

var _ CombinedDocument = &StuffCombineDocument{}
var _ chain.BaseChain = &StuffCombineDocument{}

// StuffCombineDocument chain to feed text document to LLM with specified prompt
type StuffCombineDocument struct {
	prompt            *prompt.PromptTemplate
	llmChain          *llm_chain.LLMChain
	callbackManager   *callback.Manager
	promptTemplateKey string
}

// NewStuffCombineDocument creates new instance of StuffCombineDocument
func NewStuffCombineDocument(callbackManager *callback.Manager, prompt *prompt.PromptTemplate,
	templateKey string, llmChain *llm_chain.LLMChain, verbose bool) *StuffCombineDocument {

	if verbose {
		callbackManager.RegisterCallback(chain.CallbackChainEnd, callback.VerboseCallback)
	}

	return &StuffCombineDocument{
		prompt:            prompt,
		llmChain:          llmChain,
		callbackManager:   callbackManager,
		promptTemplateKey: templateKey,
	}
}

// Combine concatenate the given document and then feed to LLM
func (S *StuffCombineDocument) Combine(ctx context.Context, docs []string, options ...func(*model.Option)) (output string, err error) {
	//concat all docs into 1 string
	var doc string
	for _, item := range docs {
		doc += item + "\n"
	}
	templateData := map[string]string{S.promptTemplateKey: doc}

	prompt, err := S.prompt.FormatPrompt(templateData)
	if err != nil {
		return
	}
	output, err = S.llmChain.SimpleRun(ctx, prompt)

	return
}

// Run expect input["input"] as input, and put the result to output["output"]
func (S *StuffCombineDocument) Run(ctx context.Context, input map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if _, ok := input["input"]; !ok {
		return output, errors.New("input[\"input\"] is not specified")
	}
	//trigger callback chain start
	S.callbackManager.TriggerEvent(ctx, chain.CallbackChainStart, callback.CallbackData{
		FunctionName: "StuffCombineDocument.Run",
		Input:        input,
		Output:       output,
	})
	output = make(map[string]string)
	output["output"], err = S.Combine(ctx, []string{input["input"]})

	//trigger callback chain end
	S.callbackManager.TriggerEvent(ctx, chain.CallbackChainEnd, callback.CallbackData{
		FunctionName: "StuffCombineDocument.Run",
		Input:        input,
		Output:       output,
	})

	return
}

// SimpleRun will run the input string agains llmchain
func (S *StuffCombineDocument) SimpleRun(ctx context.Context, input string, options ...func(*model.Option)) (output string, err error) {
	//trigger callback chain start
	S.callbackManager.TriggerEvent(ctx, chain.CallbackChainStart, callback.CallbackData{
		FunctionName: "StuffCombineDocument.SimpleRun",
		Input:        map[string]string{"input": input},
		Output:       map[string]string{"output": output},
	})

	output, err = S.Combine(ctx, []string{input})

	//trigger callback chain end
	S.callbackManager.TriggerEvent(ctx, chain.CallbackChainEnd, callback.CallbackData{
		FunctionName: "StuffCombineDocument.SimpleRun",
		Input:        map[string]string{"input": input},
		Output:       map[string]string{"output": output},
	})
	return
}
