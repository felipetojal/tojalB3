package cmd

import (
	"fmt"

	"github.com/felipetojal/tojalB3/internal/engine"
	"github.com/felipetojal/tojalB3/internal/metadata"
	"github.com/felipetojal/tojalB3/internal/volume"
	"github.com/spf13/cobra"
)

var (
	eng            *engine.Engine
	destFilePath   string
	originFilePath string

	dbDirPath  = "./badger-data"
	volumePath = "./volume.dat"

	rootCmd = &cobra.Command{
		Use:   "TojalB3",
		Short: "Local Block Storage Engine",
		Long:  "TojalB3 is a Block Storage Engine implementation",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			db, err := metadata.NewDatabase(dbDirPath)
			if err != nil {
				return fmt.Errorf("failed to create database: %w", err)
			}

			v := volume.NewVolumeManager(volumePath)

			it, err := db.LoadIndexTable()
			if err != nil {
				return fmt.Errorf("failed to load index table: %w", err)
			}

			eng, err = engine.NewEngine(v, db, it)
			if err != nil {
				return fmt.Errorf("failed to create engine: %w", err)
			}

			return nil
		},
	}
)

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute: %w", err)
	}
	return nil
}
