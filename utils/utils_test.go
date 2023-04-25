package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestAutommit_ParseStringAsJson(t *testing.T) {
	type fields struct {
		OpenAiApiKey string
		Context      context.Context
		OpenAiClient openai.Client
		PgpSign      bool
		CommitInfo   Commit
	}
	type args struct {
		strSrc string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "TestEmptyJSON",
			fields: fields{},
			args: args{
				strSrc: "{}",
			},
			wantErr: true,
		},
		{
			name:   "TestInvalidJSON",
			fields: fields{},
			args: args{
				strSrc: "{",
			},
			wantErr: true,
		},
		{
			name:   "TestValidJSON",
			fields: fields{},
			args: args{
				strSrc: `{"commit_message": "test", "commit_message_long": "test"}`,
			},
		},
		{
			name:   "TestValidJSONWithWrongFields",
			fields: fields{},
			args: args{
				strSrc: `{"commeet_message": "test", "commit_message_wrong": "test"}`,
			},
			wantErr: true,
		},
		{
			name:   "TestValidJSONWithExtraFields",
			fields: fields{},
			args: args{
				strSrc: `{"commit_message": "test", "commit_message_long": "test", "changed_files": "test"}`,
			},
		},
		{
			name:   "TestOneFieldEmpty",
			fields: fields{},
			args: args{
				strSrc: `{"commit_message": "", "commit_message_long": "test"}`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Autommit{
				OpenAiApiKey: tt.fields.OpenAiApiKey,
				Context:      tt.fields.Context,
				OpenAiClient: tt.fields.OpenAiClient,
				CommitInfo:   tt.fields.CommitInfo,
			}
			if err := a.ParseStringAsJson(tt.args.strSrc); (err != nil) != tt.wantErr {
				t.Errorf("Autommit.ParseStringAsJson() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProceedSelector(t *testing.T) {
	type args struct {
		title   string
		choices []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "TestEmptyChoices",
			args: args{
				title:   "test",
				choices: []string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProceedSelector(tt.args.title, tt.args.choices)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProceedSelector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProceedSelector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrCheck(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestNilError",
			args: args{
				err: nil,
			},
		},
		{
			name: "TestExistingError",
			args: args{
				err: fmt.Errorf("test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ErrCheck(tt.args.err)
		})
	}
}

func TestProceedEditor(t *testing.T) {
	type args struct {
		title  string
		target string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "TestEmptyTitle",
			args: args{
				title:  "",
				target: "test",
			},
			wantErr: true,
		},
		{
			name: "TestEmptyTarget",
			args: args{
				title:  "test",
				target: "",
			},
			wantErr: true,
		},
		{
			name: "TestEmptyTitleAndTarget",
			args: args{
				title:  "",
				target: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProceedEditor(tt.args.title, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProceedEditor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProceedEditor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAutommit_PopulateGitUserInfo(t *testing.T) {
	type fields struct {
		OpenAiApiKey string
		Context      context.Context
		OpenAiClient openai.Client
		PgpKeyPath   string
		CommitInfo   Commit
		Type         string
		GitConfig    GitConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestNonExistingGitConfig",
			fields: fields{
				GitConfig: GitConfig{
					FilePath: "/home/testdir/.gitconfig",
				},
			},
			wantErr: true,
		},
		{
			name: "TestExistingGitConfig",
			fields: fields{
				GitConfig: GitConfig{
					FilePath: "./testdir/.gitconfig",
				},
			},
			wantErr: false,
		},
		{
			name: "TestPartialGitConfig",
			fields: fields{
				GitConfig: GitConfig{
					FilePath: "./testdir/.gitconfig.partial",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Autommit{
				OpenAiApiKey: tt.fields.OpenAiApiKey,
				Context:      tt.fields.Context,
				OpenAiClient: tt.fields.OpenAiClient,
				PgpKeyPath:   tt.fields.PgpKeyPath,
				CommitInfo:   tt.fields.CommitInfo,
				Type:         tt.fields.Type,
				GitConfig:    tt.fields.GitConfig,
			}
			if err := a.PopulateGitUserInfo(); (err != nil) != tt.wantErr {
				t.Errorf("Autommit.PopulateGitUserInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
