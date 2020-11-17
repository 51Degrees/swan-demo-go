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

package demo

import "html/template"

var marTemplate = newHTMLTemplate("mar", `
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>{{ .Title }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="icon" href="data:;base64,=">
    <style>
        body {
            margin: 0;
			font-family: Verdana,sans-serif;
        }
        header {
		    background-color: {{ .BackgroundColor }};
            line-height: 70px;
            vertical-align: middle;
            border-bottom: solid black 2px;
            padding: 5px 20px;
            margin: 0;
            position: sticky;
            top: 0;
        }
        header h1 {
            display: inline;
        }
        main {
            display: flex;
            flex-wrap: wrap;
        }
        main section {
            flex: 0 0 33.3333%;
        }
        main section h3, main section p {
            padding: 0 1em;
        }
        main section ul {
            list-style: none;
        }
        main section ul li {
            margin: 1em 0;
        }
        footer {
            padding: 5px 20px;
            position: sticky;
            bottom: 0;
            left: 0;
            right: 0;
            background-color: white;
            border-top: solid black 2px;
        }
        footer ul {
            list-style: none;
            padding: 0;
        }
        footer ul li {
            display: inline;
        }
    </style>
    <script>
        {{ if ne .JSON "" }}
        var preferences = {{ .JSON }};
        console.log(preferences);
        {{ end }}
    </script>
</head>
<body>
    <header>
        <h1>{{ .Title }}</h1>
    </header>
    <main>
    <section>
        <h3>Find out more about our great products.</h3>
        <ul>

        </ul>
    </section>           
    </main>
    <footer>
        <ul>
            {{ range $val := .Results }}
            <li>{{ $val.Key }} : {{ $val.Value }} | </li>
            {{ end }}
            <li><a href="{{ .SWANURL }}">Privacy Preferences</a></li>
        </ul>
    </footer>   
</body>
</html>`)

var pubTemplate = newHTMLTemplate("pub", `
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>Publisher | {{ .Title }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="icon" href="data:;base64,=">
    <style>
        body {
            margin: 0;
			font-family: Verdana,sans-serif;
        }
        header {
		    background-color: {{ .BackgroundColor }};
            line-height: 70px;
            vertical-align: middle;
            padding: 5px 20px;
            margin: 0 auto;
            max-width: 57rem;
        }
        header h1 {
            display: inline;
            font-size: 1.5em;
        }
        main {
            margin: 0 auto;
            max-width: 57rem;
        }
        main section {
            margin: 4em auto;
            display: block;
        }
        main section h3, main section p, main section pre, main section div {
            margin: 1em;
        }
        main section pre, main section span {
            background-color: lightgray;
            white-space: break-spaces;
            word-break: break-all;
            font-family: monospace;
        }
        main section pre {
            display: block;
            padding: 0.5em;
        }
        main section span {
            display: inline;
        }
        main section div {
            padding: 0.5em;
        }
        main section div .button {
            border-radius: 0.5em;
            background-color: lightblue;
            padding: 0.5em;
            border: lightgrey solid 1px;
        }
        main section ul {
            list-style: none;
        }
        main section ul li {
            margin: 1em 0;
        }
        footer {
            position: sticky;
            bottom: 0;
            left: 0;
            right: 0;
            background-color: white;
            border-top: solid black 2px;
        }
        footer ul {
            list-style: none;
            text-align: center;
            padding: 0.5em;
            margin: 0;
        }
        footer ul li {
            display: inline;
        }
        .logos {
            list-style: none;
            display: flex;
            padding: 0;
            align-items: center;
            flex-flow: row wrap;
            justify-content: center;
        }
        .logos li {
            display: block;
            float: left;
            margin: 1em;
        }
        .logos li img {
            width: 96px;
        }
    </style>
    <script>
        {{ if ne .JSON "" }}
        var preferences = {{ .JSON }};
        console.log(preferences);
        {{ end }}
        function verify(e, u) {
            fetch(u)
                .then(response => response.json())
                .then(data => {
                    if (data.valid == false) {
                        e.style.backgroundColor = "red";
                        e.innerText = "Bad";
                    } else {
                        e.style.backgroundColor = "green";
                        e.innerText = "Good";
                    }
                }).catch(() => e.style.backgroundColor = "red");
        }
        function creator(e, u) {
            fetch(u)
                .then(response => response.json())
                .then(data => {
                    e.innerText = data["name"];
                });
        }        
        function publicKey(e, u) {
            fetch(u)
                .then(response => response.json())
                .then(data => {
                    e.innerText = data["public-key"];
                });
        }    
        function text(e, u) {
            fetch(u)
                .then(response => response.text())
                .then(data => {
                    e.innerText = data;
                });
        }         
    </script>
</head>
<body>
    <header>
        <h1>Publisher: {{ .Title }}</h1>
    </header>
    <main>
    <section>
        <h3>Welcome to SWAN : the future of open web advertising</h3>
        <p>By enabling us to set a secure, privacy-by-design ID, we and other SWAN supporters promise to respect your privacy choices. The SWAN network is a privacy-by-design method of enabling personalized cross-publisher experiences on all publishers that adopt it.</p>
        <ul>
        <li><strong>People</strong> : enhanced transparency, persistent choices (no more consent fatigue) and "right to be forgotten".</li>
        <li><strong>Publishers</strong> : effective engagement, optimized advertising yield and accountable auditing to detect misappropriation.</li>
        <li><strong>Marketers</strong> : optimize cross publisher effectiveness, ensure you get what you pay for.</li>
        </ul>
        <p>Multiple implementors of SWAN open source form a dencentralized network. There's no single point of failure.</p>
        <p>Read on for a brief introduction to SWAN. To find out more go <a href="https://github.com/51degrees/swan">here</a>.</p>
    </section>
    <section>
        <h3>SWAN supporters</h3>
        <ul class="logos">
            <li><img src="//51degrees.com/img/logo.png"></li>
            <li><img src="//zetaglobal.com/wp-content/uploads/2017/12/Top_Logo@2x.png"></li>
            <li><img src="//www.liveintent.com/assets/img/brand-assets/LiveIntentLogo-Horiz-Orange.png"></li>
        </ul>
    </section>
    {{ if .CBID }}
    <section>
        <h3>Common Browser ID (CBID)</h3>
        <p>SWAN provides a Common Browser ID that you can easily reset at any time. Here's the SWAN CBID for this browser.<p>
        <pre>{{ .CBID.AsOWID.PayloadAsString }}</pre>
        <p>SWAN secures your ID to ensure you can have an accountable audit log. Here's the secured version:<p>
        <pre>{{ .CBID.Value }}</pre>
        <p>Anyone can confirm that this ID was created by <span><script>creator(document.scripts[document.scripts.length - 1].parentNode, '{{ .CBID.CreatorURL }}');</script></span> using this link.</p>
        <pre>{{ .CBID.VerifyURL }}</pre>
        <p>Go on. Tap the following button to check it's good.</p>
        <div><a class="button" onclick="verify(this, '{{ .CBID.VerifyURL }}')">Verify</a></div>
        <p>This shows that the domain <span>{{ .CBID.AsOWID.Domain }}</span> generated this ID on <span>{{ .CBID.AsOWID.Date }}</span>.</p>
        <p>The domain <span>{{ .CBID.AsOWID.Domain }}</span> used the following signature.</p>
        <pre>{{ .CBID.AsOWID.Signature }}</pre>
        <p>Because your online experience matters this publisher uses their public signing key so anyone can verify in microseconds.</p>
        <pre><script>publicKey(document.scripts[document.scripts.length - 1].parentNode, '{{ .CBID.CreatorURL }}');</script></pre>
    </section>
    {{ end }}
    {{ if .UUID }}
    <section>
        <h3>Signed-in ID (SID)</h3>
        <p>If you wish to preserve your preferences across multiple browsers or devices, you can use SWAN to share your signed-in ID (SID). This relies on hashing a validated email you provide to register at this site. Here's the SID SWAN generated from whatever you entered.</p>
        <pre>{{ .UUID.AsOWID.PayloadAsString }}</pre>
        <p>Just like CBID it's secured to make it verifiable. Here's the longer version.</p>
        <pre>{{ .UUID.Value }}</pre>
        <p>When all of this is decoded and verified it looks like this.</p>
        <pre><script>text(document.scripts[document.scripts.length - 1].parentNode, '{{ .UUID.DecodeAndVerifyURL }}');</script></pre>
        <p>SID and CBID are all implemented in SWAN using the Open Web ID schema. It's open source and your can find out more <a href="https://github.com/51degrees/owid">here</a>.</p>
    </section>
    {{ end }}
    {{ if .Allow }}
    <section>
        <h3>Preferences</h3>
        <p>Responsible addressable marketing requires both respecting your preferences and providing you transparency as to which organizations were involved in delivering you personalized content. SWAN has recorded your personalization preferences as.</p>
        <pre>{{ .Allow.AsOWID.PayloadAsString }}</pre>
        <p>Just like your Common Browser ID, we secure your preferences too. Your preference token is:</p>
        <pre>{{ .Allow.Value }}</pre>
        <p>You can change your preferences any time <a href="{{ .SWANURL }}">here</a>.</p>
        <p>If you want to only temporarily change your preference, you can using a new incognito or private browsing tab.</p>
    </section>
    {{ end }}
    <section>
        <h3>Improving the Web</h3>
        <p>We believe you deserve not only transparency and control, but an auditable view into the organizations involved in displaying the content on this publisher.  We need to rely on SWAN given recent announcements by browsers owned by the largest US publishers who intend to interfere with how we and other smaller publishers operate our business.</p>
        <p>To provide you access to our website, we rely on advertising paid by marketers. In exchange, they need to measure and optimize their advertising as easily as they can within the Walled Gardens. However, marketers do not need to know your offline identity and SWAN members agree to keep your offline identity distinct from your digital activity.</p>
        <p>Like the World Wide Web, SWAN is operated by an open market of hosts that do not want or have a central controller. Accordingly, there isn't a single SWAN domain, but many of them. To speed up your online experience, every browser is assigned a home domain. This reduces the number of times publishers need to ask multiple SWAN domains for your information.</p>
    </section>
    <section>
        <h3>Find out more about the open source projects used in this demo.</h3>
        <ul>
            <li><a href="https://github.com/51degrees/swift">SWIFT</a> Shared Web InFormaTion is a browser-agnostic method to share information across web domains.</li>
            <li><a href="https://github.com/51degrees/owid">OWID</a> Open Web ID (OWID) is a privacy-by-design schema for ID.</li>
            <li><a href="https://github.com/51degrees/swan">SWAN</a> Shared Web Accountable Network (SWAN) brings it all together to support digital marketing use cases.</li>
        </ul>
    </section>
    <section>
        <h3>Visit these other domains</h3>
        <p>Just visit these other domains that are part of SWAN to see how the same data is shared across multiple domains.</p>
        <ul>
            {{ range $val := .Pubs }}
            <li><a href="//{{ $val }}">{{ $val }}</a></li>
            {{ end }}
        </ul>
    </section>
    </main>
    <footer>
        <ul>
            {{ range $val := .Results }}
            <li>{{ $val.Key }} : {{ $val.AsOWID.PayloadAsString }}</li>
            {{ end }}
            <li><a href="{{ .SWANURL }}">Privacy Preferences</a></li>
        </ul>
    </footer>
</body>
</html>`)

func newHTMLTemplate(n string, h string) *template.Template {
	return template.Must(template.New(n).Parse(h))
}