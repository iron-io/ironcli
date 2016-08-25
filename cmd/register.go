package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/iron-io/iron_go3/api"
	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/iron_go3/worker"
	"github.com/urfave/cli"
)

type Register struct {
	name            *string
	config          *string
	configFile      *string
	maxConc         *int
	retries         *int
	retriesDelay    *int
	defaultPriority *int
	host            *string
	cmd             string
	codes           worker.Code

	cli.Command
}

func NewRegister(settings *config.Settings) *Register {
	register := &Register{}
	register.Command = cli.Command{
		Name:      "register",
		Usage:     "do the doo",
		UsageText: "doo - does the dooing",
		ArgsUsage: "[test]",
		Action: func(c *cli.Context) error {
			register.Run(settings)

			return nil
		},
	}

	return register
}

func (r Register) GetCmd() cli.Command {
	return r.Command
}

func (r *Register) Register() error {
	if r.name != nil && *r.name != "" {
		r.codes.Name = *r.name
	} else {
		r.codes.Name = r.codes.Image
		if strings.ContainsRune(r.codes.Name, ':') {
			arr := strings.SplitN(r.codes.Name, ":", 2)
			r.codes.Name = arr[0]
		}
	}

	r.codes.MaxConcurrency = *r.maxConc
	r.codes.Retries = *r.retries
	r.codes.RetriesDelay = *r.retriesDelay
	r.codes.Config = *r.config
	r.codes.DefaultPriority = *r.defaultPriority

	if r.host != nil && *r.host != "" {
		r.codes.Host = *r.host
	}

	if *r.configFile != "" {
		pload, err := ioutil.ReadFile(*r.configFile)
		if err != nil {
			return err
		}
		r.codes.Config = string(pload)
	}

	return nil
}

func (r *Register) pushCodes(zipName string, settings *config.Settings, args worker.Code) (*worker.Code, error) {
	// TODO i don't get why i can't write from disk to wire, but I give up
	var body bytes.Buffer
	mWriter := multipart.NewWriter(&body)
	mMetaWriter, err := mWriter.CreateFormField("data")
	if err != nil {
		return nil, err
	}

	jEncoder := json.NewEncoder(mMetaWriter)
	if err := jEncoder.Encode(args); err != nil {
		return nil, err
	}

	if zipName != "" {
		r, err := zip.OpenReader(zipName)
		if err != nil {
			return nil, err
		}
		defer r.Close()

		mFileWriter, err := mWriter.CreateFormFile("file", "worker.zip")
		if err != nil {
			return nil, err
		}
		zWriter := zip.NewWriter(mFileWriter)

		for _, f := range r.File {
			fWriter, err := zWriter.Create(f.Name)
			if err != nil {
				return nil, err
			}
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(fWriter, rc)
			rc.Close()
			if err != nil {
				return nil, err
			}
		}

		zWriter.Close()
	}
	mWriter.Close()

	req, err := http.NewRequest("POST", api.Action(*settings, "codes").URL.String(), &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip/deflate")
	req.Header.Set("Authorization", "OAuth "+settings.Token)
	req.Header.Set("Content-Type", mWriter.FormDataContentType())
	req.Header.Set("User-Agent", settings.UserAgent)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err = api.ResponseAsError(response); err != nil {
		return nil, err
	}

	var data worker.Code
	err = json.NewDecoder(response.Body).Decode(&data)
	return &data, err
}

func (r *Register) Run(settings *config.Settings) {
	if r.codes.Host != "" {
		fmt.Println(`Spinning up '` + r.codes.Name + `'`)
	} else {
		fmt.Println(`Registering worker '` + r.codes.Name + `'`)
	}

	code, err := r.pushCodes("", settings, r.codes)
	if err != nil {
		fmt.Println(err)
		return
	}

	if code.Host != "" {
		fmt.Println(`Hosted at: '` + code.Host + `'`)
	} else {
		fmt.Println(`Registered code package with id='` + code.Id + `'`)
	}
}
