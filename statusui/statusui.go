package statusui

import (
	"github.com/gen0cide/laforge"
	"github.com/rivo/tview"
)

var finderPage = "*finder*"

type statusapp struct {
	app         *tview.Application
	pages       *tview.Pages
	finderFocus tview.Primitive
	base        *laforge.Laforge
}

func newstatusapp(base *laforge.Laforge) *statusapp {
	sa := &statusapp{
		app:  tview.NewApplication(),
		base: base,
	}
	sa.finder(base)
	return sa
}

func (sa *statusapp) finder(base *laforge.Laforge) {
	objecttypes := tview.NewList().ShowSecondaryText(false)
	objecttypes.SetBorder(true).SetTitle("Object Types")
	objtree := tview.NewTreeView()
	objtree.SetBorder(true).SetTitle("Details")
	objlist := tview.NewList()
	objlist.ShowSecondaryText(false).SetDoneFunc(func() {
		objlist.Clear()
		objtree.SetRoot(nil)
		sa.app.SetFocus(objecttypes)
	})
	objlist.SetBorder(true).SetTitle("Library")

	flex := tview.NewFlex().
		AddItem(objecttypes, 0, 1, true).
		AddItem(objlist, 0, 1, false).
		AddItem(objtree, 0, 3, false)

	objecttypes.AddItem("user", "", 0, func() {
		objlist.Clear()
		objtree.SetRoot(nil)
		objlist.AddItem("local-user", "", 0, nil)
		sa.app.SetFocus(objlist)
		objlist.SetChangedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.User)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
		objlist.SetCurrentItem(0)
		objlist.SetSelectedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.User)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
	})
	objecttypes.AddItem("competition", "", 0, func() {
		objlist.Clear()
		objtree.SetRoot(nil)
		objlist.AddItem("competition", "", 0, nil)
		sa.app.SetFocus(objlist)
		objlist.SetChangedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.Competition)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
		objlist.SetCurrentItem(0)
		objlist.SetSelectedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.Competition)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
	})
	objecttypes.AddItem("environment", "", 0, func() {
		objlist.Clear()
		objtree.SetRoot(nil)
		objlist.AddItem("current", "", 0, nil)
		sa.app.SetFocus(objlist)
		objlist.SetChangedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.Environment)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
		objlist.SetCurrentItem(0)
		objlist.SetSelectedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.Environment)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
	})
	objecttypes.AddItem("hosts", "", 0, func() {
		objlist.Clear()
		objtree.SetRoot(nil)
		for name := range base.Hosts {
			objlist.AddItem(name, "", 0, nil)
		}
		sa.app.SetFocus(objlist)
		objlist.SetChangedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.Hosts[name])
			node.SetExpanded(true)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
		objlist.SetCurrentItem(0)
		objlist.SetSelectedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.Hosts[name])
			node.SetExpanded(true)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
	})
	objecttypes.AddItem("networks", "", 0, func() {
		objlist.Clear()
		objtree.SetRoot(nil)
		for name := range base.Networks {
			objlist.AddItem(name, "", 0, nil)
		}
		sa.app.SetFocus(objlist)
		objlist.SetChangedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.Networks[name])
			node.SetExpanded(true)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
		objlist.SetCurrentItem(0)
		objlist.SetSelectedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.Networks[name])
			node.SetExpanded(true)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
	})
	objecttypes.AddItem("identities", "", 0, func() {
		objlist.Clear()
		objtree.SetRoot(nil)
		for name := range base.Identities {
			objlist.AddItem(name, "", 0, nil)
		}
		sa.app.SetFocus(objlist)
		objlist.SetChangedFunc(func(i int, name string, t string, s rune) {
			node := identityTreeNode(name, base.Identities[name])
			node.SetReference(base.Identities[name])
			node.SetExpanded(true)
			node.SetSelectable(true)
			node.ExpandAll()
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
		objlist.SetCurrentItem(0)
		objlist.SetSelectedFunc(func(i int, name string, t string, s rune) {
			node := identityTreeNode(name, base.Identities[name])
			node.SetReference(base.Identities[name])
			node.SetExpanded(true)
			node.SetSelectable(true)
			node.ExpandAll()
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
	})
	objecttypes.AddItem("scripts", "", 0, func() {
		objlist.Clear()
		objtree.SetRoot(nil)
		for name := range base.Scripts {
			objlist.AddItem(name, "", 0, nil)
		}
		sa.app.SetFocus(objlist)
		objlist.SetChangedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.Scripts[name])
			node.SetExpanded(true)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
		objlist.SetCurrentItem(0)
		objlist.SetSelectedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.Scripts[name])
			node.SetExpanded(true)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
	})
	objecttypes.AddItem("commands", "", 0, func() {
		objlist.Clear()
		objtree.SetRoot(nil)
		for name := range base.Commands {
			objlist.AddItem(name, "", 0, nil)
		}
		sa.app.SetFocus(objlist)
		objlist.SetChangedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.Commands[name])
			node.SetExpanded(true)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
		objlist.SetCurrentItem(0)
		objlist.SetSelectedFunc(func(i int, name string, t string, s rune) {
			node := tview.NewTreeNode(name)
			node.SetReference(base.Commands[name])
			node.SetExpanded(true)
			objtree.SetRoot(node)
			objtree.SetCurrentNode(node)
		})
	})
	sa.pages = tview.NewPages().AddPage(finderPage, flex, true, true)
	sa.app.SetRoot(sa.pages, true)
}

// RenderLaforgeStatusUI renders an interactive console for exploring your data
func RenderLaforgeStatusUI(base *laforge.Laforge) error {
	sa := newstatusapp(base)
	return sa.app.Run()
}

func identityTreeNode(name string, i *laforge.Identity) *tview.TreeNode {
	newnode := tview.NewTreeNode(name)
	idnode := tview.NewTreeNode("ID")
	idnode.SetIndent(1)
	idnodeval := tview.NewTreeNode(i.ID)
	idnodeval.SetIndent(2)
	idnode.AddChild(idnodeval)
	newnode.AddChild(idnode)
	fnnode := tview.NewTreeNode("Firstname")
	fnnode.SetIndent(1)
	fnnodeval := tview.NewTreeNode(i.Firstname)
	fnnodeval.SetIndent(2)
	fnnode.AddChild(fnnodeval)
	newnode.AddChild(fnnode)
	lnnode := tview.NewTreeNode("Lastname")
	lnnode.SetIndent(1)
	lnnodeval := tview.NewTreeNode(i.Lastname)
	lnnodeval.SetIndent(2)
	lnnode.AddChild(lnnodeval)
	newnode.AddChild(lnnode)
	emnode := tview.NewTreeNode("Email")
	emnode.SetIndent(1)
	emnodeval := tview.NewTreeNode(i.Email)
	emnodeval.SetIndent(2)
	emnode.AddChild(emnodeval)
	newnode.AddChild(emnode)
	pwnode := tview.NewTreeNode("Password")
	pwnode.SetIndent(1)
	pwnodeval := tview.NewTreeNode(i.Password)
	pwnodeval.SetIndent(2)
	pwnode.AddChild(pwnodeval)
	newnode.AddChild(pwnode)
	dsnode := tview.NewTreeNode("Description")
	dsnode.SetIndent(1)
	dsnodeval := tview.NewTreeNode(i.Description)
	dsnodeval.SetIndent(2)
	dsnode.AddChild(dsnodeval)
	newnode.AddChild(dsnode)
	afnode := tview.NewTreeNode("Avatar File")
	afnode.SetIndent(1)
	afnodeval := tview.NewTreeNode(i.AvatarFile)
	afnodeval.SetIndent(2)
	afnode.AddChild(afnodeval)
	newnode.AddChild(afnode)
	return newnode
}
