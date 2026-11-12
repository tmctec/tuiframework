#!/usr/bin/env bash
set -euo pipefail

# ============================================================================
# makehelp.sh — Complex logic extracted from Makefile
# ============================================================================
# The Makefile defines WHAT; this script implements HOW for anything too
# complex for inline make recipes (>10 lines, OS branching, multi-step ops).
# ============================================================================

BINARY_DIR="./bin"
INSTALL_DIR="${HOME}/.local/bin"
STICKERS_DIR="../stickers"

# ---------------------------------------------------------------------------
# Helper functions
# ---------------------------------------------------------------------------

check_command() {
    if ! command -v "$1" &>/dev/null; then
        echo "Error: $1 is required but not installed" >&2
        return 1
    fi
}

# ---------------------------------------------------------------------------
# Command functions
# ---------------------------------------------------------------------------

cmd_prereqs() {
    echo "Checking prerequisites..."

    if ! command -v go &>/dev/null; then
        case "$(uname -s)" in
            Darwin)
                check_command brew || { echo "Install Homebrew first: https://brew.sh" >&2; exit 1; }
                brew install go
                ;;
            Linux)
                if command -v apt-get &>/dev/null; then
                    sudo apt-get update && sudo apt-get install -y golang
                elif command -v dnf &>/dev/null; then
                    sudo dnf install -y golang
                else
                    echo "Error: install Go from https://go.dev/dl/" >&2
                    exit 1
                fi
                ;;
            *)
                echo "Unsupported OS: $(uname -s)" >&2
                exit 1
                ;;
        esac
    fi
    echo "go: $(go version)"

    # A local replace for stickers in go.mod needs the sibling checkout to
    # exist, on the branch that holds the module go.mod names.
    if grep -qE "replace github.com/tmctec/stickers(/v2)? => ${STICKERS_DIR}" go.mod; then
        if [[ ! -d "${STICKERS_DIR}/flexbox" ]]; then
            echo "Error: go.mod points at ${STICKERS_DIR} but it is missing." >&2
            echo "Fix:   git clone git@github.com:tmctec/stickers.git ${STICKERS_DIR}" >&2
            exit 1
        fi
        if ! grep -q "^module github.com/tmctec/stickers/v2$" "${STICKERS_DIR}/go.mod"; then
            echo "Error: ${STICKERS_DIR} is not the stickers/v2 module go.mod needs." >&2
            echo "Fix:   git -C ${STICKERS_DIR} checkout main" >&2
            exit 1
        fi
        echo "stickers: ${STICKERS_DIR} (local replace)"
    fi

    if command -v golangci-lint &>/dev/null; then
        echo "golangci-lint: $(golangci-lint version 2>/dev/null | head -1)"
    else
        echo "golangci-lint: not installed (make lint falls back to go vet)"
    fi

    echo "All prerequisites satisfied"
}

cmd_run() {
    # `run` = launch the demo. First arg names the binary under ./bin; any
    # remaining args are passed straight through (make run ARGS="...").
    local name="${1:-example-demo}"
    shift || true
    local binary="${BINARY_DIR}/${name}"
    if [[ ! -x "$binary" ]]; then
        echo "Error: $binary not found (run 'make build' first)" >&2
        exit 1
    fi
    if [[ ! -t 1 ]]; then
        # No terminal is attached, as when make runs from an editor or an
        # assistant's shell. Open the demo in a window of the terminal app
        # this session belongs to, when there is one.
        if open_in_terminal_window "$(pwd)/${binary#./}" "$@"; then
            echo "No terminal is attached here, so ${name} is running in a new terminal window."
            return 0
        fi
        echo "Error: the demo draws a full-screen interface and needs a terminal." >&2
        echo "Fix:   run 'make demo' in a terminal, or 'make smoke' to run it in a pseudo-terminal." >&2
        exit 1
    fi
    exec "$binary" "$@"
}

open_in_terminal_window() {
    # Runs a command in a new window of iTerm2 or Terminal on macOS, whichever
    # the session was started from. Returns 1 when it cannot: another OS, no
    # terminal app in the environment, or NO_WINDOW=1 set by the caller.
    [[ "${NO_WINDOW:-}" != "1" ]] || return 1
    [[ "$(uname -s)" == "Darwin" ]] || return 1
    command -v osascript &>/dev/null || return 1

    # The command line, quoted for a shell and then for an AppleScript string.
    local line
    line="$(printf '%q ' "$@")"
    line="${line% }"
    line="${line//\\/\\\\}"
    line="${line//\"/\\\"}"

    case "${TERM_PROGRAM:-}" in
        iTerm.app)
            osascript -e "tell application \"iTerm2\" to create window with default profile command \"${line}\"" >/dev/null
            ;;
        Apple_Terminal)
            osascript -e "tell application \"Terminal\" to do script \"${line}\"" \
                      -e "tell application \"Terminal\" to activate" >/dev/null
            ;;
        *)
            return 1
            ;;
    esac
}

cmd_smoke() {
    # Runs the demo in a pseudo-terminal, quits it with q, and checks that
    # it drew its first page, exited 0 and left no process behind.
    local name="${1:-example-demo}"
    local binary="${BINARY_DIR}/${name}"
    if [[ ! -x "$binary" ]]; then
        echo "Error: $binary not found (run 'make build' first)" >&2
        exit 1
    fi
    check_command script || { echo "Fix:   install util-linux, which provides script" >&2; exit 1; }

    # Demo processes that are already running, such as one in another
    # window, are not this run's to judge.
    local before
    before="$({ pgrep -f "${BINARY_DIR#./}/${name}" || true; } | sort | tr '\n' ' ')"

    local out status=0
    out="$(mktemp)"
    # script's arguments differ between the BSD and the util-linux versions.
    case "$(uname -s)" in
        Darwin|*BSD) (sleep 1.5; printf 'q') | script -q /dev/null "$binary" >"$out" 2>&1 & ;;
        *)           (sleep 1.5; printf 'q') | script -qec "$binary" /dev/null >"$out" 2>&1 & ;;
    esac
    local runner=$!
    # A demo that does not quit is stopped after 15 seconds and fails the check.
    ( sleep 15; kill "$runner" 2>/dev/null ) &
    local watchdog=$!
    wait "$runner" || status=$?
    kill "$watchdog" 2>/dev/null || true

    if [[ $status -ne 0 ]]; then
        echo "smoke: FAIL: the demo exited with status $status" >&2
        tail -c 400 "$out" >&2
        rm -f "$out"
        exit 1
    fi
    if ! LC_ALL=C grep -aq "Welcome" "$out"; then
        echo "smoke: FAIL: the demo never drew its Welcome page" >&2
        rm -f "$out"
        exit 1
    fi
    rm -f "$out"
    local after
    after="$({ pgrep -f "${BINARY_DIR#./}/${name}" || true; } | sort | tr '\n' ' ')"
    if [[ "$after" != "$before" ]]; then
        echo "smoke: FAIL: a ${name} process is still running (pids before: ${before:-none}; after: ${after:-none})" >&2
        exit 1
    fi
    echo "smoke: ok: ${name} drew its first page, quit on q and left no process"
}

cmd_lint() {
    if command -v golangci-lint &>/dev/null; then
        golangci-lint run ./...
    else
        echo "golangci-lint not installed; running go vet instead" >&2
        echo "Install: brew install golangci-lint" >&2
        go vet ./...
    fi
}

cmd_install_dev() {
    # install-dev creates SYMLINKS so rebuilds work without reinstalling
    mkdir -p "$INSTALL_DIR"
    local name
    for name in "$@"; do
        local binary="${BINARY_DIR}/${name}"
        if [[ ! -x "$binary" ]]; then
            echo "Error: $binary not found (run 'make build' first)" >&2
            exit 1
        fi
        ln -sf "$(pwd)/${binary}" "${INSTALL_DIR}/${name}"
        echo "Installed: ${INSTALL_DIR}/${name} -> $(pwd)/${binary} (symlink)"
    done
}

cmd_uninstall() {
    local name
    for name in "$@"; do
        rm -f "${INSTALL_DIR}/${name}"
        echo "Uninstalled: ${INSTALL_DIR}/${name}"
    done
}

cmd_release_verify() {
    # Proves a version can be used the way users use it. From a scratch
    # module outside the repository it requires the tagged module, builds a
    # program that renders a frame with every package, and runs make check in
    # a clean clone of the tag. SOURCE=local uses the local tag (before it is
    # pushed); otherwise the Go toolchain fetches it from origin.
    local version="${1:-}" source="${2:-}"
    local module="github.com/tmctec/tuiframework"
    if [[ -z "$version" ]]; then
        echo "Error: no version given." >&2
        echo "Fix:   make release-verify VERSION=0.1.0 [SOURCE=local]" >&2
        exit 1
    fi
    local tag="v${version#v}"
    local origin
    origin="$(git remote get-url origin)"

    if [[ "$source" == "local" ]]; then
        if ! git rev-parse -q --verify "refs/tags/${tag}" >/dev/null; then
            echo "Error: tag ${tag} does not exist in this repository." >&2
            echo "Fix:   git tag -a ${tag} -m \"tuiframework ${tag}\"" >&2
            exit 1
        fi
    elif ! git ls-remote --exit-code origin "refs/tags/${tag}" >/dev/null 2>&1; then
        echo "Error: tag ${tag} is not on origin (${origin})." >&2
        echo "Fix:   git push origin ${tag}, or verify the local tag with SOURCE=local" >&2
        exit 1
    fi

    local work
    work="$(mktemp -d "${TMPDIR:-/tmp}/tuiframework-release-verify.XXXXXX")"
    trap 'rm -rf "$work"' RETURN
    export GOWORK=off GOFLAGS=-mod=mod GOPRIVATE='' GONOSUMDB='' GONOPROXY=''

    # A clean clone of the tag.
    if [[ "$source" == "local" ]]; then
        git clone -q --branch "$tag" "$(pwd)" "$work/clone" 2>/dev/null
    else
        git clone -q --depth 1 --branch "$tag" "$origin" "$work/clone" 2>/dev/null
    fi

    # A program that uses the module the way an application does.
    mkdir -p "$work/app"
    cat > "$work/app/main.go" <<'GO'
package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
	"github.com/tmctec/tuiframework/tabapp"
)

func main() {
	app := tabapp.New().UseTheme(styles.Nord).MustAddPage(page.New("Home"))
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app.ShowNotification(overlays.Info("release check"))
	if lines := strings.Count(app.Render(), "\n") + 1; lines != 24 {
		fmt.Println("frame has", lines, "lines, want 24")
		os.Exit(1)
	}
	fmt.Println("rendered an 80x24 frame")
}
GO
    (
        cd "$work/app"
        go mod init release-verify >/dev/null 2>&1
        if [[ "$source" == "local" ]]; then
            go mod edit -require="${module}@${tag}" -replace="${module}=$work/clone"
        else
            go get "${module}@${tag}" >/dev/null
        fi
        go mod tidy >/dev/null
        local resolved
        resolved="$(go list -m -f '{{.Version}}' "$module")"
        if [[ "$resolved" != "$tag" ]]; then
            echo "Error: the module resolved to ${resolved}, not ${tag}." >&2
            exit 1
        fi
        go run .
    )
    echo "release-verify: ${module} ${tag} resolves, builds and renders"

    make -C "$work/clone" check >/dev/null
    echo "release-verify: make check passes in a clean clone of ${tag}"
}

SECURITY_TOOLS_DIR="${BINARY_DIR}/tools"
GITLEAKS_VERSION="v8.30.1"
GOVULNCHECK_VERSION="v1.9.0"
GOSEC_VERSION="v2.30.0"

install_tool() {
    # Installs a Go tool at a pinned version into ./bin/tools, once.
    local name="$1" pkg="$2"
    if [[ ! -x "${SECURITY_TOOLS_DIR}/${name}" ]]; then
        echo "security: installing ${pkg}"
        GOBIN="$(pwd)/${SECURITY_TOOLS_DIR#./}" go install "$pkg"
    fi
}

cmd_security() {
    # Scans for secrets (tree and every commit), known vulnerabilities, Go
    # security defects, and names that must not appear. The names are a
    # regular expression in BANNED_NAMES or in the file BANNED_NAMES_FILE
    # (default ~/.config/tuiframework/banned-names), kept out of the
    # repository so the check does not publish the names it guards.
    mkdir -p "$SECURITY_TOOLS_DIR"
    install_tool gitleaks "github.com/zricethezav/gitleaks/v8@${GITLEAKS_VERSION}"
    install_tool govulncheck "golang.org/x/vuln/cmd/govulncheck@${GOVULNCHECK_VERSION}"
    install_tool gosec "github.com/securego/gosec/v2/cmd/gosec@${GOSEC_VERSION}"
    local failed=0

    echo "security: gitleaks (working tree)"
    "${SECURITY_TOOLS_DIR}/gitleaks" dir . --no-banner --redact --exit-code 1 || failed=1
    echo "security: gitleaks (every commit)"
    "${SECURITY_TOOLS_DIR}/gitleaks" git . --no-banner --redact --exit-code 1 --log-opts=--all || failed=1
    echo "security: govulncheck"
    "${SECURITY_TOOLS_DIR}/govulncheck" ./... || failed=1
    echo "security: gosec"
    "${SECURITY_TOOLS_DIR}/gosec" -quiet ./... || failed=1

    local names="${BANNED_NAMES:-}"
    local names_file="${BANNED_NAMES_FILE:-${HOME}/.config/tuiframework/banned-names}"
    if [[ -z "$names" && -r "$names_file" ]]; then
        names="$(grep -v '^[[:space:]]*\(#\|$\)' "$names_file" | paste -sd'|' -)"
    fi
    if [[ -z "$names" ]]; then
        echo "security: name check skipped: set BANNED_NAMES or create ${names_file}"
    else
        echo "security: banned names"
        if git grep -nIiE "$names" -- . >/dev/null; then
            echo "security: FAIL: a banned name appears in a tracked file:" >&2
            git grep -nIiE "$names" -- . | cut -d: -f1,2 >&2
            failed=1
        fi
        if git log --all --format='%an %ae %cn %ce%n%B' | grep -qiE "$names"; then
            echo "security: FAIL: a banned name appears in commit metadata" >&2
            failed=1
        fi
    fi

    if [[ $failed -ne 0 ]]; then
        echo "security: FAIL" >&2
        exit 1
    fi
    echo "security: ok"
}

# ---------------------------------------------------------------------------
# Dispatcher — one entry per makehelp.sh-delegated target
# ---------------------------------------------------------------------------

case "${1:-}" in
    prereqs)      cmd_prereqs ;;
    run)          shift; cmd_run "$@" ;;
    smoke)        shift; cmd_smoke "$@" ;;
    release-verify) shift; cmd_release_verify "$@" ;;
    lint)         cmd_lint ;;
    security)     cmd_security ;;
    install-dev)  shift; cmd_install_dev "$@" ;;
    uninstall)    shift; cmd_uninstall "$@" ;;
    *)
        echo "Usage: $0 {prereqs|run <name> [args]|smoke <name>|release-verify <version> [local]|lint|security|install-dev <names>|uninstall <names>}" >&2
        exit 1
        ;;
esac
