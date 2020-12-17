/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/spf13/cobra"
)

type Bookmark struct {
	URL  string
	Tags []string
}

type importOptions struct {
	File string
}

// importCmd represents the import command
func NewImportCmd() *cobra.Command {
	options := importOptions{}
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Hatena bookmark import",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImport(cmd, options)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&options.File, "file", "f", "", "bookmark")
	return cmd
}

func runImport(cmd *cobra.Command, o importOptions) (err error) {
	bs, err := readBookmark(o.File)
	if err != nil {
		return err
	}

	oauthClient := &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  os.Getenv("HATENA_OAUTH_KEY"),
			Secret: os.Getenv("HATENA_OAUTH_SECRET"),
		},
		TemporaryCredentialRequestURI: "https://www.hatena.com/oauth/initiate",
		ResourceOwnerAuthorizationURI: "https://www.hatena.com/oauth/authorize",
		TokenRequestURI:               "https://www.hatena.com/oauth/token",
	}

	scope := url.Values{"scope": {"read_public,write_public"}}
	tempCredentials, err := oauthClient.RequestTemporaryCredentials(nil, "oob", scope)
	if err != nil {
		log.Fatal("RequestTemporaryCredentials:", err)
	}

	u := oauthClient.AuthorizationURL(tempCredentials, nil)
	fmt.Printf("1. Go to %s\n2. Authorize the application\n3. Enter verification code:\n", u)

	var code string
	fmt.Scanln(&code)

	tokenCard, _, err := oauthClient.RequestToken(nil, tempCredentials, code)
	if err != nil {
		log.Fatal("RequestToken:", err)
	}

	accessToken := oauth.Credentials{
		Token:  tokenCard.Token,
		Secret: tokenCard.Secret,
	}

	for i := len(bs) - 1; i >= 0; i-- {
		b := bs[i]
		func() {
			form := url.Values{}
			form.Set("url", b.URL)
			form["tags"] = b.Tags
			response, err := oauthClient.Post(nil, &accessToken, "https://bookmark.hatenaapis.com/rest/1/my/bookmark", form)
			if err != nil {
				cmd.PrintErrf("Failed add bookmark: %s\n", b.URL)
				return
			}
			defer response.Body.Close()
		}()

	}

	return nil
}

func readBookmark(file string) ([]*Bookmark, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(data)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	var bookmarks []*Bookmark
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		b := &Bookmark{}
		if url, ok := s.Attr("href"); ok {
			b.URL = url
		}
		if tags, ok := s.Attr("tags"); ok {
			b.Tags = strings.Split(tags, ",")
		}
		bookmarks = append(bookmarks, b)
	})

	return bookmarks, nil
}
