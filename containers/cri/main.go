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
	verbosity = flag.Int("verbosity", 0, "logging verbosity level")
	address   = flag.String("address", "unix:///run/containerd/containerd.sock", "CRI endpoint")
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

func main() {
	flag.Parse()

	var logger logr.Logger
	{
		log := logrus.StandardLogger()
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.DateTime,
			FullTimestamp:   true,
		})
		log.SetLevel(logrus.Level(int(logrus.InfoLevel) + *verbosity))
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
				Propagation: cri.MountPropagation_PROPAGATION_PRIVATE,
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
