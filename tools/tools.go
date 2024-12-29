//go:build tools

package main

import (
	_ "github.com/bwplotka/bingo"

	// _ "github.com/elastic/crd-ref-docs"

	_ "github.com/fullstorydev/grpcurl/cmd/grpcurl"

	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"

	_ "github.com/mikefarah/yq/v4"

	_ "github.com/onsi/ginkgo/v2/ginkgo"

	_ "github.com/posener/complete/v2/gocomplete"

	_ "github.com/vadimi/grpc-client-cli/cmd/grpc-client-cli"

	// _ "github.com/dave/rebecca/cmd/becca"

	_ "golang.org/x/tools/cmd/gonew"

	// _ "golang.org/x/tools/gopls"

	_ "oras.land/oras/cmd/oras"

	_ "sigs.k8s.io/controller-runtime/tools/setup-envtest"

	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"

	_ "sigs.k8s.io/cri-tools/cmd/crictl"

	_ "sigs.k8s.io/kind"

	_ "sigs.k8s.io/kustomize/kustomize/v5"
)
