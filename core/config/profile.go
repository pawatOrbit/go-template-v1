package core_config

import (
	"html/template"
	"strings"

	"github.com/pawatOrbit/ai-mock-data-service/go/utils/runtime"
)

type Profile struct {
	Env runtime.Environment `mapstructure:"env"`
}

const configGlobalTemplate = "config/config.{{.Env}}.yaml"

func GetGlobalConfigFilePath(cfg runtime.RuntimeCfg) (string, error) {
	var result strings.Builder

	tmpl, err := template.New("configGlobalTemplate").Parse(configGlobalTemplate)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&result, Profile{
		Env: cfg.Env,
	})
	if err != nil {
		return "", err
	}

	// Store the result in a variable
	configName := result.String()
	return configName, nil
}
