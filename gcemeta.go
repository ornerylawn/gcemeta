// Package gcemeta provides access to Google Compute Engine
// metadata. It's useful for discovering things about the machine on
// which your process is running (e.g. instance name, project name,
// zone).
package gcemeta

import (
	"encoding/json"
	"net/http"
)

type Metadata struct {
	Instance *Instance `json:"instance"`
	Project  *Project  `json:"project"`
}

type Instance struct {
	Attributes        map[string]string   `json:"attributes"`
	Description       string              `json:"description"`
	Disks             []*Disk             `json:"disks"`
	Hostname          string              `json:"hostname"`
	ID                int64               `json:"id"`
	Image             string              `json:"image"`
	MachineType       string              `json:"machineType"`
	MaintenanceEvent  string              `json:"maintenanceEvent"`
	NetworkInterfaces []*NetworkInterface `json:"networkInterfaces"`
	Scheduling        *Scheduling         `json:"scheduling"`
	Tags              []string            `json:"tags"`
	Zone              string              `json:"zone"`
}

type Disk struct {
	DeviceName string `json:"deviceName"`
	Index      int    `json:"index"`
	Mode       string `json:"mode"`
	Type       string `json:"type"`
}

type NetworkInterface struct {
	AccessConfigs []*AccessConfig `json:"accessConfigs"`
	ForwardedIPs  []string        `json:"forwardedIps"`
	IP            string          `json:"ip"`
	Network       string          `json:"network"`
}

type AccessConfig struct {
	ExternalIP string `json:"externalIp"`
	Type       string `json:"type"`
}

type Scheduling struct {
	AutomaticRestart  string `json:"automaticRestart"`
	OnHostMaintenance string `json:"onHostMaintenance"`
}

type Project struct {
	Attributes       map[string]string `json:"attributes"`
	NumericProjectID int64             `json:"numericProjectId"`
	ProjectID        string            `json:"projectId"`
}

const url = "http://metadata.google.internal/computeMetadata/v1/?recursive=true"

// Get requests metadata from the metadata server.
func Get() (*Metadata, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var m *Metadata
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
