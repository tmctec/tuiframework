package tabapp

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
)

type firedMsg struct{ name string }

// recorder is an application message handler that keeps what it received.
type recorder struct {
	fired []string  // the names of the firedMsg messages
	other []tea.Msg // every other message
}

func (r *recorder) handle(msg tea.Msg) (bool, tea.Cmd) {
	if f, ok := msg.(firedMsg); ok {
		r.fired = append(r.fired, f.name)
		return true, nil
	}
	r.other = append(r.other, msg)
	return false, nil
}

// keyApp is a sized one-page app whose page binds the given keys, each to a
// firedMsg that carries the key as written.
func keyApp(keys ...string) (*TabApp, *recorder) {
	p := page.New("Keys")
	for _, k := range keys {
		p.BindKey(k, "do "+k, firedMsg{name: k})
	}
	rec := &recorder{}
	app := sizedApp(New().MustAddPage(p).SetNavBarAutoHide(false))
	app.SetCustomMessageHandler(rec.handle)
	return app, rec
}

// key-1: the space bar dispatches the binding registered as "space".
func TestSpaceBarDispatchesSpaceBinding(t *testing.T) {
	app, rec := keyApp("space")
	app.Update(KeyPress("space"))
	if len(rec.fired) != 1 || rec.fired[0] != "space" {
		t.Errorf("fired = %v, want the space binding once", rec.fired)
	}
}

// key-2: a binding written " " is stored as space on a page, on the app and
// on an overlay, fires on the space bar, and shows as "space" in the footer.
func TestSpaceWrittenAsBlankIsSpace(t *testing.T) {
	app, rec := keyApp(" ")
	p := app.currentPage()
	if !p.HasKey("space") || !p.HasKey(" ") || p.Bindings()["space"].Key() != "space" {
		t.Fatalf("page bindings = %v, want the key stored as space", p.Bindings())
	}
	app.Update(KeyPress("space"))
	if len(rec.fired) != 1 {
		t.Errorf("page binding fired %d times, want 1", len(rec.fired))
	}
	if footer := ansiSequence.ReplaceAllString(app.Render(), ""); !strings.Contains(footer, "space do") {
		t.Errorf("the footer should name the key space:\n%s", footer)
	}

	p.UnbindKey(" ")
	if p.HasKey("space") {
		t.Error(`UnbindKey(" ") should remove the space binding`)
	}

	// The app: a universal binding.
	app.MustBindKey(" ", "Everywhere", firedMsg{name: "app"})
	if !app.IsUniversalKey("space") || !app.IsUniversalKey(" ") {
		t.Fatal("the universal binding should be stored as space")
	}
	app.Update(KeyPress("space"))
	if got := rec.fired[len(rec.fired)-1]; got != "app" {
		t.Errorf("last fired = %q, want the universal binding", got)
	}
	app.UnbindKey(" ")
	if app.IsUniversalKey("space") {
		t.Error(`UnbindKey(" ") should remove the universal binding`)
	}

	// An overlay: a registered key, a map of keys, and the exit key.
	o := overlays.New().SetExitKey(" ").
		RegisterKey(" ", firedMsg{name: "overlay"}).
		RegisterKeys(map[string]tea.Msg{" ": firedMsg{name: "overlay map"}})
	if o.ExitKey() != "space" || len(o.Keys()) != 1 || o.Keys()["space"] == nil {
		t.Errorf("overlay exit %q keys %v, want both stored as space", o.ExitKey(), o.Keys())
	}
	if got := (&overlays.NotificationSettings{}).SetExitKey(" ").ExitKey(); got != "space" {
		t.Errorf("notification settings exit key = %q, want space", got)
	}

	// The tab-navigation keys.
	app.SetTabNavKeyForward(" ").SetTabNavKeyReverse(" ")
	if app.tabNavKeyFwd != "space" || app.tabNavKeyRev != "space" {
		t.Errorf("tab-nav keys = %q %q, want space", app.tabNavKeyFwd, app.tabNavKeyRev)
	}
}

// key-3: a key release dispatches nothing. The terminal reports releases
// only on request, and a release is not a press.
func TestKeyReleaseDispatchesNothing(t *testing.T) {
	app, rec := keyApp("s")
	_, cmd := app.Update(tea.KeyReleaseMsg(tea.Key{Code: 's', Text: "s"}))
	if cmd != nil || len(rec.fired) != 0 {
		t.Errorf("a release fired %v (cmd nil = %v), want nothing", rec.fired, cmd == nil)
	}
	if len(rec.other) != 1 {
		t.Errorf("the handler received %d other messages, want the release itself", len(rec.other))
	}

	// The quit key released does not quit either.
	if _, cmd := app.Update(tea.KeyReleaseMsg(tea.Key{Code: 'q', Text: "q"})); cmd != nil {
		t.Error("releasing q should return no command")
	}
}

// key-4: on a terminal that tells the chords apart, shift+enter dispatches
// its own binding and leaves the enter binding alone.
func TestShiftEnterDispatchesItsOwnBinding(t *testing.T) {
	app, rec := keyApp("enter", "shift+enter", "ctrl+h", "backspace")
	app.Update(KeyPress("shift+enter"))
	app.Update(KeyPress("ctrl+h"))
	if got := strings.Join(rec.fired, ","); got != "shift+enter,ctrl+h" {
		t.Errorf("fired = %q, want shift+enter then ctrl+h", got)
	}
}

// key-5: a terminal that sends plain enter fires the enter binding, never
// the shift+enter one.
func TestPlainEnterDoesNotFireShiftEnter(t *testing.T) {
	app, rec := keyApp("enter", "shift+enter", "ctrl+h", "backspace")
	app.Update(KeyPress("enter"))
	app.Update(KeyPress("backspace"))
	if got := strings.Join(rec.fired, ","); got != "enter,backspace" {
		t.Errorf("fired = %q, want the fallback bindings enter then backspace", got)
	}
}

// key-6: pasted text is not keys. Pasting "q" does not quit; the paste
// reaches the application's handler as one message.
func TestPasteIsNotKeys(t *testing.T) {
	app, rec := keyApp("s")
	_, cmd := app.Update(tea.PasteMsg{Content: "qs"})
	if cmd != nil {
		t.Error("a paste should return no command; it must not quit")
	}
	if len(rec.fired) != 0 {
		t.Errorf("a paste fired %v, want no binding", rec.fired)
	}
	if len(rec.other) != 1 {
		t.Fatalf("the handler received %d messages, want the paste", len(rec.other))
	}
	if p, ok := rec.other[0].(tea.PasteMsg); !ok || p.Content != "qs" {
		t.Errorf("the handler received %#v, want the paste", rec.other[0])
	}
}

// key-7: Model.KeyDisambiguation follows the terminal's report.
func TestKeyDisambiguationFollowsTheTerminal(t *testing.T) {
	app, rec := keyApp()
	if app.Model().KeyDisambiguation {
		t.Fatal("a new app has no report: KeyDisambiguation should be false")
	}
	app.Update(tea.KeyboardEnhancementsMsg{Flags: ansi.KittyDisambiguateEscapeCodes})
	if !app.Model().KeyDisambiguation {
		t.Error("KeyDisambiguation should be true after the terminal reports it")
	}
	if len(rec.other) != 1 {
		t.Errorf("the handler received %d messages, want the report", len(rec.other))
	}
	app.Update(tea.KeyboardEnhancementsMsg{})
	if app.Model().KeyDisambiguation {
		t.Error("a report without the feature turns KeyDisambiguation off")
	}
}

// key-8: a page key that the app binds globally is a conflict however the
// space bar is written.
func TestSpaceConflictIsReported(t *testing.T) {
	app := New().MustBindKey("space", "Everywhere", NoOpMsg{})
	err := app.AddPage(page.New("Clash").BindKey(" ", "Here", NoOpMsg{}))
	if err == nil {
		t.Fatal(`a page binding " " should conflict with the universal space binding`)
	}
	if !strings.Contains(err.Error(), "space") {
		t.Errorf("error %q should name the key", err)
	}

	// Written the other way round it is the same conflict.
	app = New().MustBindKey(" ", "Everywhere", NoOpMsg{})
	if app.AddPage(page.New("Clash").BindKey("space", "Here", NoOpMsg{})) == nil {
		t.Error(`a page binding space should conflict with a universal " " binding`)
	}

	// A hardcoded key stays refused.
	if New().BindKey("ctrl+c", "Copy", NoOpMsg{}) == nil {
		t.Error("ctrl+c is hardcoded and cannot be bound")
	}
}
