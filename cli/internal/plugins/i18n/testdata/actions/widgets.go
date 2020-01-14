package actions

import (
	"github.com/gobuffalo/buffalo"
)

type WidgetsResource struct{}

func (v WidgetsResource) List(c buffalo.Context) error {
	return nil
}

func (v WidgetsResource) Show(c buffalo.Context) error {
	return nil
}

func (v WidgetsResource) Create(c buffalo.Context) error {
	msg := "CREATED!"
	c.Flash().Add("success", msg)
	return nil
}

func (v WidgetsResource) Edit(c buffalo.Context) error {
	return nil
}

func (v WidgetsResource) Update(c buffalo.Context) error {
	c.Flash().Add("success", "Updated!")
	return nil
}

func (v WidgetsResource) Destroy(c buffalo.Context) error {
	c.Flash().Add("success", "Destroyed!")
	return nil
}
