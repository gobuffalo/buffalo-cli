package actions

import (
	"github.com/gobuffalo/buffalo"
)

type GadgetsResource struct{}

func (v GadgetsResource) List(c buffalo.Context) error {
	return nil
}

func (v GadgetsResource) Show(c buffalo.Context) error {
	return nil
}

func (v GadgetsResource) Create(c buffalo.Context) error {
	msg := "CREATED!"
	c.Flash().Add("success", msg)
	return nil
}

func (v GadgetsResource) Edit(c buffalo.Context) error {
	return nil
}

func (v GadgetsResource) Update(c buffalo.Context) error {
	c.Flash().Add("success", "UPDATED!")
	return nil
}

func (v GadgetsResource) Destroy(c buffalo.Context) error {
	c.Flash().Add("success", "DESTROYED!")
	return nil
}
