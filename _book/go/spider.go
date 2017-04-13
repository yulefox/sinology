package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"

	"github.com/lunny/html2md"
)

var cfgBlob = []byte(`{
	"host": "xn--5rtnx620bw5s.tw",
	"fmt_url": "http://%s/%s/%s",
	"in_ext": ".htm",
	"out_ext": ".md"
}`)

var cfg config

type config struct {
	Host      string `json:"host"`
	FormatURL string `json:"fmt_url"`
	InExt     string `json:"in_ext"`
	OutExt    string `json:"out_ext"`
	LocalDir  string
	URL       string
	File      string
}

func init() {
	err := json.Unmarshal(cfgBlob, &cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	html2md.AddRule("span", &html2md.Rule{
		Patterns: []string{"span"},
		Replacement: func(innerHTML string, attrs []string) string {
			if len(attrs) > 1 {
				return "`" + attrs[1] + "`"
			}
			return ""
		},
	})
	html2md.AddConvert(link)
	html2md.AddConvert(cleanUp)
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
func link(ct string) string {
	re := regexp.MustCompile(`(\([^\)]*)\.htm[l]?\)`)
	ct = re.ReplaceAllString(ct, "${1}.md)")
	return ct
}

func cleanUp(ct string) string {
	ct = regexp.MustCompile(".*漢川草廬.*\n").ReplaceAllString(ct, "")
	ct = regexp.MustCompile(`.*\[Menu\]\(#nav\).*\n`).ReplaceAllString(ct, "")

	// trim leading/trailing whitespace
	ct = regexp.MustCompile("^[\t\r\n]+|[\t\r\n]+$").ReplaceAllString(ct, "")

	ct = regexp.MustCompile(`\n\{3,}`).ReplaceAllString(ct, "\n\n")
	return ct
}

func conv() {
	inFile := cfg.File + cfg.InExt
	outFile := cfg.File + cfg.OutExt
	html, _ := ioutil.ReadFile(inFile)
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

func toc() {
	download()

	outFile := cfg.File + cfg.OutExt
	md, _ := ioutil.ReadFile(outFile)
	files := regexp.MustCompile(`\(([\w]+)\.md\)`).FindAllStringSubmatch(string(md), -1)
	for _, v := range files {
		cfg.File = v[1]
		cfg.URL = fmt.Sprintf(cfg.FormatURL, cfg.Host, cfg.LocalDir, cfg.File+cfg.InExt)
		download()
	}
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("spider cate name")
	}
	cate := os.Args[1]
	name := os.Args[2]
	cfg.File = name
	cfg.LocalDir = path.Join(cate, name)
	cfg.URL = fmt.Sprintf(cfg.FormatURL, cfg.Host, cfg.LocalDir, cfg.File+cfg.InExt)
	chdir()
	toc()
}
