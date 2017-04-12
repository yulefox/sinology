package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/lunny/html2md"
)

var cfgBlob = []byte(`{
	"host": "xn--5rtnx620bw5s.tw",
	"fmt_dir": "a/a%02d",
	"fmt_url": "http://%s/%s/%s",
	"fmt_file": "%03d",
	"in_ext": ".htm",
	"out_ext": ".md",
	"arg_0": 3,
	"arg_1a": [1,1],
	"arg_1e": []
}`)

var cfg config

type config struct {
	Host       string   `json:"host"`
	FormatDir  string   `json:"fmt_dir"`
	FormatURL  string   `json:"fmt_url"`
	FormatFile string   `json:"fmt_file"`
	InExt      string   `json:"in_ext"`
	OutExt     string   `json:"out_ext"`
	Arg0       int      `json:"arg_0"`
	Arg1A      []int    `json:"arg_1a"`
	Arg1E      []string `json:"arg_1e"`
	LocalDir   string
	URL        string
	File       string
}

func init() {
	err := json.Unmarshal(cfgBlob, &cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func chdir() {
	if err := os.MkdirAll(cfg.LocalDir, 0777); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := os.Chdir(cfg.LocalDir); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func conv() {
	inFile := cfg.File + cfg.InExt
	outFile := cfg.File + cfg.OutExt
	html, _ := ioutil.ReadFile(inFile)
	html2md.AddRule("span", &html2md.Rule{
		Patterns: []string{"span"},
		Replacement: func(innerHTML string, attrs []string) string {
			if len(attrs) > 1 {
				return "`" + attrs[1] + "`"
			}
			return ""
		},
	})
	md := html2md.Convert(string(html))
	ioutil.WriteFile(outFile, []byte(md), 0644)
}

func download() {
	fmt.Println(cfg.URL)
	cmd := fmt.Sprintf("curl -O %s", cfg.URL)
	_, err := exec.Command("/bin/sh", "-c", cmd).Output()

	if err != nil {
		log.Println(err)
	}
	conv()
}

func main() {
	if len(os.Args) > 1 {
		cfg.Arg0, _ = strconv.Atoi(os.Args[1])
	}
	fmt.Printf("%+v\n", cfg)
	cfg.LocalDir = fmt.Sprintf(cfg.FormatDir, cfg.Arg0)
	chdir()
	for i := cfg.Arg1A[0]; i <= cfg.Arg1A[1]; i++ {
		cfg.File = fmt.Sprintf(cfg.FormatFile, i)
		cfg.URL = fmt.Sprintf(cfg.FormatURL, cfg.Host, cfg.LocalDir, cfg.File+cfg.InExt)
		download()
	}
	e := fmt.Sprintf("a%02d", cfg.Arg0)
	cfg.Arg1E = append(cfg.Arg1E, e)
	for _, e := range cfg.Arg1E {
		cfg.File = e
		cfg.URL = fmt.Sprintf(cfg.FormatURL, cfg.Host, cfg.LocalDir, cfg.File+cfg.InExt)
		download()
	}
}
