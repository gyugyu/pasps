/*
Copyright © 2023 Yusuke Yagyu <yu@hoaxster.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package helper

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetClient(profile string) (*http.Client, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	profilePath := filepath.Join(home, ".gapy", profile)
	credentialsJson, err := os.ReadFile(filepath.Join(profilePath, "credentials.json"))
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(credentialsJson, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}

	tokenFile, err := os.Open(filepath.Join(profilePath, "token.json"))
	if err != nil {
		return nil, err
	}

	defer tokenFile.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(tokenFile).Decode(token)
	if err != nil {
		return nil, err
	}

	return config.Client(context.Background(), token), nil
}

func CreateValues(rows [][]string, padding int) [][]interface{} {
	var values [][]interface{}

	for i := 0; i < padding; i++ {
		values = append(values, []interface{}{})
	}

	for _, row := range rows {
		var r []interface{}
		for _, item := range row {
			r = append(r, item)
		}
		values = append(values, r)
	}
	return values
}
