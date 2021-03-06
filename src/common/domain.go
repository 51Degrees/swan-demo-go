/* ****************************************************************************
 * Copyright 2020 51 Degrees Mobile Experts Limited (51degrees.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 * ***************************************************************************/

package common

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"owid"
	"path/filepath"
	"strings"
	"swan"
	"swift"
)

// Domain represents the information held in the domain configuration file
// commonly represented in the demo in config.json.
type Domain struct {
	Category            string // Category of the domain
	Name                string // Common name for the domain
	Bad                 bool   // True if this domain is a bad actor for the demo
	Host                string // The host name for the domain
	SwanMessage         string // Message if used with SWAN
	SwanBackgroundColor string // Background color if used with SWAN
	SwanMessageColor    string // Message text color if used with SWAN
	SwanProgressColor   string // Message progress color if used with SWAN
	// The domain of the access node used with SWAN (only set for CMPs)
	SWANAccessNode string
	SWANAccessKey  string // The access key to use when communicating with SWAN.
	// The domain of the CMP that will in turn access the SWAN Network via an Operator
	CMP       string
	Suppliers []string           // Suppliers used by the domain operator
	Adverts   []Advert           // Adverts the domain can serve
	Config    *Configuration     // Configuration for the server
	folder    string             // Location of the directory
	templates *template.Template // HTML templates
	OWID      *owid.Creator      // The OWID creator associated with the domain if any
	owidStore owid.Store         // The connection to the OWID store
	// The HTTP handler to use for this domain
	handler func(d *Domain, w http.ResponseWriter, r *http.Request)
}

// NewDomain creates a new instance of domain information from the file
// provided.
func NewDomain(c *Configuration, folder string) (*Domain, error) {
	var d Domain

	// Read the configuration for the folder provided.
	configFile, err := os.Open(filepath.Join(folder, "config.json"))
	defer configFile.Close()
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&d)

	// Set the private members.
	d.Config = c
	d.Host = filepath.Base(folder)
	d.folder = folder
	d.templates, err = d.parseHTML()
	if err != nil {
		return nil, err
	}
	d.owidStore = c.owid

	return &d, nil
}

// SetHandler adds a HTTP handler to the domain.
func (d *Domain) SetHandler(fn func(
	d *Domain,
	w http.ResponseWriter,
	r *http.Request)) {
	d.handler = fn
}

func (d *Domain) parseHTML() (*template.Template, error) {
	var t *template.Template
	files, err := ioutil.ReadDir(d.folder)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".html" {
			s, err := ioutil.ReadFile(filepath.Join(d.folder, file.Name()))
			if err != nil {
				return nil, err
			}
			if t == nil {
				t, err = template.New(file.Name()).Funcs(
					template.FuncMap{"role": infoRole}).Parse(
					removeHTMLWhiteSpace(string(s)))
				if err != nil {
					return nil, err
				}
			} else {
				t, err = t.New(file.Name()).Funcs(
					template.FuncMap{"role": infoRole}).Parse(
					removeHTMLWhiteSpace(string(s)))
				if err != nil {
					return nil, err
				}
			}

		}
	}
	return t, nil
}

// LookupHTML based on the templates available to the domain.
func (d *Domain) LookupHTML(p string) *template.Template {
	if d.templates == nil {
		return nil
	}

	// Try to find the template that relates to the file path.
	t := d.templates.Lookup(filepath.Base(p))

	// If no template can be found try finding one for the category of the
	// domain.
	if t == nil {
		t = d.templates.Lookup(strings.ToLower(d.Category) + ".html")
	}

	// Finally, if no template is found try the default one.
	if t == nil {
		t = d.templates.Lookup("default.html")
	}
	return t
}

func (d *Domain) setCommon(r *http.Request, q *url.Values) {

	// Set the access key
	q.Set("accessKey", d.SWANAccessKey)

	// Set the user interface title, message and colours.
	q.Set("title", "Your Preference Management")

	// Add the headers that are relevant to the home node calculation.
	swift.SetHomeNodeHeaders(r, q)
}

// CallSWANURL constructs a URL, gets the response, and then returns the
// response as a byte array. If an error occurs then an API error is returned.
func (d *Domain) CallSWANURL(
	action string,
	addParams func(*url.Values) error) ([]byte, *SWANError) {
	if d.SWANAccessNode == "" {
		return nil, &SWANError{fmt.Errorf(
			"Verify '%s' config.json for missing SWANAccessNode",
			d.Host), nil}
	}
	if d.SWANAccessKey == "" {
		return nil, &SWANError{fmt.Errorf(
			"Verify '%s' config.json for missing SWANAccessKey",
			d.Host), nil}
	}
	var u url.URL
	u.Scheme = d.Config.Scheme
	u.Host = d.SWANAccessNode
	u.Path = "/swan/api/v1/" + action
	q := u.Query()
	q.Set("accessKey", d.SWANAccessKey)
	err := addParams(&q)
	if err != nil {
		return nil, &SWANError{err, nil}
	}
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return nil, &SWANError{err, nil}
	}
	if res.StatusCode != http.StatusOK {
		return nil, NewSWANError(d.Config, res)
	}
	b, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return nil, &SWANError{e, nil}
	}
	return b, nil
}

// CreateSWANURL returns a URL from SWAN to pass to the web browser navigation.
func (d *Domain) CreateSWANURL(
	r *http.Request,
	returnURL string,
	action string,
	addParams func(*url.Values)) (string, *SWANError) {
	b, err := d.CallSWANURL(action, func(q *url.Values) error {
		d.setCommon(r, q)
		// If an explicit return URL was provided then use that. Otherwise use the
		// page for the current request.
		if returnURL != "" {
			q.Set("returnUrl", returnURL)
		} else {
			q.Set("returnUrl", getCurrentPage(d.Config, r))
		}

		// Add user interface parameters for the SWAN operation and the user
		// interface.
		if d.SwanMessage != "" {
			q.Set("message", d.SwanMessage)
		}
		if d.SwanBackgroundColor != "" {
			q.Set("backgroundColor", d.SwanBackgroundColor)
		}
		if d.SwanProgressColor != "" {
			q.Set("progressColor", d.SwanProgressColor)
		}
		if d.SwanMessageColor != "" {
			q.Set("messageColor", d.SwanMessageColor)
		}

		// Add any additional parameters needed by the action if a function was
		// provided.
		if addParams != nil {
			addParams(q)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// GetOWIDCreator returns the OWID creator from the OWID store for the the
// domain.
func (d *Domain) GetOWIDCreator() (*owid.Creator, error) {
	var err error
	if d.OWID == nil {
		d.OWID, err = d.owidStore.GetCreator(d.Host)
		if err != nil {
			return nil, err
		}
		if d.OWID == nil {
			return nil, fmt.Errorf(
				"Domain '%s' is not a registered OWID creator. Register the "+
					"domain for the SWAN demo using http[s]://%s/owid/register",
				d.Host,
				d.Host)
		}
	}
	return d.OWID, nil
}

func infoRole(s interface{}) string {
	_, fok := s.(*swan.Failed)
	_, bok := s.(*swan.Bid)
	_, eok := s.(*swan.Empty)
	_, ook := s.(*swan.Offer)
	if fok {
		return "Failed"
	}
	if bok {
		return "Bid"
	}
	if eok {
		return "Empty"
	}
	if ook {
		return "Offer"
	}
	return ""
}

// Removes white space from the HTML string provided whilst retaining valid
// HTML.
func removeHTMLWhiteSpace(h string) string {
	var sb strings.Builder
	for i, r := range h {

		// Only write out runes that are not control characters.
		if r != '\r' && r != '\n' && r != '\t' {

			// Only write this rune if the rune is not a space, or if it is a
			// space the preceding rune is not a space.
			if i == 0 || r != ' ' || h[i-1] != ' ' {
				sb.WriteRune(r)
			}
		}
	}
	return sb.String()
}
