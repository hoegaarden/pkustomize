package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type ObjectMeta struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name        string            `yaml:"name"`
		Labels      map[string]string `yaml:"labels,omitempty"`
		Annotations map[string]string `yaml:"annotations,omitempty"`
	} `yaml:"metadata"`
}

type pluginConfig struct {
	ObjectMeta `yaml:",inline"`
	CF         struct {
		Org      string   `yaml:"org"`
		Space    string   `yaml:"space"`
		Services []string `yaml:"services"`
	} `yaml:"cf"`
}

type outputResource struct {
	ObjectMeta `yaml:",inline"`
	Data       struct {
		VcapSerices string `yaml:"vcap-service"`
	} `yaml:"data"`
}

const (
	pluginConfigEnv = "KUSTOMIZE_PLUGIN_CONFIG_STRING"
	mockWorkDelay   = 500 * time.Millisecond
)

var errPluginConfigEnv = fmt.Errorf("env var %q not set", pluginConfigEnv)

func getPluginConfig() (pluginConfig, error) {
	c, ok := os.LookupEnv(pluginConfigEnv)
	if !ok {
		return pluginConfig{}, errPluginConfigEnv
	}

	p := pluginConfig{}
	if err := yaml.Unmarshal([]byte(c), &p); err != nil {
		return pluginConfig{}, err
	}

	return p, nil
}

// genSecret is a function which mocks what you cold do in real live: go and
// talk to external services and fetch some data.
//
// One example, seen at a customer:
// They need to talk to a PAS, pulling in some services data from a specific
// org/space, massage the data a bit and finally put it into a k8s secret.
func genSecret(c pluginConfig) (outputResource, error) {
	org, space := c.CF.Org, c.CF.Space
	res := []string{}

	for _, s := range c.CF.Services {
		log.Printf("doing something in %v/%v for service %v ...\n", org, space, s)
		time.Sleep(mockWorkDelay)

		res = append(res, fmt.Sprintf("did something in %v/%v for service %v", org, space, s))
	}

	secret := outputResource{ObjectMeta: ObjectMeta{
		APIVersion: "v1",
		Kind:       "Secret",
		Metadata:   c.Metadata,
	}}
	secret.Data.VcapSerices = base64.StdEncoding.EncodeToString([]byte(
		strings.Join(res, "\n"),
	))

	return secret, nil
}

func bailOnErr(err error) {
	if err == nil {
		return
	}

	log.Fatal(err)
}

func main() {
	c, err := getPluginConfig()
	bailOnErr(err)

	s, err := genSecret(c)
	bailOnErr(err)

	err = yaml.NewEncoder(os.Stdout).Encode(s)
	bailOnErr(err)
}
