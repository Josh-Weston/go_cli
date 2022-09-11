package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type content struct {
	Title    string
	Body     template.HTML
	FileName string
}

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		<title>{{ .Title }}</title>
	</head>
	<body>
<h3>Now previewing: {{ .FileName }}</h3>
{{ .Body }}
	</body>
</html>
`
)

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")
	flag.Parse()
	// if no file provided, show the usage text
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, *tFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename, tFname string, out io.Writer, skipPreview bool) error {
	// read all data and check for errors
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	htmlData, err := parseContent(input, tFname)
	if err != nil {
		return err
	}

	// Create temporary file for preview
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()
	out.Write([]byte(outName))

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}
	defer os.Remove(outName)
	return preview(outName)
}

func parseContent(input []byte, tFname string) ([]byte, error) {
	// parse the markdown file through blackfriday and bluemonday
	// to generate a valid and safe HTML
	output := blackfriday.Run(input)                     // creates the HTML
	body := bluemonday.UGCPolicy().SanitizeBytes(output) // sanitizes the HTML

	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	c := content{
		Title:    "Markdown Preview Tool",
		Body:     template.HTML(body), // convert the bytes to string of already sanitized
		FileName: "Default Template",
	}

	if tFname != "" || os.Getenv("TEMPLATE_FILE") != "" {
		var fname string
		// command flags take precedence over environment variables
		if tFname != "" {
			fname = tFname
		} else {
			fname = os.Getenv("TEMPLATE_FILE")
		}
		t, err = template.ParseFiles(fname)
		if err != nil {
			return nil, err
		}
		c.FileName = fname
	}

	var bb bytes.Buffer

	if err := t.Execute(&bb, c); err != nil {
		return nil, err
	}

	return bb.Bytes(), nil
}

func saveHTML(outFname string, data []byte) error {
	return os.WriteFile(outFname, data, 0644) // read and write by owner, but only readable by others
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	// Define executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	// Append filename to parameters slice
	cParams = append(cParams, fname)

	// Located executable in PATH
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// Open the file using the default program
	err = exec.Command(cPath, cParams...).Run()
	time.Sleep(2 * time.Second) // delay to give the browser time to open the file (we will learn how to use signals later in the book to avoid these dirty delays)
	return err

}
