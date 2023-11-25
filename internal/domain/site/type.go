package site

// Site represents a site in the build. This is currently a very narrow interface,
// but the actual implementation will be richer, see hugolib.SiteInfo.
type Site interface {

	// Title Returns the configured title for this Site.
	Title() string
}
