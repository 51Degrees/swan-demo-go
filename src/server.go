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

package main

import (
	"demo"
	"fmt"
	"log"
	"net/http"
	"os"
)

func getPort() string {
	var port string
	if os.Getenv("HTTP_PLATFORM_PORT") != "" {
		// Get the port environment variable from Azure App Services.
		port = os.Getenv("HTTP_PLATFORM_PORT")
	} else if os.Getenv("PORT") != "" {
		// Get the port environment variable from Amazon Web Services.
		port = os.Getenv("PORT")
	} else {
		// If there is no environment variable use 5000, the default for AWS.
		port = "5000"
	}
	return port
}

func main() {
	var settingsFile string

	// Get the path to the settings file.
	if len(os.Args) >= 2 {
		settingsFile = os.Args[1]
	} else {
		settingsFile = "appsettings.json"
	}

	// Add the SWAN handlers.
	demo.AddHandlers(settingsFile)

	// Start the web server on the port provided.
	port := getPort()
	log.Printf("Listenning on port: %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
