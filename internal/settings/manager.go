package settings

import (
	"errors"
	"sync"
)

const defaultProfileCollectionPath = "profiles.json"

var profileCollection *ProfileCollection
var currentSettings *Settings
var currentProfile Profile
var lock sync.RWMutex

func InitManager() {
	lock.Lock()
	defer lock.Unlock()
	profileCollection = loadProfileCollection(defaultProfileCollectionPath)
	currentProfile = profileCollection.GetSelectedProfile()
	currentSettings = retrieveSettings(currentProfile.Path)
}

func Get() *Settings {
	lock.RLock()
	defer lock.RUnlock()
	return currentSettings
}

func Save() {
	SaveSettings(currentProfile, currentSettings)
}

func LoadSettings(profile Profile) *Settings {
	lock.Lock()
	defer lock.Unlock()
	if profile.Name == currentProfile.Name {
		return currentSettings
	} else {
		return retrieveSettings(profile.Path)
	}
}

func ResetSettings(profile Profile) *Settings {
	lock.Lock()
	defer lock.Unlock()

	newSettings := NewSettings()

	if profile.Name == currentProfile.Name {
		currentSettings = newSettings
	}

	persistSettings(profile.Path, newSettings)
	return newSettings
}

func SaveSettings(profile Profile, settings *Settings) {
	lock.Lock()
	defer lock.Unlock()
	if profile.Name == currentProfile.Name {
		currentSettings = settings
	}
	persistSettings(profile.Path, settings)
}

func GetProfileNames() []string {
	lock.RLock()
	defer lock.RUnlock()

	names := make([]string, len(profileCollection.Profiles))
	for i, p := range profileCollection.Profiles {
		names[i] = p.Name
	}
	return names
}

func GetProfiles() []Profile {
	lock.RLock()
	defer lock.RUnlock()
	return profileCollection.Profiles
}

func GetProfile() Profile {
	lock.RLock()
	defer lock.RUnlock()
	return currentProfile
}

func SelectProfile(value string) {
	lock.Lock()
	defer lock.Unlock()
	if currentProfile.Name == value {
		return
	}

	for i, p := range profileCollection.Profiles {
		if p.Name == value {
			profileCollection.SelectProfile(i)
			saveProfileCollection(defaultProfileCollectionPath, profileCollection)
			currentProfile = profileCollection.GetSelectedProfile()
			currentSettings = retrieveSettings(currentProfile.Path)
			break
		}
	}
}

func DeleteProfile(name string) {
	lock.Lock()
	defer lock.Unlock()
	if currentProfile.Name == name {
		return
	}
	for i, p := range profileCollection.Profiles {
		if p.Name == name {
			profileCollection.DeleteProfile(i)
			saveProfileCollection(defaultProfileCollectionPath, profileCollection)
			break
		}
	}
}

func AddProfile(profile Profile) error {
	lock.Lock()
	defer lock.Unlock()
	for _, p := range profileCollection.Profiles {
		if p.Name == profile.Name {
			return errors.New("duplicated profile name")
		}
	}
	profileCollection.AddProfile(profile)
	saveProfileCollection(defaultProfileCollectionPath, profileCollection)
	return nil
}
