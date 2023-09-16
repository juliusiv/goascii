package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"goascii/ascii_converter"
)

var (
	//go:embed templates/*
	resources embed.FS
	//go:embed css/build.css
	css embed.FS
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func handleUploadImage(c echo.Context) error {
	multipart_file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	file, err := multipart_file.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	tmpfile, err := os.CreateTemp("tmp", fmt.Sprintf("upload-*-%s", multipart_file.Filename))
	if err != nil {
		return err
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	// Copy
	if _, err = io.Copy(tmpfile, file); err != nil {
		return err
	}

	ascii := ascii_converter.ConvertToAscii(tmpfile)
	ascii_converter.PrintAscii(ascii)

	data := map[string][][]string{
		"Ascii": ascii,
	}
	return c.Render(http.StatusOK, "ascii-grid.html", data)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tmpls := &Template{
		templates: template.Must(template.ParseFS(resources, "templates/*")),
	}

	e := echo.New()
	e.Renderer = tmpls

	e.Use(middleware.Logger())

	e.Static("/css", "css")

	e.GET("/", func(c echo.Context) error {
		data := map[string]string{
			"Region": os.Getenv("FLY_REGION"),
		}
		return c.Render(http.StatusOK, "index.html", data)
	})

	e.POST("/upload-image", handleUploadImage)

	e.Logger.Fatal(e.Start(":8080"))
}
