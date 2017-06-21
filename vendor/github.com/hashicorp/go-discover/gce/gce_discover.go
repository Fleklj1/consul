// Package gce provides node discovery for Google Cloud.
package gce

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/go-discover/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

// Discover returns the private ip addresses of all Google Cloud
// instances in some or all zones of a project with a certain tag value.
//
// cfg contains the configuration in "key=val key=val ..." format. The
// values are URL encoded.
//
// The supported keys are:
//
//   project_name     : The name of the project. discovered if not set
//   zone_pattern     : A RE2 regular expression for filtering zones, e.g. us-west1-.*, or us-(?west|east).*
//   tag_value        : The tag value for filtering instances
//   credentials_file : The path to the credentials file. See below for more details
//
// Authentication is handled in the following order:
//
//  1. Use credentials from "credentials_file", if provided.
//  2. Use JSON file from GOOGLE_APPLICATION_CREDENTIALS environment variable.
//  3. Use JSON file in a location known to the gcloud command-line tool.
//     On Windows, this is %APPDATA%/gcloud/application_default_credentials.json.
//     On other systems, $HOME/.config/gcloud/application_default_credentials.json.
//  4. On Google Compute Engine, use credentials from the metadata
//     server. In this final case any provided scopes are ignored.
//
// Example:
//
//  project_name=test zone_pattern=us-(?west|east).* tag_value=consul-server credentials_file=xxx
//
func Discover(cfg string, l *log.Logger) ([]string, error) {
	m, err := config.Parse(cfg)
	if err != nil {
		return nil, err
	}

	project := m["project_name"]
	zone := m["zone_pattern"]
	creds := m["credentials_file"]
	tagValue := m["tag_value"]

	// determine the project name
	if project == "" {
		l.Println("[INFO] discover-gce: Looking up project name")
		p, err := lookupProject()
		if err != nil {
			return nil, err
		}
		project = p
	}
	l.Printf("[INFO] discover-gce: Project name is %q", project)

	// create an authenticated client
	if creds != "" {
		l.Printf("[INFO] discover-gce: Loading credentials from %s", creds)
	}
	client, err := client(creds)
	if err != nil {
		return nil, err
	}
	svc, err := compute.New(client)
	if err != nil {
		return nil, err
	}

	// lookup the project zones to look in
	if zone != "" {
		l.Printf("[INFO] discover-gce: Looking up zones matching %s", zone)
	} else {
		l.Printf("[INFO] discover-gce: Looking up all zones")
	}
	zones, err := lookupZones(svc, project, zone)
	if err != nil {
		return nil, err
	}
	l.Printf("[INFO] discover-gce: Found zones %v", zones)

	// lookup the instance addresses
	var addrs []string
	for _, zone := range zones {
		a, err := lookupAddrs(svc, project, zone, tagValue)
		if err != nil {
			return nil, err
		}
		l.Printf("[INFO] discover-gce: Zone %q has %v", zone, a)
		addrs = append(addrs, a...)
	}
	return addrs, nil
}

// client returns an authenticated HTTP client for use with GCE.
func client(path string) (*http.Client, error) {
	if path == "" {
		return google.DefaultClient(oauth2.NoContext, compute.ComputeScope)
	}

	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	jwtConfig, err := google.JWTConfigFromJSON(key, compute.ComputeScope)
	if err != nil {
		return nil, err
	}

	return jwtConfig.Client(oauth2.NoContext), nil
}

// lookupProject retrieves the project name from the metadata of the current node.
func lookupProject() (string, error) {
	req, err := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Metadata-Flavor", "Google")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("discover-gce: invalid status code %d when fetching project id", resp.StatusCode)
	}

	project, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(project), nil
}

// lookupZones retrieves the zones of the project and filters them by pattern.
func lookupZones(svc *compute.Service, project, pattern string) ([]string, error) {
	call := svc.Zones.List(project)
	if pattern != "" {
		call = call.Filter("name eq " + pattern)
	}

	var zones []string
	f := func(page *compute.ZoneList) error {
		for _, v := range page.Items {
			zones = append(zones, v.Name)
		}
		return nil
	}

	if err := call.Pages(oauth2.NoContext, f); err != nil {
		return nil, err
	}
	return zones, nil
}

// lookupAddrs retrieves the private ip addresses of all instances in a given
// project and zone which have a matching tag value.
func lookupAddrs(svc *compute.Service, project, zone, tag string) ([]string, error) {
	var addrs []string
	f := func(page *compute.InstanceList) error {
		for _, v := range page.Items {
			if len(v.NetworkInterfaces) == 0 || v.NetworkInterfaces[0].NetworkIP == "" {
				continue
			}
			for _, t := range v.Tags.Items {
				if t == tag {
					addrs = append(addrs, v.NetworkInterfaces[0].NetworkIP)
					break
				}
			}
		}
		return nil
	}

	call := svc.Instances.List(project, zone)
	if err := call.Pages(oauth2.NoContext, f); err != nil {
		return nil, err
	}
	return addrs, nil
}
