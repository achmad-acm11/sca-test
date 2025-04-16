package dto

import "time"

type ResultOutputFile struct {
	SchemaVersion int    `json:"SchemaVersion"`
	ArtifactName  string `json:"ArtifactName"`
	ArtifactType  string `json:"ArtifactType"`
	Metadata      struct {
		ImageConfig struct {
			Architecture string    `json:"architecture"`
			Created      time.Time `json:"created"`
			Os           string    `json:"os"`
			Rootfs       struct {
				Type    string      `json:"type"`
				DiffIds interface{} `json:"diff_ids"`
			} `json:"rootfs"`
			Config struct {
			} `json:"config"`
		} `json:"ImageConfig"`
	} `json:"Metadata"`
	Results []struct {
		Target          string `json:"Target"`
		Class           string `json:"Class"`
		Type            string `json:"Type"`
		Vulnerabilities []struct {
			VulnerabilityID  string `json:"VulnerabilityID"`
			PkgName          string `json:"PkgName"`
			InstalledVersion string `json:"InstalledVersion"`
			FixedVersion     string `json:"FixedVersion"`
			Layer            struct {
			} `json:"Layer"`
			SeveritySource string   `json:"SeveritySource"`
			PrimaryURL     string   `json:"PrimaryURL"`
			Title          string   `json:"Title"`
			Description    string   `json:"Description"`
			Severity       string   `json:"Severity"`
			CweIDs         []string `json:"CweIDs"`
			Cvss           struct {
				Nvd struct {
					V2Vector string  `json:"V2Vector"`
					V3Vector string  `json:"V3Vector"`
					V2Score  float64 `json:"V2Score"`
					V3Score  float64 `json:"V3Score"`
				} `json:"nvd"`
			} `json:"CVSS"`
			References       []string  `json:"References"`
			PublishedDate    time.Time `json:"PublishedDate"`
			LastModifiedDate time.Time `json:"LastModifiedDate"`
		} `json:"Vulnerabilities"`
	} `json:"Results"`
}
