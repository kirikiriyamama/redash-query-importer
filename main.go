package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

type Definition struct {
	Api struct {
		Base string `validate:"required"`
		Key  string `validate:"required"`
	} `validate:"required,dive"`
	Queries []struct {
		DataSourceId int    `yaml:"data_source_id" json:"data_source_id" validate:"required"`
		Query        string `json:"query" validate:"required"`
		Name         string `json:"name" validate:"required"`
		Description  string `json:"description,omitempty"`
		Schedule     string `json:"schedule,omitempty"`
	} `validate:"required,dive"`
}

func run() int {
	var err error

	file := flag.String("f", "", "Path to definition file")
	flag.Parse()

	if *file == "" {
		fmt.Fprintln(os.Stderr, "the required flag `-f' was not provided")
		return 64 // EX_USAGE
	}

	buf, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 66 // EX_NOINPUT
	}

	definition := Definition{}
	err = yaml.Unmarshal(buf, &definition)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 65 // EX_DATAERR
	}
	err = validator.New().Struct(definition)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 78 // EX_CONFIG
	}

	for _, query := range definition.Queries {
		reqBody, err := json.Marshal(query)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 70 // EX_SOFTWARE
		}

		req, err := http.NewRequest("POST", definition.Api.Base+"/api/queries", bytes.NewReader(reqBody))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 78 // EX_CONFIG
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Key "+definition.Api.Key)

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 69 // EX_UNAVAILABLE
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || 299 < res.StatusCode {
			fmt.Fprintln(os.Stderr, res.Status)
			return 76 // EX_PROTOCOL
		}
	}

	return 0 // EX_OK
}

func main() {
	os.Exit(run())
}
