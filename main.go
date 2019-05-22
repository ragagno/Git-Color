package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const UrlGithub = "https://raw.githubusercontent.com/github/linguist/master/lib/linguist/languages.yml"
const UrlGitlab = "https://gitlab.com/gitlab-org/linguist/raw/master/lib/linguist/languages.yml"

const ArgsLongGithub = "--github"
const ArgsLongGitlab = "--gitlab"

type Language struct {
	name  string
	color string
}

func main() {
	var languages []Language

	var github = false
	var gitlab = false

	for _, element := range os.Args[1:] {
		switch element {
		case ArgsLongGithub:
			if github || gitlab {
				_, _ = fmt.Fprintf(os.Stderr, "remote already defined")
				os.Exit(-1)
			}

			github = true
		case ArgsLongGitlab:
			if github || gitlab {
				_, _ = fmt.Fprintf(os.Stderr, "remote already defined")
				os.Exit(-1)
			}

			gitlab = true
		default:
			_, _ = fmt.Fprintf(os.Stderr, "unknown argument: %s", element)
			os.Exit(-1)
		}
	}

	var url string

	if github {
		url = UrlGithub
	} else if gitlab {
		url = UrlGitlab
	} else {
		url = UrlGitlab
	}

	if reader, err := getYaml(url); err == nil {
		var currentLanguage string

		for {
			if line, _, err := reader.ReadLine(); err == nil {
				if line != nil && len(line) > 0 && line[0] != '#' && line[0] != '-' {
					if line[0] == ' ' {
						if strings.HasPrefix(string(line), "  color: ") {
							languages = append(languages, Language{
								currentLanguage,
								strings.Trim(strings.TrimPrefix(string(line), "  color: "), "\""),
							})
						}
					} else {
						currentLanguage = strings.TrimSuffix(string(line), ":")
					}
				}
			} else {
				if err == io.EOF {
					break
				} else {
					_, _ = fmt.Fprintf(os.Stderr, "unknown error reading body: %s", err.Error())
					os.Exit(-2)
				}
			}
		}
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "could not fetch body: %s", err.Error())
		os.Exit(-2)
	}

	for _, language := range languages {
		fmt.Printf("\"%s\": %s,\n", strings.ToLower(strings.Replace(language.name, " ", "-", -1)), strings.ToUpper(language.color))
	}
}

func getYaml(url string) (*bufio.Reader, error) {
	if response, err := http.Get(url); err == nil {
		if response != nil {
			return bufio.NewReader(response.Body), nil
		} else {
			return nil, fmt.Errorf("nil response")
		}
	} else {
		return nil, err
	}
}
