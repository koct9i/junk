package main

import (
	"context"
	"flag"
	"fmt"

	// "net"
	"os"
	"os/signal"

	// "path/filepath"
	// "strings"
	"time"

	// "golang.org/x/mod/semver"
	// "golang.org/x/sys/unix"

	"github.com/urfave/cli/v3"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"

	cri "k8s.io/cri-api/pkg/apis/runtime/v1"

	"gopkg.in/yaml.v3"
)

var (
	address = flag.String("address", "unix:///run/containerd/containerd.sock", "CRI endpoint")
)

func grpcLogger(logger logr.Logger) grpclogging.Logger {
	return grpclogging.LoggerFunc(func(ctx context.Context, level grpclogging.Level, msg string, fields ...any) {
		l := logger.WithValues(fields...)
		switch level {
		case grpclogging.LevelError:
			l.Error(nil, msg)
		case grpclogging.LevelDebug:
			l.V(2).Info(msg)
		default:
			l.Info(msg)
		}
	})
}

type Client interface {
	cri.RuntimeServiceClient
	cri.ImageServiceClient
}

func newClient(logger logr.Logger, address string) (Client, error) {
	glogLogger := grpcLogger(logger.WithName("grpc"))
	glogOpts := []grpclogging.Option{
		grpclogging.WithLogOnEvents(grpclogging.StartCall, grpclogging.FinishCall),
		grpclogging.WithLevels(grpclogging.DefaultClientCodeToLevel),
	}

	cc, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclogging.UnaryClientInterceptor(glogLogger, glogOpts...),
		),
		grpc.WithChainStreamInterceptor(
			grpclogging.StreamClientInterceptor(glogLogger, glogOpts...),
		),
	)
	if err != nil {
		logger.Error(err, "Cannot create gRPC client")
		return nil, err
	}

	return struct {
		cri.RuntimeServiceClient
		cri.ImageServiceClient
	}{
		cri.NewRuntimeServiceClient(cc),
		cri.NewImageServiceClient(cc),
	}, nil
}

func format(data any) string {
	text, err := yaml.Marshal(data)
	if err != nil {
		return "# " + err.Error()
	}
	return string(text)
}

func mainx() {
	flag.Parse()

	var logger logr.Logger
	{
		log := logrus.StandardLogger()
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.DateTime,
			FullTimestamp:   true,
		})
		log.SetLevel(logrus.Level(int(logrus.InfoLevel)))
		logger = logrusr.New(log)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client, err := newClient(logger, *address)
	if err != nil {
		os.Exit(1)
	}

	{
		rsp, _ := client.Status(ctx, &cri.StatusRequest{Verbose: true})
		fmt.Print("Status", "status", format(rsp))
	}

	{
		rsp, _ := client.ListPodSandbox(ctx, &cri.ListPodSandboxRequest{
			Filter: &cri.PodSandboxFilter{
				LabelSelector: map[string]string{
					"test.pod.namespace": "default",
					"test.pod.name":      "test-pod",
				},
			},
		})
		fmt.Print("Status", "status", format(rsp))
		for _, pod := range rsp.Items {
			rsp, err := client.RemovePodSandbox(ctx, &cri.RemovePodSandboxRequest{
				PodSandboxId: pod.Id,
			})
			fmt.Print("Result", err, format(rsp))
		}
	}

	podConfig := cri.PodSandboxConfig{
		Metadata: &cri.PodSandboxMetadata{
			Namespace: "default",
			Name:      "test-pod",
			Uid:       "test-pod-uid",
			Attempt:   1,
		},
		Labels: map[string]string{
			"test.pod.namespace": "default",
			"test.pod.name":      "test-pod",
		},
		Linux: &cri.LinuxPodSandboxConfig{
			SecurityContext: &cri.LinuxSandboxSecurityContext{
				Privileged: false,
				NamespaceOptions: &cri.NamespaceOption{
					Network: cri.NamespaceMode_NODE,
				},
			},
			CgroupParent: "/test.slice/test-pod.slice",
			Resources: &cri.LinuxContainerResources{
				MemoryLimitInBytes: 256 << 20,
				Unified: map[string]string{
					"memory.oom.group": "1",
				},
			},
		},
	}

	containerConfig := cri.ContainerConfig{
		Metadata: &cri.ContainerMetadata{
			Name:    "test",
			Attempt: 1,
		},
		Labels: map[string]string{
			"test.container.namespace": "default",
			"test.container.name":      "test",
		},
		Image: &cri.ImageSpec{
			Image: "docker.io/library/ubuntu:jammy",
		},
		Command:    []string{"sleep"},
		Args:       []string{"inf"},
		Envs:       []*cri.KeyValue{},
		WorkingDir: "/",
		Mounts: []*cri.Mount{
			&cri.Mount{
				ContainerPath: "/dev/shm",
				HostPath:      "/dev/shm/test",
				Propagation:   cri.MountPropagation_PROPAGATION_PRIVATE,
			},
		},
		Linux: &cri.LinuxContainerConfig{
			Resources: &cri.LinuxContainerResources{
				MemoryLimitInBytes: 128 << 20,
				Unified:            map[string]string{},
			},
			SecurityContext: &cri.LinuxContainerSecurityContext{
				Privileged: false,
				NamespaceOptions: &cri.NamespaceOption{
					Network: cri.NamespaceMode_POD,
					Pid:     cri.NamespaceMode_POD,
					Ipc:     cri.NamespaceMode_CONTAINER,
				},
				NoNewPrivs: true,
			},
		},
	}

	var podID string
	{
		rsp, err := client.RunPodSandbox(ctx, &cri.RunPodSandboxRequest{
			Config: &podConfig,
		})
		if err != nil {
			os.Exit(1)
		}
		fmt.Print("Result", format(rsp))
		podID = rsp.PodSandboxId
	}

	var containerID string
	{
		rsp, err := client.CreateContainer(ctx, &cri.CreateContainerRequest{
			PodSandboxId:  podID,
			Config:        &containerConfig,
			SandboxConfig: &podConfig,
		})
		if err != nil {
			os.Exit(1)
		}
		fmt.Print("Result", format(rsp))
		containerID = rsp.ContainerId
	}

	{
		rsp, err := client.StartContainer(ctx, &cri.StartContainerRequest{
			ContainerId: containerID,
		})
		if err != nil {
			os.Exit(1)
		}
		fmt.Print("Result", format(rsp))
	}
}

type ContainerCLI struct {
	Command *cli.Command

	Verbosity int

	RuntimeAddress string

	ImageAddress string
	Logger       logr.Logger
	Client       Client

	PodSandboxConfig	*cri.PodSandboxConfig
	ContianerConfig 	*cri.ContainerConfig
}

func (c *ContainerCLI) Connect() error {
	{
		log := logrus.StandardLogger()
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.DateTime,
			FullTimestamp:   true,
		})
		log.SetLevel(logrus.Level(int(logrus.InfoLevel) + c.Verbosity))
		c.Logger = logrusr.New(log)
	}

	var err error
	c.Client, err = newClient(c.Logger, c.RuntimeAddress)
	if err != nil {
		return err
	}
	return nil
}

func NewContainerCLI() *ContainerCLI {
	c := &ContainerCLI{}

	c.Command = &cli.Command{
		EnableShellCompletion: true,

		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "verbosity",
				Aliases: []string{"v"},
				Value:   0,
				Usage:   "logging verbosity `level`",
			},
			&cli.StringFlag{
				Name:        "address",
				Aliases:     []string{"a"},
				Usage:       "CRI endpoint",
				Value:       "unix:///run/containerd/containerd.sock",
				Sources:     cli.EnvVars("CONTAINER_RUNTIME_ENDPOINT", "CONTAINERD_ADDRESS"),
				Destination: &c.RuntimeAddress,
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			if err := c.Connect(); err != nil {
				return nil, err
			}
			return nil, nil
		},
		Commands: []*cli.Command{
			{
				Category: "service",
				Name:     "status",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					respose, err := c.Client.Status(ctx, &cri.StatusRequest{})
					if err != nil {
						return err
					}
					c.Logger.Info("status", "reponse", respose)
					return nil
				},
			},
		},
	}

	containerCommands := []*cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l", "ls"},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				respose, err := c.Client.ListContainers(ctx, &cri.ListContainersRequest{})
				if err != nil {
					return err
				}
				for _, ct := range respose.Containers {
					c.Logger.Info("ct", "id", ct.Id, "name", ct.Metadata.Name, "labels", ct.Labels)
				}
				return nil
			},
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "id",
					Min:  1,
					Max:  -1,
				},
			},
			ShellComplete: func(ctx context.Context, cmd *cli.Command) {
				if err := c.Connect(); err != nil {
					return
				}
				respose, err := c.Client.ListContainers(ctx, &cri.ListContainersRequest{})
				if err != nil {
					return
				}
				for _, ct := range respose.Containers {
					fmt.Println(ct.Id)
					fmt.Println(ct.Metadata.Name)
				}
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				for _, id := range cmd.Args().Slice() {
					respose, err := c.Client.ListContainers(ctx, &cri.ListContainersRequest{
						Filter: &cri.ContainerFilter{
							Id: id,
						},
					})
					if err != nil {
						return err
					}
					for _, ct := range respose.Containers {
						c.Logger.Info("ct", "id", ct.Id, "name", ct.Metadata.Name)
					}
				}
				return nil
			},
		},
		{
			Name:    "run",
			Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
				c.PodSandboxConfig = &cri.PodSandboxConfig{}
				c.ContianerConfig = &cri.ContainerConfig{}
				return nil, nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "name",
					Action: func(ctx context.Context, cmd *cli.Command, value string) error {
						c.ContianerConfig.Metadata.Name = value
						return nil
					},
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return nil
			},
		},
	}

	c.Command.Commands = append(c.Command.Commands, &cli.Command{
		Category: "container",
		Name:     "container",
		Aliases:  []string{"c", "ct"},
		Usage: "manage containers",
		Commands: containerCommands,
	})

	return c
}

func main() {
	containerCLI := NewContainerCLI()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := containerCLI.Command.Run(ctx, os.Args); err != nil {
		os.Exit(1)
	}
}
