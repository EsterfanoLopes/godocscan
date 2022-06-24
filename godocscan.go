package main

import (
	"fmt"
	"regexp"

	"github.com/otiai10/gosseract/v2"
)

const (
	// Regex_ddMMMyyyy to RG document dates
	Regex_ddMMMyyyy = `(?m)(([0-9])|([0-2][0-9])|([3][0-1]))\/(JAN|FEV|MAR|ABR|MAI|JUN|JUL|AGO|SET|OUT|NOV|DEZ)\/\d{4}`

	// Regex_CPF to CPF document number
	Regex_CPF = `(?m)^\d{3}.?\d{3}.?\d{3}-?\d{2}`

	// Regex_RG to RG document number
	Regex_RG = `(?m)(^\d{1,2}).?(\d{3}).?(\d{3})-?(\d{1,2}|X|x$)`
)

func Setup(language string) (client *gosseract.Client, err error) {
	if language == "" {
		language = "por"
	}

	client = gosseract.NewClient()
	err = client.SetLanguage(language)
	if err != nil {
		return
	}

	return
}

type RegistroGeral struct {
	Client    *gosseract.Client
	Number    string
	IssueDate string
	BirthDate string
	CpfNumber string
	Raw       string
}

func (rg *RegistroGeral) New(language string) (err error) {
	rg.Client, err = Setup(language)
	if err != nil {
		return
	}

	return
}

func (rg *RegistroGeral) SetImage(imagePath string) (err error) {
	err = rg.Client.SetImage(imagePath)
	if err != nil {
		return
	}

	return
}

func (rg *RegistroGeral) ScanDates() (err error) {
	text, err := rg.Client.Text()
	if err != nil {
		return
	}

	result := rg.matchDate(text)

	rg.IssueDate = result[0]
	rg.BirthDate = result[1]

	return
}

func (rg *RegistroGeral) ScanRaw() (err error) {
	rg.Raw, err = rg.Client.Text()
	return
}

func (rg *RegistroGeral) ScanCPF() (err error) {
	text, err := rg.Client.Text()
	if err != nil {
		return
	}

	result := rg.matchCpfNumber(text)

	if len(result) > 0 {
		rg.CpfNumber = result[0]
	}

	return
}

func (rg *RegistroGeral) ScanRgNumber() (err error) {
	text, err := rg.Client.Text()
	if err != nil {
		return
	}

	result := rg.matchRgNumber(text)

	if len(result) > 0 {
		rg.Number = result[0]
	}

	return
}

func (rg *RegistroGeral) ScanAll() (err error) {
	err = rg.ScanRaw()
	if err != nil {
		return err
	}

	err = rg.ScanDates()
	if err != nil {
		return err
	}

	err = rg.ScanCPF()
	if err != nil {
		return err
	}

	err = rg.ScanRgNumber()
	if err != nil {
		return err
	}

	return
}

func (rg *RegistroGeral) matchDate(text string) []string {
	return regexp.MustCompile(Regex_ddMMMyyyy).FindAllString(text, -1)
}

func (rg *RegistroGeral) matchCpfNumber(text string) []string {
	return regexp.MustCompile(Regex_CPF).FindAllString(text, -1)
}

func (rg *RegistroGeral) matchRgNumber(text string) []string {
	return regexp.MustCompile(Regex_RG).FindAllString(text, -1)
}

func main() {
	rg := RegistroGeral{}
	defer func() {
		err := rg.Client.Close()
		if err != nil {
			fmt.Errorf(`Error closing client: %s`, err)
			return
		}
	}()

	rg.New("por")

	err := rg.SetImage("docs/examples/rgexample.jpg")
	if err != nil {
		fmt.Printf("Error setting image: %s", err)
		return
	}

	err = rg.ScanAll()
	if err != nil {
		fmt.Printf("Error scanning dates: %s", err)
		return
	}

	fmt.Printf(`Results: %v`, rg)
	return
}
