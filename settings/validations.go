package settings

type TypeAllowLists struct {
	SemVerMinorIncreasers []string `yaml:"semver_minor_increases"`
	SemverPatchIncreasers []string `yaml:"semver_patch_increases"`
	SemverNotAffects      []string `yaml:"semver_is_not_affected"`
}

func (a TypeAllowLists) GetAllTypes() []string {
	list := []string{}
	list = append(list, a.SemVerMinorIncreasers...)
	list = append(list, a.SemverPatchIncreasers...)
	list = append(list, a.SemverNotAffects...)
	return list
}

type ValidationScheme struct {
	Sections struct {
		Hook struct {
			CommitMsg struct {
				Enabled bool `yaml:"enabled"`
			} `yaml:"commitMsg"`
		} `yaml:"hook"`
		// TODO, add ability to disable Changelog validations?
		// Changelog struct {
		// 	Enabled bool `yaml:"enabled"`
		// } `yaml:"changelog"`
	} `yaml:"sections"`
	Rules struct {
		Issue struct {
			Present bool `yaml:"present"`
		} `yaml:"issue"`
		Header struct {
			MaxLength int `yaml:"maxLength"`
			Type      struct {
				Lowercase  bool           `yaml:"lowercase"`
				Allowlists TypeAllowLists `yaml:"allowlists"`
			} `yaml:"type"`
			Scope struct {
				Present   bool     `yaml:"present"`
				Lowercase bool     `yaml:"lowercase"`
				Allowlist []string `yaml:"allowlist"`
			} `yaml:"scope"`
			Subject struct {
				MinWords int `yaml:"minWords"`
			} `yaml:"subject"`
		} `yaml:"header"`
	} `yaml:"rules"`
}
