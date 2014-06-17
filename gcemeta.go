// Package gcemeta provides access to Google Compute Engine
// metadata. It's useful for discovering things about the machine on
// which your process is running (e.g. instance name, project name,
// zone).
//
//  meta, err := gcemeta.Get()
//  ...
//  fmt.Println(meta.Project.ProjectID)
//
// Documentation on each of the fields can be found here:
// https://developers.google.com/compute/docs/metadata
package gcemeta

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
)

type Meta struct {
	Instance *Instance `json:"instance"`
	Project  *Project  `json:"project"`
}

// InstanceURL computes the fully qualified URL of the instance,
// intended to be used with the Google Compute Engine API.
func (m *Meta) InstanceURL() string {
	return fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s",
		m.Project.ID, m.Instance.ShortZone(), m.Instance.Name())
}

type Instance struct {
	Attributes        map[string]string   `json:"attributes"`
	Description       string              `json:"description"`
	Disks             []*Disk             `json:"disks"`
	Hostname          string              `json:"hostname"`
	ID                *big.Int            `json:"id"`
	Image             string              `json:"image"`
	MachineType       string              `json:"machineType"`
	MaintenanceEvent  string              `json:"maintenanceEvent"`
	NetworkInterfaces []*NetworkInterface `json:"networkInterfaces"`
	Scheduling        *Scheduling         `json:"scheduling"`
	Tags              []string            `json:"tags"`
	Zone              string              `json:"zone"`
}

// Name parses the instance's name from its hostname.
func (i *Instance) Name() string {
	j := strings.Index(i.Hostname, ".")
	if j < 0 {
		return i.Hostname
	}
	return i.Hostname[:j]
}

// ShortZone parses the zone for the common short representation
// (e.g. "us-central1-a").
func (i *Instance) ShortZone() string {
	j := strings.LastIndex(i.Zone, "/")
	if j < 0 {
		return i.Zone
	}
	return i.Zone[j+1:]
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
	Attributes map[string]string `json:"attributes"`
	ID         string            `json:"projectId"`
	NumericID  *big.Int          `json:"numericProjectId"`
}

const url = "http://metadata.google.internal/computeMetadata/v1/?recursive=true"

// Get requests metadata from the metadata server.
func Get() (*Meta, error) {
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
	var m *Meta
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
