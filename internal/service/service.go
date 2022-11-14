package service

import (
	"bannerhelps/internal"
	"bannerhelps/internal/client/pdf"
	"bannerhelps/internal/client/voice"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
)

// compile-time proofs of service interface implementation
var _ internal.Service = (*Service)(nil)

// Service represents service
type Service struct {
	l   log.Logger
	env string
	voc voice.Client
	pfc pdf.Client
}

// NewService creates and returns service
func NewService(env string, voc voice.Client, pfc pdf.Client) internal.Service {
	return &Service{
		env: env,
		voc: voc,
		pfc: pfc,
	}
}

func (s *Service) PdfToVoice(_ context.Context, req internal.PdfToVoiceRequest) internal.PdfToVoiceResponse {
	res := internal.PdfToVoiceResponse{}
	clientID := uuid.New().String()

	pdfFilePath, err := s.pfc.Save(req.File, clientID)
	if err != nil {
		res.IsSuccessfully = false
		res.Error = &internal.ExternalAError{
			Code:    400,
			Message: err.Error(),
		}

		return res
	}

	text, err := s.pfc.ConvertToText(pdfFilePath)
	if err != nil {
		res.IsSuccessfully = false
		res.Error = &internal.ExternalAError{
			Code:    400,
			Message: err.Error(),
		}

		return res
	}

	files, err := s.voc.Speak(text, req.CurrentLanguage, clientID)
	filePath, err := s.voc.MergeAllSpeak(req.FileName, files, clientID)
	if err != nil {
		res.IsSuccessfully = false
		res.Error = &internal.ExternalAError{
			Code:    400,
			Message: err.Error(),
		}

		return res
	}

	err = s.pfc.Delete(clientID, pdfFilePath)
	if err != nil {
		fmt.Println(err.Error())
	}

	res.IsSuccessfully = true
	res.FilePath = filePath

	return res
}
