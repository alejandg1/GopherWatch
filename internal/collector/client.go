package collector

import (
	"fmt"
	"os"

	"github.com/docker/docker/client"
)

func NewDockerClient() (*client.Client, error) {
	host := getDockerHost()
	opts := []client.Opt{
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	}

	if host != "" {
		opts = append(opts, client.WithHost(host))
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func getDockerHost() string {
	// 1. Si ya existe variable de entorno, dejamos que el SDK la use (retornando vacío aquí o chequeándola)
	if os.Getenv("DOCKER_HOST") != "" {
		return ""
	}

	// 2. Chequear sockets comunes
	candidates := []string{
		"unix:///var/run/docker.sock",
		"unix:///run/podman/podman.sock",
	}

	// 3. Chequear rootless podman (XDG_RUNTIME_DIR)
	if xdg := os.Getenv("XDG_RUNTIME_DIR"); xdg != "" {
		candidates = append(candidates, fmt.Sprintf("unix://%s/podman/podman.sock", xdg))
		candidates = append(candidates, fmt.Sprintf("unix://%s/docker.sock", xdg))
	}

	// 4. Fallback manual con UID si XDG no está (común en servidores/ssh)
	uid := os.Getuid()
	candidates = append(candidates, fmt.Sprintf("unix:///run/user/%d/podman/podman.sock", uid))

	for _, path := range candidates {
		// Quitamos el prefijo unix:// para chequear existencia de archivo
		cleanPath := path[7:]
		if _, err := os.Stat(cleanPath); err == nil {
			return path
		}
	}

	return ""
}
