package analyzer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
	"github.com/stretchr/testify/require"
)

func compareCommit(c *semrel.Commit, t, s string, change *semrel.Change) bool {
	if c.Type != t || c.Scope != s {
		return false
	}
	if c.Change.Major != change.Major ||
		c.Change.Minor != change.Minor ||
		c.Change.Patch != change.Patch {
		return false
	}
	return true
}

func createRawCommit(sha, message string) *semrel.RawCommit {
	return &semrel.RawCommit{
		SHA:        sha,
		RawMessage: message,
		Annotations: map[string]string{
			"author_name": "test",
		},
	}
}

func TestAnnotations(t *testing.T) {
	defaultAnalyzer := &DefaultCommitAnalyzer{}
	rawCommit := createRawCommit("a", "fix: bug #123 and #243\nthanks @Test-user for providing this fix\n\nCloses #22")
	commit := defaultAnalyzer.analyzeSingleCommit(rawCommit)
	require.Equal(t, rawCommit.SHA, commit.SHA)
	require.Equal(t, rawCommit.RawMessage, strings.Join(commit.Raw, "\n"))
	require.Equal(t, "test", commit.Annotations["author_name"])
	require.Equal(t, "123,243,22", commit.Annotations["mentioned_issues"])
	require.Equal(t, "Test-user", commit.Annotations["mentioned_users"])
}

func TestDefaultAnalyzer(t *testing.T) {
	testCases := []struct {
		RawCommit *semrel.RawCommit
		Type      string
		Scope     string
		Change    *semrel.Change
	}{
		{
			createRawCommit("a", "feat: new feature"),
			"feat",
			"",
			&semrel.Change{Major: false, Minor: true, Patch: false},
		},
		{
			createRawCommit("b", "feat(web): new feature"),
			"feat",
			"web",
			&semrel.Change{Major: false, Minor: true, Patch: false},
		},
		{
			createRawCommit("c", "new feature"),
			"",
			"",
			&semrel.Change{Major: false, Minor: false, Patch: false},
		},
		{
			createRawCommit("d", "chore: break\nBREAKING CHANGE: breaks everything"),
			"chore",
			"",
			&semrel.Change{Major: true, Minor: false, Patch: false},
		},
		{
			createRawCommit("e", "feat!: modified login endpoint"),
			"feat",
			"",
			&semrel.Change{Major: true, Minor: false, Patch: false},
		},
		{
			createRawCommit("f", "fix!: fixed a typo"),
			"fix",
			"",
			&semrel.Change{Major: true, Minor: false, Patch: false},
		},
		{
			createRawCommit("g", "refactor(parser)!: drop support for Node 6\n\nBREAKING CHANGE: refactor to use JavaScript features not available in Node 6."),
			"refactor",
			"parser",
			&semrel.Change{Major: true, Minor: false, Patch: false},
		},
		{
			createRawCommit("h", "docs: added more documentation"),
			"docs",
			"",
			&semrel.Change{Major: false, Minor: false, Patch: false},
		},
		{
			createRawCommit("i", "chore: moved README.md to root"),
			"chore",
			"",
			&semrel.Change{Major: false, Minor: false, Patch: false},
		},
	}

	defaultAnalyzer := &DefaultCommitAnalyzer{}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("AnalyzeCommitMessage: %s", tc.RawCommit.RawMessage), func(t *testing.T) {
			require.True(t, compareCommit(defaultAnalyzer.analyzeSingleCommit(tc.RawCommit), tc.Type, tc.Scope, tc.Change))
		})
	}
}
