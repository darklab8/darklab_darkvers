package semanticgit

import (
	"autogit/semanticgit/git"
	"autogit/semanticgit/semver"
	"testing"

	"github.com/stretchr/testify/assert"
)

// For debug on workdir
// func TestIngeration(t *testing.T) {
// 	git := (&git.Repository{}).NewRepoIntegration()
// 	gitSemantic := (&SemanticGit{}).NewRepo(git)
// 	vers := gitSemantic.CalculateNextVersion(gitSemantic.GetCurrentVersion()).ToString()
// 	assert.Equal(t, "v0.2.0", vers)
// }

func TestCurrentNextRegularVersion(t *testing.T) {
	gitInMemory := (&git.Repository{}).TestNewRepo()
	gitSemantic := (&SemanticGit{}).NewRepo(gitInMemory)
	gitInMemory.TestCommit("feat: init")

	assert.Equal(t, "v0.0.0", gitSemantic.GetCurrentVersion().ToString())

	gitInMemory.TestCreateTag("v0.0.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: test2")

	assert.Equal(t, "v0.0.1", gitSemantic.GetCurrentVersion().ToString())

	assert.Equal(t, "v0.1.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())

	gitInMemory.TestCreateTag("v0.1.0", gitInMemory.TestCommit("fix: thing"))

	gitInMemory.TestCommit("fix: test2")
	assert.Equal(t, "v0.1.1", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())

	gitInMemory.TestCommit("feat: test2")

	assert.Equal(t, "v0.2.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())

	// Semantic version should be same if no new comments
	gitInMemory.TestCreateTag("v0.2.0", gitInMemory.TestCommit("feat: new thing"))
	assert.Equal(t, "v0.2.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())
}

func TestGetChangelogs(t *testing.T) {
	gitInMemory := (&git.Repository{}).TestNewRepo()
	gitSemantic := (&SemanticGit{}).NewRepo(gitInMemory)
	gitInMemory.TestCommit("feat: init")

	gitInMemory.TestCommit("feat: test3")
	gitInMemory.TestCommit("feat: test5")
	gitInMemory.TestCreateTag("v0.0.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat(api): test")
	gitInMemory.TestCreateTag("v0.0.2", gitInMemory.TestCommit("feat(api): test2"))
	gitInMemory.TestCommit("fix: test1")
	gitInMemory.TestCommit("fix: test2")
	gitInMemory.TestCommit("fix: test3")

	logs1 := gitSemantic.GetChangelogByTag("", true).Logs
	assert.Len(t, logs1, 3)

	logs2 := gitSemantic.GetChangelogByTag("v0.0.2", true).Logs
	assert.Len(t, logs2, 2)

	logs3 := gitSemantic.GetChangelogByTag("v0.0.1", true).Logs
	assert.Len(t, logs3, 4)
}

func TestTestPrereleaseVersions(t *testing.T) {
	gitInMemory := (&git.Repository{}).TestNewRepo()
	gitSemantic := (&SemanticGit{}).NewRepo(gitInMemory)

	gitInMemory.TestCommit("feat: init")
	assert.Equal(t, "v0.1.0-a.1", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true}).ToString())

	gitInMemory.TestCreateTag("v0.1.0-a.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: thing")
	assert.Equal(t, "v0.1.0-a.2", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true}).ToString())

	gitInMemory.TestCommit("feat: test5")
	gitInMemory.TestCreateTag("v0.1.0", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("feat: thing")
	assert.Equal(t, "v0.2.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())
	assert.Equal(t, "v0.2.0-a.1.b.1", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0-a.1", gitInMemory.TestCommit("fix: thing"))
	assert.Equal(t, "v0.2.0", gitSemantic.GetNextVersion(semver.OptionsSemVer{}).ToString())
	assert.Equal(t, "v0.2.0-a.1.b.1", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0-a.1.b.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("fix: thing")
	assert.Equal(t, "v0.2.0-a.1.b.2", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0-a.1.b.2", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("fix: thing")
	assert.Equal(t, "v0.2.0-a.1.b.3", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true, Beta: true}).ToString())

	gitInMemory.TestCreateTag("v0.2.0-rc.1", gitInMemory.TestCommit("fix: thing"))
	gitInMemory.TestCommit("fix: thing")
	assert.Equal(t, "v0.2.0-rc.2", gitSemantic.GetNextVersion(semver.OptionsSemVer{Rc: true}).ToString())
	assert.Equal(t, "v0.2.0-a.1.b.3", gitSemantic.GetNextVersion(semver.OptionsSemVer{Alpha: true, Beta: true}).ToString())
}
