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

package marketer

import (
	"bytes"
	"common"
	"fmt"
	"html/template"
	"owid"
	"strings"
	"swan"
)

// MarketerModel used with HTML templates.
type MarketerModel struct {
	common.PageModel
	offer *owid.Node // The offer and tree associated with the page
}

// Stop returns true if the request includes the key Stop to indicate that the
// advert should no longer be displayed.
func (m *MarketerModel) Stop() bool {
	for k := range m.Request.Form {
		if k == "stop" {
			return true
		}
	}
	return false
}

// TreeAsJSON return the transaction as JSON.
func (m *MarketerModel) TreeAsJSON() (template.HTML, error) {
	if m.offer == nil {
		return template.HTML("<p>Advert not source of request.</p>"), nil
	}

	b, err := m.offer.AsJSON()
	if err != nil {
		return "", err
	}

	var html bytes.Buffer
	html.WriteString("<p style=\"word-break:break-all\">")
	html.WriteString(string(b))
	html.WriteString("</p>")
	return template.HTML(html.String()), nil
}

// OfferID returns the offer ID string for the advert.
func (m *MarketerModel) OfferID() (string, error) {
	if m.offer == nil {
		return "", nil
	}
	o, err := m.offer.GetOWID()
	if err != nil {
		return "", err
	}
	return o.AsString(), nil
}

// OfferIDUnpacked returns the unpacked Offer ID
func (m *MarketerModel) OfferIDUnpacked() (template.HTML, error) {
	if m.offer == nil {
		return template.HTML("<p>Advert not source of request.</p>"), nil
	}
	o, err := m.offer.GetOWID()
	if err != nil {
		return template.HTML("<p>" + err.Error() + "</p>"), nil
	}
	s, err := swan.OfferFromOWID(o)
	if err != nil {
		return template.HTML("<p>" + err.Error() + "</p>"), nil
	}

	var html bytes.Buffer
	html.WriteString("<table class=\"table\">")
	html.WriteString("<thead><tr>")
	html.WriteString("<th>Field</th>")
	html.WriteString("<th>Value</th>")
	html.WriteString("</tr></thead><tbody>")
	html.WriteString(fmt.Sprintf(
		"<tr><td>Version</td><td>%d</td></tr>",
		o.Version))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Domain</td><td>%s</td></tr>",
		o.Domain))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Signature</td><td style=\"word-break:break-all\">%s</td></tr>",
		convertToString(o.Signature)))
	html.WriteString(fmt.Sprintf(
		"<tr><td>CBID</td><td>%s</td></tr>",
		s.CBIDAsString()))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Allow</td><td>%s</td></tr>",
		s.PreferencesAsString()))
	html.WriteString(fmt.Sprintf(
		"<tr><td>SID</td><td>%s</td></tr>",
		s.SIDAsString()))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Pub. domain</td><td>%s</td></tr>",
		s.PubDomain))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Placement</td><td>%s</td></tr>",
		s.Placement))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Unique</td><td style=\"word-break:break-all\">%s</td></tr>",
		convertToString(s.UUID)))
	html.WriteString(fmt.Sprintf(
		"<tr><td>Stopped Ads.</td><td style=\"word-break:break-all\">%s</td></tr>",
		strings.Join(s.StoppedAsArray(), ",")))
	htmlAddFooter(&html)
	return template.HTML(html.String()), nil
}

// AuditWinnerHTML returns the audit information from the bid used in the advert
// that resulted in the request to this page.
func (m *MarketerModel) AuditWinnerHTML() (template.HTML, error) {
	var html bytes.Buffer
	if m.offer == nil {
		return template.HTML("<p>Advert not source of request.</p>"), nil
	}

	w, err := swan.WinningNode(m.offer)
	if err != nil {
		return "", nil
	}

	if err != nil {
		return "", nil
	}
	htmlAddHeader(&html)
	err = appendParents(&html, w)
	if err != nil {
		return "", err
	}
	htmlAddFooter(&html)
	return template.HTML(html.String()), nil
}

// AuditFullHTML returns the audit information from the bid used in the advert
// that resulted in the request to this page.
func (m *MarketerModel) AuditFullHTML() (template.HTML, error) {
	if m.offer == nil {
		return template.HTML("<p>Advert not source of request.</p>"), nil
	}

	w, err := swan.WinningNode(m.offer)
	if err != nil {
		return "", err
	}

	var html bytes.Buffer
	htmlAddHeader(&html)
	err = appendOWIDAndChildren(&html, m.offer, w, 0)
	if err != nil {
		return template.HTML("<p>" + err.Error() + "</p>"), nil
	}
	htmlAddFooter(&html)
	return template.HTML(html.String()), nil
}

func convertToString(b []byte) string {
	return fmt.Sprintf("%x", b)
}

func htmlAddHeader(html *bytes.Buffer) {
	html.WriteString("<table class=\"table\">\r\n")
	html.WriteString("<thead>\r\n<tr>\r\n")
	html.WriteString("<th>Organization</th>\r\n")
	html.WriteString("<th>Audit Result</th>\r\n")
	html.WriteString("<th>\r\n</th>\r\n")
	html.WriteString("<th>\r\n</th>\r\n")
	html.WriteString("</tr>\r\n</thead>\r\n<tbody>\r\n")
}

func htmlAddFooter(html *bytes.Buffer) {
	html.WriteString("</tbody>\r\n</table>\r\n")
}

func appendParents(html *bytes.Buffer, w *owid.Node) error {
	var n []*owid.Node
	p := w
	for p != nil {
		n = append(n, p)
		p = p.GetParent()
	}
	i := len(n) - 1
	for i >= 0 {
		err := appendHTML(html, w, n[i], 0)
		if err != nil {
			return err
		}
		i--
	}
	return nil
}

func appendOWIDAndChildren(
	html *bytes.Buffer,
	o *owid.Node,
	w *owid.Node,
	level int) error {
	appendHTML(html, w, o, level)
	if len(o.Children) > 0 {
		for _, c := range o.Children {
			err := appendOWIDAndChildren(html, c, w, level+1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func appendHTML(
	html *bytes.Buffer,
	w *owid.Node,
	o *owid.Node,
	level int) error {

	s, err := swan.FromNode(o)
	if err != nil {
		return err
	}

	html.WriteString("<tr>\r\n")
	html.WriteString(fmt.Sprintf(
		"<td style=\"padding-left:%dem;\" class=\"text-left\">\r\n"+
			"<script>new owid().appendName(document.currentScript.parentNode,\"%s\")</script></td>\r\n",
		level,
		o.GetOWIDAsString()))

	var r string
	if o.GetRoot() != o {
		r = o.GetRoot().GetOWIDAsString()
	}

	html.WriteString(fmt.Sprintf(
		"<td style=\"text-align:center;\">\r\n"+
			"<script>new owid().appendAuditMark(document.currentScript.parentNode,\"%s\",\"%s\");</script>\r\n"+
			"<noscript>JavaScript needed to audit</noscript></td>\r\n",
		r,
		o.GetOWIDAsString()))
	if w == o {
		html.WriteString("<td>\r\n<img style=\"width:32px\" src=\"noun_rosette_470370.svg\">\r\n</td>\r\n")
	} else {
		f, fok := s.(*swan.Failed)
		_, bok := s.(*swan.Bid)
		if fok {
			html.WriteString(fmt.Sprintf("<td style=\"color:lightpink\">\r\n%s&nbsp;%s</td>\r\n",
				f.Host,
				f.Error))
		} else if bok {
			html.WriteString("<td>\r\n<img style=\"width:32px\" src=\"noun_movie ticket_1807397.svg\">\r\n</td>\r\n")
		} else {
			html.WriteString("<td>\r\n</td>\r\n")
		}
	}

	if o != nil && r != "" {
		html.WriteString(fmt.Sprintf(
			"<td style=\"text-align:center;\">\r\n"+
				"<script>new owid().appendComplaintEmail(document.currentScript.parentNode,\"%s\",\"%s\", \"noun_complaint_376466.svg\");</script>\r\n"+
				"<noscript>JavaScript needed to audit</noscript></td>\r\n",
			r,
			o.GetOWIDAsString()))
	} else {
		html.WriteString("<td>\r\n</td>\r\n")
	}
	html.WriteString("</tr>\r\n")
	return nil
}
