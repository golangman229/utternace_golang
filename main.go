package main

import (
	"fmt"

	"github.com/Utterance/utterance"
)

func main() {

	folder := "utterances"

	utterances, err := utterance.ReadUtterances(folder)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	// Process the utterances and get the transcript
	transcript := utterance.ProcessUtterances(utterances)

	// Write the transcript to the output file
	err = utterance.WriteTranscript(transcript, "data/output.json")
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}

	fmt.Println("Successfully generated complete transcript.", utterances)
}
