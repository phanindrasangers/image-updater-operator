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
	"fmt"
	"strings"
	"text/template"
)

// DefaultCommitMessageTemplate renders the commit message when a workload does
// not set the git-commit-message annotation. It is a conventional-commits style
// subject naming the workload, followed by a body listing every container that
// changed and the path it was written to.
const DefaultCommitMessageTemplate = `chore(images): update {{.Name}}

{{range .Changes}}{{.Container}}: {{.OldImage}} -> {{.Image}} ({{.File}})
{{end}}`

// Change describes a single image edit applied during a write-back.
type Change struct {
	// File is the repo-relative path that was edited.
	File string
	// Container is the name of the container whose image changed.
	Container string
	// Repository is the image repository, without a tag.
	Repository string
	// Tag is the selected tag.
	Tag string
	// Image is the full reference, repository:tag.
	Image string
	// OldImage is the reference the container carried before the update.
	OldImage string
}

// CommitData is the template context passed to a commit-message template. The
// singular Container/Repository/Tag/Image/OldImage fields mirror the first
// change, a convenience for the common single-image case so templates can read
// {{.Tag}} without ranging.
type CommitData struct {
	// Kind is the workload kind, e.g. Deployment.
	Kind string
	// Name and Namespace identify the workload being updated.
	Name      string
	Namespace string
	// Changes is every image edit included in this commit.
	Changes []Change

	Container  string
	Repository string
	Tag        string
	Image      string
	OldImage   string
}

// WithPrimary returns d with the singular convenience fields populated from the
// first change, if any.
func (d CommitData) WithPrimary() CommitData {
	if len(d.Changes) > 0 {
		c := d.Changes[0]
		d.Container, d.Repository, d.Tag, d.Image, d.OldImage =
			c.Container, c.Repository, c.Tag, c.Image, c.OldImage
	}
	return d
}

// RenderCommitMessage renders tmpl with data, falling back to the default
// template when tmpl is empty. The result has trailing whitespace trimmed. A
// template that fails to parse, fails to execute, or renders to nothing returns
// an error so the caller can warn and fall back to the default.
func RenderCommitMessage(tmpl string, data CommitData) (string, error) {
	if strings.TrimSpace(tmpl) == "" {
		tmpl = DefaultCommitMessageTemplate
	}
	t, err := template.New("commit").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parsing commit message template: %w", err)
	}
	var b strings.Builder
	if err := t.Execute(&b, data); err != nil {
		return "", fmt.Errorf("rendering commit message template: %w", err)
	}
	msg := strings.TrimRight(b.String(), " \t\r\n")
	if strings.TrimSpace(msg) == "" {
		return "", fmt.Errorf("commit message template produced empty output")
	}
	return msg, nil
}
