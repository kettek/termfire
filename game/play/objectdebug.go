package play

import "github.com/rivo/tview"

type ObjectDebugView struct {
	container *tview.Flex
	list      *tview.List
}

func (v *ObjectDebugView) Init() {
	v.container = tview.NewFlex()
	v.container.SetDirection(tview.FlexRow)
	v.container.SetTitle("Object Debug")
	v.container.SetBorder(true)

	v.list = tview.NewList()
	v.container.AddItem(v.list, 0, 1, true)
}

func (v *ObjectDebugView) GetContainer() *tview.Flex {
	return v.container
}

func (v *ObjectDebugView) Refresh() {
	v.list.Clear()

	for _, r := range FaceToRuneMap {
		v.list.AddItem(string(r.R), r.F.String()+" "+r.B.String(), 0, nil)
	}
}
