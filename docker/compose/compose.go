package compose

import (
	"context"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

type Compose struct {
	Project project.APIProject
}

func New(name, cpath string) (*Compose, error) {
	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{cpath},
			ProjectName:  name,
		},
	}, nil)

	if err != nil {
		return nil, err
	}

	return &Compose{Project: project}, nil
}

func (c *Compose) Up() error {
	return c.Project.Up(context.Background(), options.Up{})
}

func (c *Compose) Down(clean bool) error {
	var opts options.Down

	if clean {
		opts.RemoveImages = "local"
		opts.RemoveOrphans = true
		opts.RemoveVolume = true
	}

	return c.Project.Down(context.Background(), opts)
}

func (c *Compose) Start(services []string) error {
	return c.Project.Start(context.Background(), services...)
}

func (c *Compose) Stop(services []string) error {
	return c.Project.Stop(context.Background(), 60, services...)
}

func (c *Compose) Kill(services []string) error {
	return c.Project.Kill(context.Background(), "SIGKILL", services...)
}

func (c *Compose) Scale(services map[string]int) error {
	return c.Project.Scale(context.Background(), 120, services)
}

func (c *Compose) Logs(services []string, follow bool) error {
	return c.Project.Log(context.Background(), follow, services...)
}

func (c *Compose) Build(services []string, nocache bool) error {
	opts := options.Build{NoCache: nocache}

	for _, service := range services {
		err := c.Project.Build(context.Background(), opts, service)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compose) Rebuild(services []string) error {
	err := c.Stop(services)
	if err != nil {
		return err
	}

	err = c.Build(services, true)
	if err != nil {
		return err
	}

	return c.Start(services)
}