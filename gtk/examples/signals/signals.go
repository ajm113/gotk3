package main

import (
	"fmt"
	"github.com/conformal/gotk3/glib"
	"github.com/conformal/gotk3/gtk"
	"log"
)

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.Add(windowWidget())
	win.ShowAll()

	gtk.Main()
}

func windowWidget() *gtk.Widget {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	entry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry:", err)
	}
	s, _ := entry.GetText()
	label, err := gtk.LabelNew(s)
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	grid.Add(entry)
	entry.SetHExpand(true)
	grid.AttachNextTo(label, entry, gtk.POS_RIGHT, 1, 1)
	label.SetHExpand(true)

	// Connects this entry's "activate" signal (which is emitted whenever
	// Enter is pressed when the Entry is activated) to an anonymous
	// function that gets the current text of the entry and sets the text of
	// the label beside it with it.  Unlike with native GTK callbacks,
	// (*glib.Object).Connect() supports closures.  In this example, this is
	// demonstrated by using the label variable.  Without closures, a
	// pointer to the label would need to be passed in as user data
	// (demonstrated in the next example).
	entry.Connect("activate", func() {
		s, _ := entry.GetText()
		label.SetText(s)
	})

	sb, err := gtk.SpinButtonNewWithRange(0, 1, 0.1)
	if err != nil {
		log.Fatal("Unable to create spin button:", err)
	}
	pb, err := gtk.ProgressBarNew()
	if err != nil {
		log.Fatal("Unable to create progress bar:", err)
	}
	grid.Add(sb)
	sb.SetHExpand(true)
	grid.AttachNextTo(pb, sb, gtk.POS_RIGHT, 1, 1)
	label.SetHExpand(true)

	// Pass in a ProgressBar and the target ScrollButton as user data rather
	// than using the sb and pb variables scoped to the anonymous func.
	// This can be useful when passing in a closure that has already been
	// generated, but when you still wish to connect the callback with some
	// variables only visible in this scope.
	//
	// (*glib.CallbackContext).Target() returns an interface{} which is
	// actually a *glib.Object, so to use as a gtk.SpinButton, a new
	// composite literal must be created.
	sb.ConnectWithData("value-changed", func(ctx *glib.CallbackContext) {
		obj, ok := ctx.Target().(*glib.Object)
		if !ok {
			fmt.Println("here")
			return
		}
		pb, ok := ctx.Data().(*gtk.ProgressBar)
		if !ok {
			return
		}

		sb := &gtk.SpinButton{gtk.Entry{gtk.Widget{
			glib.InitiallyUnowned{obj}}}}
		pb.SetFraction(sb.GetValue() / 1)
	}, pb)

	label, err = gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	s = "Hyperlink to <a href=\"https://www.cyphertite.com/\">Cyphertite</a> for your clicking pleasure"
	label.SetMarkup(s)
	grid.AttachNextTo(label, sb, gtk.POS_BOTTOM, 2, 1)

	// Some GTK callback functions require arguments, such as the
	// 'gchar *uri' argument of GtkLabel's "activate-link" signal.  To use
	// these arguments, pass in a *glib.CallbackContext as an argument, and
	// access by calling (*glib.CallbackContext).Arg(n) for the nth
	// argument.
	label.Connect("activate-link", func (ctx *glib.CallbackContext) {
		uri := ctx.Arg(0).String()
		fmt.Println("you clicked a link to:", uri)
	})

	return &grid.Container.Widget
}