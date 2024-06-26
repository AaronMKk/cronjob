package main

import (
	"awesomeProject/config"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/opensourceways/server-common-lib/logrusutil"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"

	liboptions "github.com/opensourceways/server-common-lib/options"
)

const (
	component = "merlin-cronjob"

	testGetURL = "https://dsapi.test.osinfra.cn/query/modelfoundry/download/count"
	ProdGetURL = "https://dsapi.test.osinfra.cn/query/modelfoundry/download/count"

	updateRepoURL = "http://172.28.223.236:8888/internal/v1/space/%s"
)

type DownloadData struct {
	Code int `json:"code"`
	Data []struct {
		Name     string `json:"name"`
		Download int    `json:"download"`
		RepoID   string `json:"id"`
	} `json:"data"`
}

type UpdateRepo struct {
	DownloadCount int `json:"download_count"`
}

func fetchDownloadCounts(cfg config.Config) (*DownloadData, error) {
	resp, err := http.Get(testGetURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data DownloadData
	err = json.Unmarshal(body, &data)
	return &data, err
}

func updateRepo(id string, count int) error {
	client := &http.Client{}
	data := UpdateRepo{DownloadCount: count}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf(updateRepoURL, id), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("TOKEN", "12345")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

type options struct {
	service   liboptions.ServiceOptions
	removeCfg bool
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.removeCfg, "rm-cfg", false,
		"whether remove the cfg file after initialized.",
	)

	if err := fs.Parse(args); err != nil {
		fs.PrintDefaults()

		logrus.Fatalf("failed to parse cmdline %s", err)
	}

	return o
}

func main() {
	logrusutil.ComponentInit(component)

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)

	// cfg
	cfg, err := config.LoadConfig(o.service.ConfigFile, o.removeCfg)

	data, err := fetchDownloadCounts(cfg)
	if err != nil {
		fmt.Printf("Error fetching download counts: %v\n", err)
		return
	}

	for _, codeRepo := range data.Data {
		err := updateRepo(codeRepo.RepoID, codeRepo.Download)

		if err != nil {
			fmt.Printf("Error updating download counts: %v for repo id: %s", err, codeRepo.RepoID)
		}

	}
}
