package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ledongthuc/pdf"
	"net/http"
	"os"
)

const (
	BasePath  = "./docs/pdf"
	Extension = ".pdf"
)

// Client defines behaviors of pdf client
type Client interface {
	ConvertToText(path string) (string, error)
	Save(file []byte, clientID string) (string, error)
	Delete(clientID string, filePath string) error
}

type client struct {
	c *http.Client
}

// NewClient creates and returns pdf client
func NewClient(c *http.Client) Client {
	cli := &client{
		c: c,
	}

	return cli
}

func (c client) ConvertToText(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	_, _ = buf.ReadFrom(b)
	text := buf.String()

	err = f.Close()
	if err != nil {
		return "", err
	}

	return text, err
}

func (c client) Save(file []byte, clientID string) (string, error) {
	folder := fmt.Sprintf("%s/%s", BasePath, clientID)
	fName := fmt.Sprintf("%s%s", uuid.New().String(), Extension)
	filePath := fmt.Sprintf("%s/%s", folder, fName)

	if _, err := os.Stat(folder); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	_, err = f.Write(file)
	if err != nil {
		return "", err
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (c client) Delete(clientID string, filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	folderPath := fmt.Sprintf("%s/%s", BasePath, clientID)
	err = os.Remove(folderPath)
	if err != nil {
		return err
	}

	return nil
}
