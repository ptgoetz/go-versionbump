package config

const DefaultConfigTemplate = `# The current version of the project. This is the source of truth for the project version. 
# Set this once and let VersionBump manage it.
version: "{{.Version}}" 

# Git configuration (optional)
git-commit: {{ .GitCommit }} # Whether to create a git commit for the version bump.
git-sign: {{ .GitSign }}   # Whether to sign the git commit and tag.
git-tag: {{ .GitTag }}    # Whether to create a git tag for the version bump.

# Git commit and tag templates. These are the templates used for the git commit and tag messages.
git-commit-template: "{{ .GitCommitTemplate }}" # The template for the git commit message.
git-tag-template: "{{ .GitTagTemplate }}" # The template for the git tag name.
git-tag-message-template: "{{ .GitTagMessageTemplate }}" # The template for the git tag message.

# Prerelease labels. These are the labels that will be used for prerelease versions.
# VersionBump with sort these labels in ascending order when determining the next version.
# If the bump type is 'prerelease-next'', the next label will be used. Attempting to bump past the last label 
# will result in an error.
prerelease-labels:
{{- range .PreReleaseLabels }}
  - "{{ . }}"
{{- end }}

# The build label. This is the label that will be used for build versions.
build-label: "{{.BuildLabel}}"

files:  # The files to update with the new version (i.e. "Tracked files").
# The following example will replace all occurrences of the old version with the new version in the README.md file.
#   - path: "README.md"
#    replace:
#      - "v{version}"
{{- range .Files }}
  - path: "{{ .Path }}"
    replace:
    {{- range .Replace }}
      - "{{ . }}"
    {{- end }}
{{- end }}
`
