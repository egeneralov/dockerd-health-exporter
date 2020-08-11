package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricsListen = "0.0.0.0:8909"
	namespace     = "dockerd"
	cli           *client.Client
	err           error
	registry      = prometheus.NewRegistry()
)

func main() {
	flag.StringVar(&metricsListen, "listen-addr", "0.0.0.0:8909", "bind http server")
	flag.StringVar(&namespace, "metrics-namespace", "dockerd", "metric namespace")
	flag.Parse()

	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			containers, err := getContainers()
			if err != nil {
				panic(err)
				// continue
			}
			for _, container := range containers {
				cname := strings.Split(container.Names[0], "/")
				if err := registry.Register(prometheus.NewGaugeFunc(
					prometheus.GaugeOpts{
						Namespace: namespace,
						Name:      "container_health_" + cname[1],
						Help:      "1 = healthy",
						ConstLabels: prometheus.Labels{
							"id":      container.ID,
							"image":   container.Image,
							"imageid": container.ImageID,
						},
					},
					func() float64 {
						container, err := getContainerByName(cname[1])
						if err != nil {
							return float64(0)
						}
						matched, err := regexp.MatchString(`^Up .* \(healthy\)$`, container.Status)
						if err == nil {
							if matched {
								return float64(1)
							} else {
  							return float64(0)
							}
						} else {
							return float64(0)
						}
						return float64(0)
					},
				)); err != nil {
					continue
				}
			}
			time.Sleep(time.Second)
		}
	}()

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	fmt.Println(http.ListenAndServe(metricsListen, nil))
}

func getContainers() ([]types.Container, error) {
	return cli.ContainerList(context.Background(), types.ContainerListOptions{})
}

func getContainerByName(name string) (types.Container, error) {
	containers, err := getContainers()
	if err != nil {
		return types.Container{}, err
	}
	for _, c := range containers {
		if c.Names[0] == "/"+name {
			return c, nil
		}
	}
	return types.Container{}, errors.New("404 not found")
}
