package rule

import (
	"strings"
	"testing"
)

func TestAddKeepsComments(t *testing.T) {
	t.Parallel()

	const input = `# top level
config:
  # existing comment
  - name: foo
`

	got, err := Add(input, "bar")
	if err != nil {
		t.Fatalf("Add returned error: %v", err)
	}

	if !strings.Contains(got, "# top level") {
		t.Fatalf("expected top comment to remain, got:\n%s", got)
	}
	if !strings.Contains(got, "# existing comment") {
		t.Fatalf("expected inline comment to remain, got:\n%s", got)
	}
	if !strings.Contains(got, "- name: bar") {
		t.Fatalf("expected new rule to be added, got:\n%s", got)
	}
}
