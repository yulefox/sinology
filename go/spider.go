package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
	Files     [][]string
	Root      bool
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
	ct = regexp.MustCompile(".*#### 二十四史.*\n").ReplaceAllString(ct, "")
	ct = regexp.MustCompile(`.*\[Menu\]\(#nav\).*\n`).ReplaceAllString(ct, "")
	ct = regexp.MustCompile(`\n.*\(\D+.*\.md\).*`).ReplaceAllString(ct, "")

	if cfg.Root {
		md := ct
		// trim leading/trailing whitespace
		md = regexp.MustCompile(`(\[.*\]\(\d+.*\.md\))`).ReplaceAllString(md, "\n* ${1}\n")
		md = regexp.MustCompile(`\n[\* ]+\n`).ReplaceAllString(md, "\n")
		md = regexp.MustCompile(`\n[^#\*].*\n`).ReplaceAllString(md, "\n")
		md = regexp.MustCompile(`\n[\t ]+|[\t ]+\n|\n[\t 　]+\n`).ReplaceAllString(md, "\n")
		md = regexp.MustCompile(`\n{3,}`).ReplaceAllString(md, "\n\n")
		ioutil.WriteFile("SUMMARY.md", []byte(md), 0644)
	}

	// trim leading/trailing whitespace
	ct = regexp.MustCompile(`\n.*\(.*\.md\).*`).ReplaceAllString(ct, "")
	ct = regexp.MustCompile(`\n[\t ]+|[\t ]+\n|\n[\t 　]+\n`).ReplaceAllString(ct, "\n")
	ct = regexp.MustCompile(`\n{3,}`).ReplaceAllString(ct, "\n\n")
	if cfg.Root {
		ct = regexp.MustCompile(`\n[\* ]+\n+\* ]+\n`).ReplaceAllString(ct, "\n")
		ioutil.WriteFile("README.md", []byte(ct), 0644)
	}

	return ct
}

func conv() {
	inFile := cfg.File + cfg.InExt
	outFile := cfg.File + cfg.OutExt
	html, _ := ioutil.ReadFile(inFile)
	if cfg.Root {
		exp := `href="([\w]+)\` + cfg.InExt + `"`
		cfg.Files = regexp.MustCompile(exp).FindAllStringSubmatch(string(html), -1)
	}
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

	for _, v := range cfg.Files {
		cfg.File = v[1]
		cfg.Root = false
		cfg.URL = fmt.Sprintf(cfg.FormatURL, cfg.Host, cfg.LocalDir, cfg.File+cfg.InExt)
		download()
	}
}

func main() {
	n := len(os.Args)
	if n < 3 {
		log.Fatal("spider cate name [ext]")
	}
	if n > 3 {
		cfg.InExt = "." + os.Args[3]
	}
	cfg.LocalDir = os.Args[1]
	cfg.File = os.Args[2]
	cfg.Root = true
	cfg.URL = fmt.Sprintf(cfg.FormatURL, cfg.Host, cfg.LocalDir, cfg.File+cfg.InExt)
	chdir()
	toc()
}
