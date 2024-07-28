package handlers

import (
	"eatwo/services"
	"eatwo/views/layouts"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// This custom render replaces Echo's echo.Context.render() with templ's templ.Component.render().
func render(c echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(c.Request().Context(), buf); err != nil {
		return err
	}

	return c.HTML(statusCode, buf.String())
}

func renderHTMX(c echo.Context, statusCode int, t templ.Component, claims *services.CustomClaims) error {
	if c.Request().Header.Get("HX-Request") != "true" {
		t = layouts.Base(claims, t)
	}
	return render(c, statusCode, t)
}
