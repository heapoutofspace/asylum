package kit

func init() {
	Register(&Kit{
		Name:        "apt",
		Description: "Extra apt packages installed in the project image",
		ConfigSnippet: `  # apt:                # System packages installed via apt-get
  #   packages:
  #     - imagemagick
  #     - ffmpeg
`,
	})
}
