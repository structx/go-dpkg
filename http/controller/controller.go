// Package controller exposed http controllers
package controller

import (
	"github.com/labstack/echo/v4"
)

// Controller exposed http controller on service handlers
type Controller interface {
	RegisterRoutesV1(g *echo.Group)
}

// RootController exposed http controller on root handler
type RootController interface {
	RegisterRoutesV0(g *echo.Echo)
}
