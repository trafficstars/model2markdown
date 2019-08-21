package main

import (
	"flag"
	"fmt"
	"github.com/trafficstars/model2markdown"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func fatalSyntax() {
	_, _ = fmt.Fprintf(os.Stderr, "syntax: %v [options (see below)] <input file/dir [input file/dir [...]]>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(int(syscall.EINVAL))
}

func parseArgs(cfg *Config) {
	if cfg == nil {
		panic(`cfg == nil`)
	}

	outDir := flag.String(`output-directory`, `/tmp`, `the directory to put resulting markdown documents to`)

	flag.Parse()

	if flag.NArg() < 1 {
		fatalSyntax()
	}

	cfg.OutputDirectory = *outDir

	for _, filePath := range flag.Args() {
		stat, err := os.Lstat(filePath)
		if err != nil {
			log.Fatal(err)
			continue // won't be executed, actually
		}
		if !stat.IsDir() {
			cfg.Jobs = append(cfg.Jobs, Job{
				InputFile: filePath,
			})
		}

		err = filepath.Walk(filePath, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
				return err
			}
			if !strings.HasSuffix(filePath, `.go`) {
				return nil
			}

			cfg.Jobs = append(cfg.Jobs, Job{
				InputFile: filePath,
			})

			return nil
		})
		if err != nil {
			log.Fatal(err)
			continue // won't be executed, actually
		}
	}
}

func main() {
	var cfg Config
	parseArgs(&cfg)

	library := model2markdown.NewLibrary()

	for _, job := range cfg.Jobs {
		file, err := model2markdown.ParseFile(job.InputFile)
		if err != nil {
			log.Fatal(err)
			continue // won't be executed, actually
		}

		library.AddFile(file)
	}

	err := library.GenerateMarkdownsToDirectory(cfg.OutputDirectory)
	if err != nil {
		log.Fatal(err)
		return // won't be executed, actually
	}
}
