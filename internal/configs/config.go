package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Route struct {
	Name            string `json:"name"`
	Endpoint        string `json:"endpoint"`
	Port            string `json:"port"`
	Destination_URL string `json:"destination_url"`
}
type Configuration struct {
	Server struct {
		Host        string `json:"host"`
		Listen_port string `json:"listen_port"`
	} `json:"server"`
	Routes []Route `json:"routes"`
}

func NewConfiguration(configFilename string) (*Configuration, error) {
	fmt.Printf("Load config file: %s\n", configFilename)

	contents, err := os.ReadFile(configFilename)
	if err != nil {
		return nil, fmt.Errorf("could not load configuration file: %v", err)
	}

	var config Configuration

	decoder := json.NewDecoder(strings.NewReader(string(contents)))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		return nil, fmt.Errorf("failed to parse JSON config: %v", err)
	}

	err = validateConfig(&config)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Config loaded Successfully\n\n")
	return &config, nil
}

func validateConfig(config *Configuration) error {
	if strings.TrimSpace(config.Server.Host) == "" {
		return fmt.Errorf("server.host not defined")
	}
	if strings.TrimSpace(config.Server.Listen_port) == "" {
		return fmt.Errorf("server.listen_port not defined")
	}
	for i, route := range config.Routes {
		if strings.TrimSpace(route.Name) == "" {
			return fmt.Errorf("routes[%v].name not defined", i)
		}
		if strings.TrimSpace(route.Endpoint) == "" {
			return fmt.Errorf("routes[%v].endpoint not defined", i)
		}
		if strings.TrimSpace(route.Destination_URL) == "" {
			return fmt.Errorf("routes[%v].destination_url not defined", i)
		}
	}
	return nil
}
