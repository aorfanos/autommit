package utils

import (
	"context"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestAutommit_GetOpenPGPKeyring(t *testing.T) {
	type fields struct {
		Version      string
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
			name: "TestInvalidKeyPath",
			fields: fields{
				PgpKeyPath: "nonexistent",
			},
			wantErr: true,
		},
		{
			name: "TestValidKeyPath",
			fields: fields{
				PgpKeyPath: "./testdir/pgp.key.valid",
			},
			wantErr: false,
		},
		{
			name: "TestValidKeyPathWithInvalidKey",
			fields: fields{
				PgpKeyPath: "./testdir/pgp.key.invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Autommit{
				Version:      tt.fields.Version,
				OpenAiApiKey: tt.fields.OpenAiApiKey,
				Context:      tt.fields.Context,
				OpenAiClient: tt.fields.OpenAiClient,
				PgpKeyPath:   tt.fields.PgpKeyPath,
				CommitInfo:   tt.fields.CommitInfo,
				Type:         tt.fields.Type,
				GitConfig:    tt.fields.GitConfig,
			}
			if err := a.GetOpenPGPKeyring(); (err != nil) != tt.wantErr {
				t.Errorf("Autommit.GetOpenPGPKeyring() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
