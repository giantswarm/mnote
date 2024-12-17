package models

import (
	"testing"

	"github.com/giantswarm/mnote/internal/config"
)

func TestGetWhisperModel(t *testing.T) {
	cfg := config.DefaultConfig()

	tests := []struct {
		name    string
		lang    string
		want    string
		wantErr bool
	}{
		{
			name: "english language",
			lang: "en",
			want: DefaultEnglishModel,
		},
		{
			name: "german language",
			lang: "de",
			want: DefaultLargeModel,
		},
		{
			name: "auto detection",
			lang: "auto",
			want: DefaultEnglishModel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetWhisperModel(tt.lang, cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWhisperModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetWhisperModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateLanguage(t *testing.T) {
	tests := []struct {
		name string
		lang string
		want bool
	}{
		{"valid auto", "auto", true},
		{"valid english", "en", true},
		{"valid german", "de", true},
		{"valid spanish", "es", true},
		{"valid french", "fr", true},
		{"invalid language", "invalid", false},
		{"empty language", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateLanguage(tt.lang); got != tt.want {
				t.Errorf("ValidateLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSupportedLanguages(t *testing.T) {
	want := []string{"auto", "en", "de", "es", "fr"}
	got := GetSupportedLanguages()

	if len(got) != len(want) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(got), len(want))
	}

	for i, lang := range want {
		if got[i] != lang {
			t.Errorf("GetSupportedLanguages()[%d] = %v, want %v", i, got[i], lang)
		}
	}
}
