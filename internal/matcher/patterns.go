package matcher

// Cloud service identifiers
const (
	ServiceDrive     = "drive"
	ServiceSharePoint = "sharepoint"
	ServiceDropbox   = "dropbox"
	ServiceOneDrive  = "onedrive"
	ServiceBox       = "box"
	ServiceICloud    = "icloud"
)

// Pattern defines a regex pattern for a cloud storage service
type Pattern struct {
	Service string
	Regex   string
}

// AllPatterns returns all cloud storage URL patterns
func AllPatterns() []Pattern {
	return []Pattern{
		// Google Drive - drive.google.com, docs.google.com
		{
			Service: ServiceDrive,
			Regex:   `https?://(?:drive|docs)\.google\.com/[^\s"'>]+`,
		},
		// SharePoint - *.sharepoint.com (my, team subdomains)
		{
			Service: ServiceSharePoint,
			Regex:   `https?://(?:[a-z0-9\-]+\.)?(?:my\.|team\.)?[a-z0-9\-]+\.sharepoint\.com/[^\s"'>]+`,
		},
		// Dropbox - dropbox.com, dropboxusercontent.com, db.tt
		{
			Service: ServiceDropbox,
			Regex:   `https?://(?:[^\s"'>]+\.)?(?:dropbox\.com|dropboxusercontent\.com|db\.tt)/[^\s"'>]+`,
		},
		// OneDrive - onedrive.live.com, 1drv.ms
		{
			Service: ServiceOneDrive,
			Regex:   `https?://(?:onedrive\.live\.com|1drv\.ms)/[^\s"'>]+`,
		},
		// Box - app.box.com, *.box.com
		{
			Service: ServiceBox,
			Regex:   `https?://(?:app\.)?[a-z0-9\-]*\.?box\.com/[^\s"'>]+`,
		},
		// iCloud - icloud.com shared links
		{
			Service: ServiceICloud,
			Regex:   `https?://(?:[^\s"'>]+\.)?icloud\.com/[^\s"'>]+`,
		},
	}
}
