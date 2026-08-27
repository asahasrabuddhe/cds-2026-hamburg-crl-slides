// Command snippets pulls Go code out of the vendored cds-2026-hamburg-crl-code
// submodule by anchor comment, so the deck never contains pasted Go. It fails
// loudly if slides.md references an anchor that does not exist in the
// submodule, since a silently missing snippet is worse than a broken build.
//
// The source directory can be overridden with CRL_CODE_DIR, which is how you
// run this against a sibling checkout before the submodule exists.
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	defaultVendorDir = "vendor/crl-code"
	snippetsDir      = "snippets"
	slidesFile       = "slides.md"
)

var (
	startRe     = regexp.MustCompile(`//\s*snippet:start\s+(\S+)`)
	endRe       = regexp.MustCompile(`//\s*snippet:end\s+(\S+)`)
	referenceRe = regexp.MustCompile(`<<<\s*@/snippets/([\w-]+)\.go`)
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "snippets:", err)
		os.Exit(1)
	}
}

func vendorDir() string {
	if dir := os.Getenv("CRL_CODE_DIR"); dir != "" {
		return dir
	}
	return defaultVendorDir
}

func run() error {
	src := vendorDir()

	referenced, err := referencedAnchors(slidesFile)
	if err != nil {
		return fmt.Errorf("scan %s for snippet references: %w", slidesFile, err)
	}

	found, err := extractAnchors(src)
	if err != nil {
		return fmt.Errorf("extract anchors from %s: %w", src, err)
	}

	if err := os.MkdirAll(snippetsDir, 0o755); err != nil {
		return fmt.Errorf("create %s: %w", snippetsDir, err)
	}

	var missing []string
	for name := range referenced {
		path := filepath.Join(snippetsDir, name+".go")
		body, ok := found[name]
		if !ok {
			missing = append(missing, name)
			// Slidev's markdown-it importer throws if the file does not
			// exist at all, which kills the whole render rather than one
			// slide. A placeholder keeps the other slides previewable while
			// this command still fails loudly.
			body = fmt.Sprintf("// snippet %q not found in %s\n", name, src)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
		fmt.Printf("snippets: wrote %s\n", path)
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("slides.md references snippets with no matching anchor in %s: %s",
			src, strings.Join(missing, ", "))
	}
	return nil
}

// referencedAnchors returns every snippet name slides.md imports via Slidev's
// `<<< @/snippets/<name>.go` code-import syntax.
func referencedAnchors(path string) (map[string]struct{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := map[string]struct{}{}
	for _, m := range referenceRe.FindAllStringSubmatch(string(data), -1) {
		out[m[1]] = struct{}{}
	}
	return out, nil
}

// extractAnchors walks the source tree and collects every region delimited by
// `// snippet:start <name>` and `// snippet:end <name>`. The anchor comments
// themselves are not included in the extracted body, and the common leading
// indentation is stripped so the slide does not start every line with a tab.
func extractAnchors(root string) (map[string]string, error) {
	out := map[string]string{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		open := map[string][]string{}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()

			if m := startRe.FindStringSubmatch(line); m != nil {
				open[m[1]] = []string{}
				continue
			}
			if m := endRe.FindStringSubmatch(line); m != nil {
				if body, ok := open[m[1]]; ok {
					out[m[1]] = dedent(body)
					delete(open, m[1])
				}
				continue
			}
			for name := range open {
				open[name] = append(open[name], line)
			}
		}
		return scanner.Err()
	})
	return out, err
}

// dedent removes the shortest leading run of tabs common to every non-blank
// line, so a snippet lifted out of a function body sits flush on the slide.
func dedent(lines []string) string {
	depth := -1
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		n := len(l) - len(strings.TrimLeft(l, "\t"))
		if depth == -1 || n < depth {
			depth = n
		}
	}
	if depth < 1 {
		return strings.Join(lines, "\n") + "\n"
	}
	for i, l := range lines {
		if len(l) >= depth {
			lines[i] = l[depth:]
		}
	}
	return strings.Join(lines, "\n") + "\n"
}
