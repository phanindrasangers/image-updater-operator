/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gitwriteback

import (
	"strings"
	"testing"
)

func sampleData() CommitData {
	return CommitData{
		Kind:      "Deployment",
		Name:      "web",
		Namespace: "default",
		Changes: []Change{
			{File: "values.yaml", Container: "app", Repository: "docker.io/acme/web", Tag: "v2", Image: "docker.io/acme/web:v2", OldImage: "docker.io/acme/web:v1"},
			{File: "values.yaml", Container: "side", Repository: "docker.io/acme/side", Tag: "9", Image: "docker.io/acme/side:9", OldImage: "docker.io/acme/side:8"},
		},
	}.WithPrimary()
}

func TestRenderCommitMessage_DefaultListsEveryChange(t *testing.T) {
	msg, err := RenderCommitMessage("", sampleData())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(msg, "chore(images): update web") {
		t.Errorf("subject = %q, want it to start with the default subject", msg)
	}
	for _, want := range []string{"app: docker.io/acme/web:v1 -> docker.io/acme/web:v2 (values.yaml)", "side: docker.io/acme/side:8 -> docker.io/acme/side:9 (values.yaml)"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message missing %q\n--- got ---\n%s", want, msg)
		}
	}
	if strings.HasSuffix(msg, "\n") {
		t.Errorf("message should have trailing whitespace trimmed, got %q", msg)
	}
}

func TestRenderCommitMessage_CustomTemplateUsesPrimaryFields(t *testing.T) {
	msg, err := RenderCommitMessage("ci: bump {{.Container}} to {{.Tag}}", sampleData())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "ci: bump app to v2" {
		t.Errorf("got %q, want %q", msg, "ci: bump app to v2")
	}
}

func TestRenderCommitMessage_RangeAndFields(t *testing.T) {
	tmpl := "release {{.Name}}\n{{range .Changes}}- {{.Image}}\n{{end}}"
	msg, err := RenderCommitMessage(tmpl, sampleData())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "release web\n- docker.io/acme/web:v2\n- docker.io/acme/side:9"
	if msg != want {
		t.Errorf("got %q, want %q", msg, want)
	}
}

func TestRenderCommitMessage_InvalidTemplateErrors(t *testing.T) {
	if _, err := RenderCommitMessage("{{.Nope", sampleData()); err == nil {
		t.Fatal("expected a parse error for a malformed template")
	}
}

func TestRenderCommitMessage_EmptyOutputErrors(t *testing.T) {
	if _, err := RenderCommitMessage("{{if false}}x{{end}}", sampleData()); err == nil {
		t.Fatal("expected an error when the template renders to nothing")
	}
}
