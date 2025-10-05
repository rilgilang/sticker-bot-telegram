package connector

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

type OCRResult struct {
	ParsedResults []struct {
		TextOverlay struct {
			Lines []struct {
				LineText string `json:"LineText"`
				Words    []struct {
					WordText string  `json:"WordText"`
					Left     float64 `json:"Left"`
					Top      float64 `json:"Top"`
					Height   float64 `json:"Height"`
					Width    float64 `json:"Width"`
				} `json:"Words"`
				MaxHeight float64 `json:"MaxHeight"`
				MinTop    float64 `json:"MinTop"`
			} `json:"Lines"`
			HasOverlay bool `json:"HasOverlay"`
		} `json:"TextOverlay"`
		TextOrientation   string `json:"TextOrientation"`
		FileParseExitCode int    `json:"FileParseExitCode"`
		ParsedText        string `json:"ParsedText"`
		ErrorMessage      string `json:"ErrorMessage"`
		ErrorDetails      string `json:"ErrorDetails"`
	} `json:"ParsedResults"`
	OCRExitCode                  int    `json:"OCRExitCode"`
	IsErroredOnProcessing        bool   `json:"IsErroredOnProcessing"`
	ProcessingTimeInMilliseconds string `json:"ProcessingTimeInMilliseconds"`
	SearchablePDFURL             string `json:"SearchablePDFURL"`
}

type APIKeys string

func (k *APIKeys) ProcessOCR(rest *resty.Client, filePath *string, fileType string) (*OCRResult, error) {

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	contentType := ""
	switch fileType {
	case "webp":
		contentType = "image/webp"
	case "webm":
		contentType = "video/webm"
	}

	result := &OCRResult{}

	response, err := rest.R().
		SetResult(&result).
		//SetHeader("apikey", string(*k)).
		SetHeader("apikey", "donotstealthiskey_ip1").
		SetFormData(map[string]string{
			"url":                          "",
			"language":                     "eng",
			"isOverlayRequired":            "true",
			"FileType":                     ".Auto",
			"IsCreateSearchablePDF":        "false",
			"isSearchablePdfHideTextLayer": "true",
			"detectOrientation":            "false",
			"isTable":                      "false",
			"scale":                        "true",
			"OCREngine":                    "5",
			"detectCheckbox":               "false",
			"checkboxTemplate":             "0",
		}).
		SetMultipartField("file", file.Name(), contentType, file).
		Post(`https://api8.ocr.space/parse/image`)

	if err != nil {
		return nil, err
	}

	if response.StatusCode() >= 400 {
		return nil, errors.New(fmt.Sprintf(`error processing ocr: %s`, string(response.Body())))
	}

	return result, nil
}
