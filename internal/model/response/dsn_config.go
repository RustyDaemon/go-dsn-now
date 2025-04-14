package response

type DSNConfig struct {
	Sites       []DSNConfigSite       `xml:"sites>site"`
	Spacecrafts []DSNConfigSpacecraft `xml:"spacecraftMap>spacecraft"`
}

type DSNConfigSite struct {
	Name         string              `xml:"name,attr"`
	FriendlyName string              `xml:"friendlyName,attr"`
	Longitude    string              `xml:"longitude,attr"`
	Latitude     string              `xml:"latitude,attr"`
	Dishes       []DSNConfigSiteDish `xml:"dish"`
}

type DSNConfigSiteDish struct {
	Name         string `xml:"name,attr"`
	FriendlyName string `xml:"friendlyName,attr"`
	Type         string `xml:"type,attr"`
}

type DSNConfigSpacecraft struct {
	Name         string `xml:"name,attr"`
	ExplorerName string `xml:"explorerName,attr"`
	FriendlyName string `xml:"friendlyName,attr"`
}
