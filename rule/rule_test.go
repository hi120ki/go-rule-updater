package rule

import (
	"reflect"
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

func TestGetAddedRules(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		old     string
		new     string
		want    []string
		wantErr bool
	}{
		{
			name: "simple addition",
			old: `config:
  - name: dc898828-aedf-4928-8778-e1015014967c
`,
			new: `config:
  - name: dc898828-aedf-4928-8778-e1015014967c
  - name: 123e4567-e89b-12d3-a456-426614174000
`,
			want:    []string{"123e4567-e89b-12d3-a456-426614174000"},
			wantErr: false,
		},
		{
			name: "multiple additions",
			old: `config:
  - name: dc898828-aedf-4928-8778-e1015014967c
`,
			new: `config:
  - name: dc898828-aedf-4928-8778-e1015014967c
  - name: 0a8f85f4-821a-4346-9d03-62d8e056534a
  - name: ab3c43a8-fd03-410c-87e1-819bc236b3d4
`,
			want:    []string{"0a8f85f4-821a-4346-9d03-62d8e056534a", "ab3c43a8-fd03-410c-87e1-819bc236b3d4"},
			wantErr: false,
		},
		{
			name: "no additions",
			old: `config:
  - name: dc898828-aedf-4928-8778-e1015014967c
`,
			new: `config:
  - name: dc898828-aedf-4928-8778-e1015014967c
`,
			want:    []string{},
			wantErr: false,
		},
		{
			name: "with new rules",
			old: `config:
  - name: dc898828-aedf-4928-8778-e1015014967c
  - name: 0a8f85f4-821a-4346-9d03-62d8e056534a
`,
			new: `config:
  - name: dc898828-aedf-4928-8778-e1015014967c
  - name: ab3c43a8-fd03-410c-87e1-819bc236b3d4
`,
			want:    []string{"ab3c43a8-fd03-410c-87e1-819bc236b3d4"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := GetAddedRules(tt.old, tt.new)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetAddedRules() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetAddedRules() succeeded unexpectedly")
			}
			if reflect.DeepEqual(got, tt.want) == false {
				t.Errorf("GetAddedRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
