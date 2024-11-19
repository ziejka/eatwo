package handlers

import (
	"eatwo/views/components"
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

func renderHTMX(c echo.Context, statusCode int, t templ.Component, isAuthenticated bool) error {
	if c.Request().Header.Get("HX-Request") != "true" {
		t = layouts.Base(isAuthenticated, t)
	}
	return render(c, statusCode, t)
}

func renderError(c echo.Context, statusCode int, message string) error {
	c.Response().Header().Set("HX-Retarget", "#error-message")
	c.Response().Header().Set("HX-Reswap", "outerHTML")
	return render(c, statusCode, components.ErrorMsg(message))
}
