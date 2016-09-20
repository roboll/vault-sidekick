/*
Copyright 2015 Home Office All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"regexp"
	"time"
)

const (
	// optionFilename option to set the filename of the resource
	optionFilename = "file"
	// optionFormat set the output format (yaml, xml, json)
	optionFormat = "fmt"
	// optionTemplatePath is the full path to a template
	optionTemplatePath = "tpl"
	// optionRenewal sets the duration to renew the resource
	optionRenewal = "renew"
	// optionRevoke revokes an old lease when retrieving a new one
	optionRevoke = "revoke"
	// optionRevokeDelay
	optionsRevokeDelay = "delay"
	// optionUpdate overrides the lease of the resource
	optionUpdate = "update"
	// optionsExec executes something on a change
	optionExec = "exec"
	// optionCreate creates a secret if it doesn't exist
	optionCreate = "create"
	// optionSize sets the initial size of a password secret
	optionSize = "size"
	// defaultSize sets the default size of a generic secret
	defaultSize = 20
)

var (
	resourceFormatRegex = regexp.MustCompile("^(yaml|yml|json|env|awsenv|ini|txt|cert|bundle|csv)$")

	// a map of valid resource to retrieve from vault
	validResources = map[string]bool{
		"raw":       true,
		"pki":       true,
		"aws":       true,
		"secret":    true,
		"mysql":     true,
		"tpl":       true,
		"postgres":  true,
		"transit":   true,
		"cubbyhole": true,
		"cassandra": true,
	}
)

func defaultVaultResource() *VaultResource {
	return &VaultResource{
		format:    "yaml",
		renewable: false,
		revoked:   false,
		options:   make(map[string]string, 0),
		size:      defaultSize,
	}
}

// VaultResource is the structure which defined a resource set from vault
type VaultResource struct {
	// the namespace of the resource
	resource string
	// the name of the resource
	path string
	// the format of the resource
	format string
	// whether the resource should be renewed?
	renewable bool
	// whether the resource should be revoked?
	revoked bool
	// the revoke delay
	revokeDelay time.Duration
	// the lease duration
	update time.Duration
	// whether the resource should be created?
	create bool
	// the size of a secret to create
	size int64
	// the filename to save the secret
	filename string
	// the template file
	templateFile string
	// the path to an exec to run on a change
	execPath string
	// additional options to the resource
	options map[string]string
}

// GetFilename generates a resource filename by default the resource name and resource type, which
// can override by the OPTION_FILENAME option
func (r VaultResource) GetFilename() string {
	if r.filename != "" {
		return r.filename
	}

	return fmt.Sprintf("%s.%s", r.path, r.resource)
}

// IsValid checks to see if the resource is valid
func (r *VaultResource) IsValid() error {
	// step: check the resource type
	if _, found := validResources[r.resource]; !found {
		return fmt.Errorf("unsupported resource type: %s", r.resource)
	}

	// step: check is have all the required options to this resource type
	if err := r.isValidResource(); err != nil {
		return fmt.Errorf("invalid resource: %s, %s", r, err)
	}

	return nil
}

// isValidResource validates the resource meets the requirements
func (r *VaultResource) isValidResource() error {
	switch r.resource {
	case "pki":
		if _, found := r.options["common_name"]; !found {
			return fmt.Errorf("pki resource requires a common name specified")
		}
	case "transit":
		if _, found := r.options["ciphertext"]; !found {
			return fmt.Errorf("transit requires a ciphertext option")
		}
	case "tpl":
		if _, found := r.options[optionTemplatePath]; !found {
			return fmt.Errorf("template resource requires a template path option")
		}
	}

	return nil
}

// String returns a string representation of the struct
func (r VaultResource) String() string {
	return fmt.Sprintf("type: %s, path:%s", r.resource, r.path)
}
