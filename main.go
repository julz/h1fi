package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gernest/front"
	"gopkg.in/yaml.v2"
)

func main() {
	path := flag.String("path", "docs", "")
	flag.Parse()

	m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)

	if err := filepath.WalkDir(*path, func(path string, _ fs.DirEntry, _ error) error {
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		f, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		s := string(f)
		if !strings.HasPrefix(s, "---\n") {
			return nil
		}

		s = strings.TrimPrefix(s, "---\n")
		split := strings.SplitN(s, "---\n", 2)
		front, body := split[0], split[1]

		header := make(map[string]interface{})
		if err := yaml.Unmarshal([]byte(front), &header); err != nil {
			return err
		}

		if _, ok := header["title"]; !ok {
			return nil
		}

		if strings.TrimSpace(body) == "" {
			return nil
		}

		if strings.HasPrefix(strings.TrimLeft(body, "\n\r"), "# ") {
			log.Println("skip ", path)
			return nil
		}

		w, err := os.OpenFile(path, os.O_WRONLY, 0)
		if err != nil {
			return err
		}

		fmt.Fprint(w, "---\n")
		fmt.Fprint(w, front)
		fmt.Fprint(w, "---\n\n")
		fmt.Fprintf(w, "# %v\n", header["title"])
		fmt.Fprint(w, body)

		return nil
	}); err != nil {
		log.Fatal("FAILED:", err)
	}
}
