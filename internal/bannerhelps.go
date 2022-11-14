package internal

import "context"

type Service interface {
	PdfToVoice(context.Context, PdfToVoiceRequest) PdfToVoiceResponse
}

type PdfToVoiceRequest struct {
	FileName        string
	File            []byte
	CurrentLanguage string
}

type PdfToVoiceResponse struct {
	IsSuccessfully bool
	FilePath       string
	Error          *ExternalAError
}

type ExternalAError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
