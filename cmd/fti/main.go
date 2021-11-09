package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"cloud.google.com/go/firestore"
	_ "github.com/BurntSushi/toml"
	"github.com/gcp-kit/fti/common"
	"github.com/gcp-kit/fti/pkg/config"
	"github.com/gcp-kit/fti/pkg/files"
	"github.com/gcp-kit/fti/pkg/inserter"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
	"golang.org/x/xerrors"
	_ "gopkg.in/yaml.v2"
)

var (
	configPath  = flag.String("c", "config.yaml", "-c config.yaml")
	versionFlag = flag.Bool("v", false, "show version")
)

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Println(common.AppVersion)
		return
	}

	if !files.Exists(*configPath) {
		configFullPath, _ := filepath.Abs(*configPath)
		log.Fatalf("not found configuration file: %s", configFullPath)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg := config.Config{}
	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend(*configPath),
		flags.NewBackend(),
	)
	err := loader.Load(ctx, &cfg)
	if err != nil {
		log.Fatalf("cannot load configuration: %+v", err)
	}

	client, err := initFirestore(ctx, &cfg)
	if err != nil {
		log.Fatalf("failed to init firestore: %+v", err)
	}

	is := inserter.NewInserter(client)
	err = is.Execute(context.Background(), &cfg)
	if err != nil {
		log.Fatalf("failed to execute insert: %+v", err)
	}
}

func initFirestore(ctx context.Context, cfg *config.Config) (*firestore.Client, error) {
	projectID := config.GetProjectID(cfg)
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize firestore client: %w", err)
	}

	return client, nil
}
