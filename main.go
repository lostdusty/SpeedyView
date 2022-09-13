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

	"github.com/akamensky/argparse"
	"github.com/jchv/go-webview-selector"
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

func createView() {

	window := webview.New(D) //Set the value from the debug flag (defaults to false)
	if window == nil {
		log.Fatalln("Failed to load webview.")
	}
	defer window.Destroy()
	window.SetTitle(T)      //Set the title value (default to SpeedyView)
	window.SetSize(W, H, S) //Now define Width, Height of the window, and also the sizing hint
	window.Navigate(U)      //And finally, the url
	window.Run()            //Starts the window

	//TODO: Install function
}
