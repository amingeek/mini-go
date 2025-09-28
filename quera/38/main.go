package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TemplateRender struct {
	templates *template.Template
}

type Author struct {
	FirstName string
	LastName  string
	Age       int
}

func (t *TemplateRender) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func aboutHandler(c echo.Context) error {
	jkRowling := Author{
		FirstName: "Joanne",
		LastName:  "Rowling",
		Age:       58,
	}

	data := map[string]interface{}{
		"author": jkRowling,
		"books": []string{"Harry Potter and the Sorcerer's Stone", "Harry Potter and the Chamber of Secrets",
			"Harry Potter and the Prisoner of Azkaban", "Harry Potter and the Goblet of Fire", "Harry Potter and the Order of the Phoenix",
			"Harry Potter and the Half-Blood Prince", "Harry Potter and the Deathly Hallows",
		},
		"isFamous": true,
	}
	return c.Render(http.StatusOK, "about.html", data)
}

func main() {
	e := echo.New()
	e.Renderer = &TemplateRender{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.GET("/about", aboutHandler)

	log.Fatal(e.Start(":8080"))
}
