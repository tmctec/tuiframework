package tabapp

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/internal/keyname"
)

// namedKeys maps the names of the keys that are not a character to their
// key codes. The names are the ones tea.KeyPressMsg.String returns.
var namedKeys = func() map[string]rune {
	keys := map[string]rune{
		"enter": tea.KeyEnter, "tab": tea.KeyTab, "esc": tea.KeyEscape,
		"backspace": tea.KeyBackspace, "space": tea.KeySpace,
		"up": tea.KeyUp, "down": tea.KeyDown, "left": tea.KeyLeft, "right": tea.KeyRight,
		"home": tea.KeyHome, "end": tea.KeyEnd, "pgup": tea.KeyPgUp, "pgdown": tea.KeyPgDown,
		"insert": tea.KeyInsert, "delete": tea.KeyDelete,
	}
	for i, code := range []rune{
		tea.KeyF1, tea.KeyF2, tea.KeyF3, tea.KeyF4, tea.KeyF5, tea.KeyF6,
		tea.KeyF7, tea.KeyF8, tea.KeyF9, tea.KeyF10, tea.KeyF11, tea.KeyF12,
	} {
		keys["f"+strconv.Itoa(i+1)] = code
	}
	return keys
}()

// keyMods maps the modifier names in a key name to their flags.
var keyMods = map[string]tea.KeyMod{
	"ctrl": tea.ModCtrl, "alt": tea.ModAlt, "shift": tea.ModShift,
	"meta": tea.ModMeta, "hyper": tea.ModHyper, "super": tea.ModSuper,
}

// KeyPress returns the message a terminal sends when the named key is
// pressed. The name is what tea.KeyPressMsg.String returns and what BindKey
// takes: "s", "?", "space", "enter", "shift+tab", "ctrl+c". Use it to drive
// an application from a test or from code:
//
//	app.Update(tabapp.KeyPress("s"))
func KeyPress(name string) tea.KeyPressMsg {
	name = keyname.Normalize(name)

	// The last part is the key and the parts before it are modifiers. A name
	// that ends in "+" names the plus key, as in "+" and "ctrl++".
	parts := strings.Split(name, "+")
	keyName, mods := parts[len(parts)-1], parts[:len(parts)-1]
	if keyName == "" && len(mods) > 0 {
		keyName, mods = "+", mods[:len(mods)-1]
	}

	var k tea.Key
	for _, m := range mods {
		k.Mod |= keyMods[m]
	}

	code, named := namedKeys[keyName]
	r, size := utf8.DecodeRuneInString(keyName)
	switch {
	case named:
		k.Code = code
		if code == tea.KeySpace && k.Mod == 0 {
			k.Text = " "
		}
	case size != len(keyName):
		// More than one character that is no key name: text typed at once.
		k.Code, k.Text = tea.KeyExtended, keyName
	case k.Mod != 0:
		k.Code = r
	case unicode.IsUpper(r):
		// A capital arrives as its key with shift held and the capital as text.
		k.Code, k.ShiftedCode, k.Mod, k.Text = unicode.ToLower(r), r, tea.ModShift, keyName
	default:
		k.Code, k.Text = r, keyName
	}
	return tea.KeyPressMsg(k)
}
