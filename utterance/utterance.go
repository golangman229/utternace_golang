package utterance

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Utterance represents a transcript event with speaker, text, and timestamp
type Utterance struct {
	Speaker     string `json:"speaker"`
	Text        string `json:"text"`
	TimestampMs int64  `json:"timestampMs"`
}

// Transcript represents a collection of utterances
type Transcript struct {
	Utterances []Utterance `json:"utterances"`
}

// IsFragment checks if a word is a sentence fragment or not
func IsFragment(word string) bool {
	// A simple heuristic is to check if the word is lowercase and does not end with punctuation
	return !strings.ContainsAny(word, ".?!")
}

// FixAttribution fixes the speaker attribution for a given utterance and the previous one
func FixAttribution(utterance *Utterance, prev *Utterance) {
	// Split the text into words
	words := strings.Split(prev.Text, " ")

	// If the first word is a fragment, append it to the previous utterance and remove it from the current one
	for IsFragment(words[len(words)-1]) && len(words) > 1 {
		if words[len(words)-1] != "" {
			utterance.Text = words[len(words)-1] + " " + utterance.Text
		}
		words = words[:len(words)-1]
	}

	if len(words) == 1 {
		if IsFragment((words[0])) && words[0] != "" {
			utterance.Text = words[0] + " " + utterance.Text
			words[0] = ""
		}
	}

	prev.Text = strings.Join(words, " ")
}

// ProcessUtterances processes a slice of utterances and returns a transcript
func ProcessUtterances(utterances []Utterance) Transcript {
	// Create an empty transcript
	transcript := Transcript{}
	length := len(utterances)

	// Loop through the utterances
	for i, utterance := range utterances {
		// If this is not the first utterance, fix the attribution with the previous one
		if i < length-1 {
			FixAttribution(&utterances[i+1], &utterance)
		}

		// Append the utterance to the transcript
		if utterance.Text != "" {
			transcript.Utterances = append(transcript.Utterances, utterance)
		}
	}

	return transcript
}

// ReadUtterances reads a JSON file and returns a slice of utterances
func ReadUtterances(foldername string) ([]Utterance, error) {
	// Loop through the files in the folder
	var utterances []Utterance

	err := filepath.Walk(foldername, func(path string, info os.FileInfo, err error) error {
		// Check if the file has a .json extension
		if filepath.Ext(path) == ".json" {
			// Open the file
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			// Close the file when done
			defer f.Close()
			// Read the file contents
			data, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			// Unmarshal the JSON data
			var m Utterance
			err = json.Unmarshal(data, &m)
			if err != nil {
				return err
			}
			// Append the map to the slice
			utterances = append(utterances, m)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return utterances, nil
}

// WriteTranscript writes a transcript to a JSON file
func WriteTranscript(transcript Transcript, filename string) error {
	data, err := json.MarshalIndent(transcript, "", "")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
