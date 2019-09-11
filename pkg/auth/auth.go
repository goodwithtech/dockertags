package auth

import (
	"fmt"
	"strings"

	"github.com/goodwithtech/image-tag-sorter/pkg/log"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/types"
)

const (
	// DefaultDockerRegistry is the default docker registry address.
	DefaultDockerRegistry = "https://registry-1.docker.io"
)

// GetAuthConfig returns the docker registry AuthConfig.
func GetAuthConfig(username, password, registry string) (types.AuthConfig, error) {
	if username != "" && password != "" && registry != "" {
		return types.AuthConfig{
			Username:      username,
			Password:      password,
			ServerAddress: registry,
		}, nil
	}
	dcfg, err := config.Load(config.Dir())
	if err != nil {
		return types.AuthConfig{}, fmt.Errorf("loading config file failed: %v", err)
	}

	// return error early if there are no auths saved
	if !dcfg.ContainsAuth() {
		// If we were passed a registry, just use that.
		if registry != "" {
			return setDefaultRegistry(types.AuthConfig{
				ServerAddress: registry,
			}), nil
		}

		// Otherwise, just use an empty auth config.
		return types.AuthConfig{}, nil
	}

	authConfigs, err := dcfg.GetAllCredentials()
	if err != nil {
		return types.AuthConfig{}, fmt.Errorf("getting credentials failed: %v", err)
	}

	// if they passed a specific registry, return those creds _if_ they exist
	if registry != "" {
		if creds, ok := authConfigs[registry]; ok {
			fixAuthConfig(&creds, registry)
			return creds, nil
		}

		if strings.HasPrefix(registry, "https://") {
			registryCleaned := strings.TrimPrefix(registry, "https://")
			if creds, ok := authConfigs[registryCleaned]; ok {
				fixAuthConfig(&creds, registryCleaned)
				return creds, nil
			}
		}

		if strings.HasPrefix(registry, "http://") {
			registryCleaned := strings.TrimPrefix(registry, "http://")
			if creds, ok := authConfigs[registryCleaned]; ok {
				fixAuthConfig(&creds, registryCleaned)
				return creds, nil
			}
		}

		if !strings.HasPrefix(registry, "https://") && !strings.HasPrefix(registry, "http://") {
			registryCleaned := "https://" + registry
			if creds, ok := authConfigs[registryCleaned]; ok {
				fixAuthConfig(&creds, registryCleaned)
				return creds, nil
			}
		}

		// Otherwise just use the registry with no auth.
		return setDefaultRegistry(types.AuthConfig{
			ServerAddress: registry,
		}), nil
	}

	// Just set the auth config as the first registryURL, username and password
	// found in the auth config.
	for _, creds := range authConfigs {
		log.Logger.Debugf("No registry passed. Using registry %q\n", creds.ServerAddress)
		return creds, nil
	}

	log.Logger.Debug("Not using any authentication")
	return types.AuthConfig{}, nil
}

func fixAuthConfig(creds *types.AuthConfig, registry string) {
	if creds.ServerAddress == "" {
		creds.ServerAddress = registry
	}
}

func setDefaultRegistry(auth types.AuthConfig) types.AuthConfig {
	if auth.ServerAddress == "docker.io" {
		auth.ServerAddress = DefaultDockerRegistry
	}

	return auth
}
