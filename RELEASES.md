# Releases

## v0.1.0 — 2026-11-11 — Built on the new Charm libraries

**New**
- Your app can set its terminal window's title, show progress in the terminal's own progress bar, and give the terminal's background the theme's colors, so the edge around the screen matches.
- Themes can follow the terminal: light on a light background, dark on a dark one.
- Make a whole theme from three colors: a background, a text color and an accent. Colors that would be hard to read together are refused.
- Borders can fade from one color to another.
- A new built-in theme, Charm, in Charm's own colors.
- Apps can tell newer key combinations apart, such as Shift+Enter, on terminals that support them. Pasted text never triggers a shortcut.
- The demo has a Terminal page that shows what your terminal supports.

**Improved**
- Screens draw faster and use less memory.
- The same app runs in a terminal, over SSH, or in an automated test.
- The demo opens in a new terminal window when started from an editor.

**Fixed**
- Italic, underlined and other styled text, and links, keep their look on every page.

**Security**
- Report a vulnerability privately through GitHub security advisories, as SECURITY.md describes.
