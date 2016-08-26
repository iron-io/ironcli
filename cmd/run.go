package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/iron-io/iron_go3/api"
	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/iron_go3/worker"
	"github.com/urfave/cli"
)

type Run struct {
	name            string
	config          string
	configFile      string
	maxConc         int
	retries         int
	retriesDelay    int
	defaultPriority int
	zip             string
	host            string
	codes           worker.Code

	cli.Command
}

func NewRun(settings *config.Settings) *Run {
	run := &Run{}
	run.Command = cli.Command{
		Name:      "run",
		Usage:     "do the doo",
		UsageText: "doo - does the dooing",
		ArgsUsage: "[image] [args]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "name",
				Usage:       "name usage",
				Destination: &run.name,
			},
			cli.StringFlag{
				Name:        "config",
				Usage:       "config usage",
				Destination: &run.config,
			},
			cli.StringFlag{
				Name:        "configFile",
				Usage:       "configFile usage",
				Destination: &run.configFile,
			},
			cli.IntFlag{
				Name:        "maxConc",
				Usage:       "maxConc usage",
				Value:       -1,
				Destination: &run.maxConc,
			},
			cli.IntFlag{
				Name:        "retries",
				Usage:       "retries usage",
				Value:       0,
				Destination: &run.retries,
			},
			cli.IntFlag{
				Name:        "retriesDelay",
				Usage:       "retriesDelay usage",
				Value:       0,
				Destination: &run.retriesDelay,
			},
			cli.IntFlag{
				Name:        "defaultPriority",
				Usage:       "defaultPriority usage",
				Value:       -3,
				Destination: &run.defaultPriority,
			},
			cli.StringFlag{
				Name:        "zip",
				Usage:       "zip usage",
				Destination: &run.zip,
			},
			cli.StringFlag{
				Name:        "host",
				Usage:       "host usage",
				Destination: &run.host,
			},
		},
		Action: func(c *cli.Context) error {
			err := run.Execute(c.Args().Tail(), c.Args().First())
			if err != nil {
				return err
			}

			if run.codes.Host != "" {
				fmt.Println(`Spinning up '` + run.codes.Name + `'`)
			} else {
				fmt.Println(`Registering worker '` + run.codes.Name + `'`)
			}

			code, err := run.pushCodes(run.zip, settings, run.codes)
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

	return run
}

func (r Run) GetCmd() cli.Command {
	return r.Command
}

func (r *Run) pushCodes(zipName string, settings *config.Settings, args worker.Code) (*worker.Code, error) {
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

func (r *Run) Execute(cmd []string, image string) error {
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

	if r.zip != "" {
		if !strings.HasSuffix(r.zip, ".zip") {
			return errors.New("file extension must be .zip, got: " + r.zip)
		}

		if _, err := os.Stat(r.zip); err != nil {
			return err
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
