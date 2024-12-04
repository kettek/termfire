package play

import (
	"github.com/kettek/termfire/messages"
	"github.com/rivo/tview"
)

type Container struct {
	items     []messages.ItemObject
	container *tview.Flex
	buttons   *tview.Flex
	listView  *tview.List
	onTrigger func(button string, object messages.ItemObject, index int)
}

func (c *Container) Init(title string, actions []string) {
	c.container = tview.NewFlex()
	c.container.SetDirection(tview.FlexRow)
	c.container.SetTitle(title)
	c.container.SetBorder(true)

	c.listView = tview.NewList()
	c.listView.ShowSecondaryText(false)
	c.container.AddItem(c.listView, 0, 1, false)

	c.buttons = tview.NewFlex()
	c.buttons.SetDirection(tview.FlexColumn)
	for _, action := range actions {
		btn := tview.NewButton(action)
		btn.SetSelectedFunc(func() {
			index := c.listView.GetCurrentItem()
			if index < 0 {
				return
			}
			if c.onTrigger != nil {
				c.onTrigger(action, c.items[index], index)
			}
		})
		c.buttons.AddItem(btn, 0, 1, false)
	}
	c.container.AddItem(c.buttons, 1, 1, false)
}

func (c *Container) Clear() {
	c.items = []messages.ItemObject{}
	c.listView.Clear()
}

func (c *Container) GetContainer() *tview.Flex {
	return c.container
}

func (c *Container) GetList() *tview.List {
	return c.listView
}

func (c *Container) SetOnTrigger(onTrigger func(button string, object messages.ItemObject, index int)) {
	c.onTrigger = onTrigger
}

func (c *Container) AddItem(obj messages.ItemObject) {
	c.items = append(c.items, obj)
	name := obj.GetName()
	if r, ok := FaceToRuneMap[uint16(obj.Face)]; ok {
		name = string(r.R) + " " + name
	}
	c.listView.AddItem(name, "", 0, nil)
}

func (c *Container) UpdateItem(obj messages.ItemObject) {
	for i, item := range c.items {
		if item.Tag == obj.Tag {
			c.items[i] = obj
			c.listView.SetItemText(i, obj.GetName(), "")
			break
		}
	}
}
