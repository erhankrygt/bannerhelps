package voice

import (
	"bytes"
	"fmt"
	htg "github.com/hegedustibor/htgo-tts"
	"github.com/hyacinthus/mp3join"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	BasePath           = "./docs/audio"
	MaxCharacterLength = 200
	Extension          = ".mp3"
	Regex              = `[^a-zA-Z0-9 ]+`
)

// Client defines behaviors of voice client
type Client interface {
	Speak(text string, currentLang string, clientID string) ([]string, error)
	MergeAllSpeak(fileName string, files []string, clientID string) (fullPath string, err error)
}

type client struct {
	c *http.Client
}

// NewClient creates and returns voice client
func NewClient(c *http.Client) Client {
	cli := &client{
		c: c,
	}

	return cli
}

func (c client) Speak(text string, currentLang string, clientID string) ([]string, error) {
	var err error

	folder := fmt.Sprintf("%s/%s", BasePath, clientID)
	speech := htg.Speech{Folder: folder, Language: currentLang}

	var partText, partNextText string
	var files []string
	text = regexp.MustCompile(Regex).ReplaceAllString(text, "")
	words := strings.Split(text, " ")

	for i := 0; i < len(words); i++ {
		partText = strings.Trim(fmt.Sprintf("%s %s ", partText, words[i]), " ")
		fileName := fmt.Sprintf("File-%d", i)

		if i+1 >= len(words) {
			f, err := speech.CreateSpeechFile(partText, fileName)
			if err != nil {
				return nil, err
			}

			files = append(files, f)
			return files, err
		}

		partNextText = strings.Trim(fmt.Sprintf("%s %s ", partText, words[i+1]), " ")
		if len(partText) <= MaxCharacterLength && len(partNextText) > MaxCharacterLength {
			f, err := speech.CreateSpeechFile(partText, fileName)
			if err != nil {
				return nil, err
			}

			files = append(files, f)
			partText, partNextText = "", ""
		}
	}

	return files, err
}

func (c client) MergeAllSpeak(fileName string, files []string, clientID string) (fullPath string, err error) {
	joiner := mp3join.New()

	for _, f := range files {
		reader, err := os.Open(f)
		if err != nil {
			return "", nil
		}

		err = joiner.Append(reader)
		if err != nil {
			return "", nil
		}

		err = os.Remove(f)
		if err != nil {
			return "", nil
		}
	}

	newFileName := strings.ReplaceAll(fileName, " ", "-")
	newFileName = strings.ReplaceAll(fileName, "&", "-")
	newFileName = strings.ReplaceAll(fileName, "%", "-")
	newFileName = strings.ReplaceAll(fileName, "+", "-")
	newFileName = strings.ReplaceAll(fileName, "'", "-")
	newFileName = strings.ReplaceAll(fileName, "<", "-")
	newFileName = strings.ReplaceAll(fileName, ">", "-")

	folder := fmt.Sprintf("%s/%s", BasePath, clientID)
	fName := fmt.Sprintf("%s%s", strings.ToLower(newFileName), Extension)
	filePath := fmt.Sprintf("%s/%s", folder, fName)

	fo, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(joiner.Reader())
	if err != nil {
		return "", err
	}

	_, err = fo.Write(buf.Bytes())
	if err != nil {
		return "", err
	}

	return filePath, nil
}
