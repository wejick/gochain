package combine_document

import (
	"context"
	"reflect"
	"testing"

	"github.com/wejick/gochain/callback"
	"github.com/wejick/gochain/chain/llm_chain"
	"github.com/wejick/gochain/model"
	"github.com/wejick/gochain/prompt"
)

var emptyPrompt, _ = prompt.NewPromptTemplate("empty", "")
var echoPrompt, _ = prompt.NewPromptTemplate("empty", "{{.echo}}")
var echoLlmChain, _ = llm_chain.NewLLMChain(&model.LLMModelMock{
	CallFunc: func(ctx context.Context, prompt string, options ...func(*model.Option)) (string, error) {
		return prompt, nil
	},
}, callback.NewManager(), nil, false)

func TestStuffCombineDocument_Combine(t *testing.T) {

	type fields struct {
		prompt            *prompt.PromptTemplate
		llmChain          *llm_chain.LLMChain
		promptTemplateKey string
	}
	type args struct {
		ctx     context.Context
		docs    []string
		options []func(*model.Option)
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOutput string
		wantErr    bool
	}{
		{
			name: "empty",
			fields: fields{
				prompt:   emptyPrompt,
				llmChain: echoLlmChain,
			},
		},
		{
			name: "crowded, jakarta",
			fields: fields{
				prompt:            echoPrompt,
				llmChain:          echoLlmChain,
				promptTemplateKey: "echo",
			},
			args: args{
				ctx:  context.Background(),
				docs: []string{"crowded", "jakarta"},
			},
			wantOutput: "crowded\njakarta\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := &StuffCombineDocument{
				prompt:            tt.fields.prompt,
				llmChain:          tt.fields.llmChain,
				promptTemplateKey: tt.fields.promptTemplateKey,
			}
			gotOutput, err := S.Combine(tt.args.ctx, tt.args.docs, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("StuffCombineDocument.Combine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput != tt.wantOutput {
				t.Errorf("StuffCombineDocument.Combine() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestStuffCombineDocument_Run(t *testing.T) {
	type fields struct {
		prompt            *prompt.PromptTemplate
		llmChain          *llm_chain.LLMChain
		promptTemplateKey string
		callbackManager   *callback.Manager
	}
	type args struct {
		ctx     context.Context
		input   map[string]string
		options []func(*model.Option)
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOutput map[string]string
		wantErr    bool
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name: "crowded,jakarta",
			fields: fields{
				prompt:            echoPrompt,
				llmChain:          echoLlmChain,
				promptTemplateKey: "echo",
				callbackManager:   callback.NewManager(),
			},
			args: args{
				input: map[string]string{"input": "crowded,jakarta"},
			},
			wantOutput: map[string]string{"output": "crowded,jakarta\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := &StuffCombineDocument{
				prompt:            tt.fields.prompt,
				llmChain:          tt.fields.llmChain,
				promptTemplateKey: tt.fields.promptTemplateKey,
				callbackManager:   tt.fields.callbackManager,
			}
			gotOutput, err := S.Run(tt.args.ctx, tt.args.input, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("StuffCombineDocument.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("StuffCombineDocument.Run() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestStuffCombineDocument_SimpleRun(t *testing.T) {
	type fields struct {
		prompt            *prompt.PromptTemplate
		llmChain          *llm_chain.LLMChain
		callbackManager   *callback.Manager
		promptTemplateKey string
	}
	type args struct {
		ctx     context.Context
		input   string
		options []func(*model.Option)
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOutput string
		wantErr    bool
	}{
		{
			name: "empty",
			fields: fields{
				prompt:          emptyPrompt,
				llmChain:        echoLlmChain,
				callbackManager: callback.NewManager(),
			},
		},
		{
			name: "crowded,jakarta",
			fields: fields{
				prompt:            echoPrompt,
				llmChain:          echoLlmChain,
				callbackManager:   callback.NewManager(),
				promptTemplateKey: "echo",
			},
			args: args{
				input: "crowded,jakarta",
			},
			wantOutput: "crowded,jakarta\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := &StuffCombineDocument{
				prompt:            tt.fields.prompt,
				llmChain:          tt.fields.llmChain,
				callbackManager:   tt.fields.callbackManager,
				promptTemplateKey: tt.fields.promptTemplateKey,
			}
			gotOutput, err := S.SimpleRun(tt.args.ctx, tt.args.input, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("StuffCombineDocument.SimpleRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput != tt.wantOutput {
				t.Errorf("StuffCombineDocument.SimpleRun() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
