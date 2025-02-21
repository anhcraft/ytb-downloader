package settings

type Profile struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

func defaultProfile() Profile {
	return Profile{
		Name: "Default",
		Path: "settings.json",
	}
}

type ProfileCollection struct {
	Profiles []Profile `json:"profiles,omitempty"`
	Selected int       `json:"selected,omitempty"`
}

func NewProfileCollection() *ProfileCollection {
	return &ProfileCollection{
		Profiles: []Profile{
			defaultProfile(),
		},
		Selected: 0,
	}
}

func (p *ProfileCollection) Normalize() {
	if len(p.Profiles) == 0 {
		p.AddProfile(defaultProfile())
	}
	if p.Selected < 0 {
		p.Selected = 0
	}
	if p.Selected >= len(p.Profiles) {
		p.Selected = len(p.Profiles) - 1
	}
}

func (p *ProfileCollection) AddProfile(profile Profile) {
	p.Profiles = append(p.Profiles, profile)
}

func (p *ProfileCollection) SelectProfile(index int) {
	if index < 0 || index >= len(p.Profiles) {
		return
	}
	p.Selected = index
}

func (p *ProfileCollection) GetSelectedProfile() Profile {
	if p.Selected >= 0 && p.Selected < len(p.Profiles) {
		return p.Profiles[p.Selected]
	}
	return defaultProfile()
}

func (p *ProfileCollection) DeleteProfile(i int) {
	if i < 0 || i >= len(p.Profiles) {
		return
	}
	p.Profiles = append(p.Profiles[:i], p.Profiles[i+1:]...)
}
