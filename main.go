/*
Copyright © 2022 Princess Mortix hi@princessmortix.link
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/akamensky/argparse"
	webview "github.com/jchv/go-webview2"
	"github.com/jchv/go-webview2/webviewloader"
	"github.com/melbahja/got"
)

var (
	U, T string
	H, W int
	S    webview.Hint
	D, I bool
)

func main() {
	// Create new parser object
	parser := argparse.NewParser("SpeedyView", "application for viewing web pages with speed. Uses Edge WebView2 under the hood.")
	// Create flags for the window configuration
	u := parser.String("u", "url", &argparse.Options{Required: true, Help: "Website url to open"})
	h := parser.Int("e", "height", &argparse.Options{Required: false, Default: 600, Help: "Window height"})
	w := parser.Int("w", "width", &argparse.Options{Required: false, Default: 800, Help: "Window width"})
	t := parser.String("t", "title", &argparse.Options{Required: false, Default: "SpeedyView", Help: "Title of the browser window"})
	s := parser.Selector("s", "sizing", []string{"none", "fixed", "min", "max"}, &argparse.Options{Required: false, Default: "fixed", Help: "Defines window sizing and resizing behavior. None specifies that width and height are default size\n\t\tFixed specifies that window size cannot be changed by the user\n\t\tMin specifies that width and height are minimum bounds\n\t\tMax specifies that width and height are maximum bounds"})
	d := parser.Flag("d", "debug", &argparse.Options{Required: false, Default: false, Help: "Makes the window debuggable"})
	i := parser.Flag("i", "install", &argparse.Options{Required: false, Default: false, Help: "Installs WebView2 if its not installed on the system."})
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	// Convert the pointers to global variables
	U = fmt.Sprint(*u)
	H = int(*h)
	W = int(*w)
	T = fmt.Sprint(*t)
	switch *s {
	case "none":
		S = webview.HintNone
	case "fixed":
		S = webview.HintFixed
	case "min":
		S = webview.HintMin
	case "max":
		S = webview.HintMax
	}
	D = bool(*d)
	I = bool(*i)

	//Call the window function to create the window
	createView()
}

func installView() {
	webviewVersion, err := webviewloader.GetInstalledVersion() //Check if the runtime is installed calling GetInstalledVersion()
	if err != nil {
		log.Println("WebView2 was not detected. SpeedyView will now install it") //If it fails to get the version, the application will:
		downloadWebview := got.New()
		err := downloadWebview.Download("https://go.microsoft.com/fwlink/p/?LinkId=2124703", "MicrosoftEdgeWebview2Setup.exe") //1: Download the runtime installer from Microsoft Website

		if err != nil {
			log.Fatalln("Failed to download the WebView runtime, SpeedyView will exit now.", err) //Exits with status code 1 if fails to download.
		}

		log.Println("Download has been completed, SpeedyView will now install the downloaded file.") //2: Install the runtime using /silent and /install arguments
		cmd := exec.Cmd{                                                                             //See https://docs.microsoft.com/en-us/microsoft-edge/webview2/concepts/distribution#online-only-deployment
			Path:        "MicrosoftEdgeWebview2Setup.exe",
			Args:        []string{"/silent", "/install"},
			Stdout:      os.Stdout,
			SysProcAttr: &syscall.SysProcAttr{HideWindow: true}, //This syscall attribute hides the console window that may appear.
		}
		err = cmd.Run() //Executes the command
		if err != nil {
			log.Fatalln("Failed to install WebView runtime, error returned by the installer:", err) //Exits with code 1 if fails to run the command.
		}
	}

	log.Println("Got WebView version", webviewVersion) //Print the WebView version.
}

func createView() {

	if I {
		installView()
	}

	window := webview.New(D) //Set the value from the debug flag (defaults to false)
	if window == nil {
		log.Fatalln("Failed to load webview.\nIf the issue persist, use --install flag.")
	}
	defer window.Destroy()
	window.SetTitle(T)      //Set the title value (default to SpeedyView)
	window.SetSize(W, H, S) //Now define Width, Height of the window, and also the sizing hint
	window.Navigate(U)      //And finally, the url
	window.Run()            //Starts the window

	//TODO: Install function
}
