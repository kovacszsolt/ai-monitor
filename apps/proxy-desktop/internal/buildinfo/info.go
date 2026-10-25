package buildinfo

const AppName = "AI Proxy Manager"

// License is the SPDX license identifier shown in About (see LICENSE in app root).
const License = "MIT"

// Version is set at link time via -ldflags (see Makefile); "dev" for local builds.
var Version = "dev"

// BuildTime is set at link time via -ldflags (see Makefile).
var BuildTime = "dev"
