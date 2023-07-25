// Package main - メインパッケージ
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
	_ "github.com/goccy/go-yaml"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
	"golang.org/x/xerrors"
)

var (
	configPath  = flag.String("c", "config.yaml", "-c config.yaml")
	versionFlag = flag.Bool("v", false, "show version")
)

const (
	refIDsFileName = "ref_ids.yaml"
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

	// refIDの一覧
	refIDs := inserter.RefIDs{}
	refFilePath := func() string {
		if cfg.StateDir != "" {
			return filepath.Join(cfg.StateDir, refIDsFileName)
		}
		return ""
	}()

	// stateDirが指定されている場合は、そこからrefIDを読み込む
	if cfg.StateDir != "" {
		err := refIDs.LoadFromFile(refFilePath)
		if err != nil {
			log.Fatalf("failed to load state file: %+v", err)
		}
	}

	is := inserter.NewInserter(client, refIDs)
	refIDs, err = is.Execute(context.Background(), &cfg)
	if err != nil {
		log.Fatalf("failed to execute insert: %+v", err)
	}

	// stateDirが指定されている場合は、そこにrefIDを書き込む
	if cfg.StateDir != "" {
		err := refIDs.SaveToFile(refFilePath)
		if err != nil {
			log.Fatalf("failed to save state file: %+v", err)
		}
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
