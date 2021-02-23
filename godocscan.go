package main

import (
	"fmt"
	"regexp"

	"github.com/otiai10/gosseract/v2"
)

const (
	// ddMMMyyyyREGEX to RG document dates
	ddMMMyyyyREGEX = `(([0-9])|([0-2][0-9])|([3][0-1]))\/(JAN|FEV|MAR|ABR|MAI|JUN|JUL|AGO|SET|OUT|NOV|DEZ)\/\d{4}`
)

func setup() (client *gosseract.Client, err error) {
	client = gosseract.NewClient()

	err = client.SetLanguage("por")
	if err != nil {
		return
	}

	return
}

func matchDate(text string) []string {
	re := regexp.MustCompile(ddMMMyyyyREGEX)

	return re.FindAllString(text, -1)
}

func ReadRG(client *gosseract.Client) (result []string, err error) {
	err = client.SetImage("docs/examples/rgexample.jpg")
	if err != nil {
		return
	}
	text, err := client.Text()
	if err != nil {
		return
	}

	result = matchDate(text)

	return
}

func main() {
	cli, err := setup()
	if err != nil {
		fmt.Printf(`Error: %s`, err.Error())
	}

	defer cli.Close()

	result, err := ReadRG(cli)
	if err != nil {
		fmt.Printf("aaa %s", err)
	}

	fmt.Printf(`Results: %v`, result)

	return
}
