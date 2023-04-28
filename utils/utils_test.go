package utils

import (
	"context"
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

func TestShowVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestVersion",
			args: args{
				version: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ShowVersion(tt.args.version)
		})
	}
}

func TestFindDotGit(t *testing.T) {
	type args struct {
		repoPath string
	}
	tests := []struct {
		name     string
		args     args
		wantPath string
		wantErr  bool
	}{
		{
			name: "TestNonExistingRepo",
			args: args{
				repoPath: "/home/testdir",
			},
			wantPath: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, err := FindDotGit(tt.args.repoPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindDotGit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPath != tt.wantPath {
				t.Errorf("FindDotGit() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func Test_getDirectoryLevelsToRoot(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "TestUtilsGHA", // adjusted to GitHub Actions' directory structure
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDirectoryLevelsToRoot(); got != tt.want {
				t.Errorf("getDirectoryLevelsToRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}
