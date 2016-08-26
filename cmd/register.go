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
	name            string
	config          string
	configFile      string
	maxConc         int
	retries         int
	retriesDelay    int
	defaultPriority int
	host            string
	codes           worker.Code

	cli.Command
}

func NewRegister(settings *config.Settings) *Register {
	register := &Register{}
	register.Command = cli.Command{
		Name:      "register",
		Usage:     "register worker in the project",
		UsageText: "doo - does the dooing",
		ArgsUsage: "[image] [command] [args]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "name",
				Usage:       "name usage",
				Destination: &register.name,
			},
			cli.StringFlag{
				Name:        "config",
				Usage:       "config usage",
				Destination: &register.config,
			},
			cli.StringFlag{
				Name:        "configFile",
				Usage:       "configFile usage",
				Destination: &register.configFile,
			},
			cli.IntFlag{
				Name:        "maxConc",
				Usage:       "maxConc usage",
				Value:       -1,
				Destination: &register.maxConc,
			},
			cli.IntFlag{
				Name:        "retries",
				Usage:       "retries usage",
				Value:       0,
				Destination: &register.retries,
			},
			cli.IntFlag{
				Name:        "retriesDelay",
				Usage:       "retriesDelay usage",
				Value:       0,
				Destination: &register.retriesDelay,
			},
			cli.IntFlag{
				Name:        "defaultPriority",
				Usage:       "defaultPriority usage",
				Value:       -3,
				Destination: &register.defaultPriority,
			},
			cli.StringFlag{
				Name:        "host",
				Usage:       "host usage",
				Destination: &register.host,
			},
		},
		Action: func(c *cli.Context) error {
			err := register.Execute(c.Args().Tail(), c.Args().First())
			if err != nil {
				return err
			}

			if register.codes.Host != "" {
				fmt.Println(`Spinning up '` + register.codes.Name + `'`)
			} else {
				fmt.Println(`Registering worker '` + register.codes.Name + `'`)
			}

			code, err := register.pushCodes("", settings, register.codes)
			if err != nil {
				return err
			}

			if code.Host != "" {
				fmt.Println(`Hosted at: '` + code.Host + `'`)
			} else {
				fmt.Println(`Registered code package with id='` + code.Id + `'`)
			}

			return nil
		},
	}

	return register
}

func (r Register) GetCmd() cli.Command {
	return r.Command
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

func (r *Register) Execute(cmd []string, image string) error {
	r.codes.Command = strings.TrimSpace(strings.Join(cmd, " "))
	r.codes.Image = image
	r.codes.Name = image

	if r.name != "" {
		r.codes.Name = r.name
	} else {
		r.codes.Name = r.codes.Image
		if strings.ContainsRune(r.codes.Name, ':') {
			arr := strings.SplitN(r.codes.Name, ":", 2)
			r.codes.Name = arr[0]
		}
	}

	r.codes.MaxConcurrency = r.maxConc
	r.codes.Retries = r.retries
	r.codes.RetriesDelay = r.retriesDelay
	r.codes.Config = r.config
	r.codes.DefaultPriority = r.defaultPriority

	if r.host != "" {
		r.codes.Host = r.host
	}

	if r.configFile != "" {
		pload, err := ioutil.ReadFile(r.configFile)
		if err != nil {
			return err
		}
		r.codes.Config = string(pload)
	}

	return nil
}
