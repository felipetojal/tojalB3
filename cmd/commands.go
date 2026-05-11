package cmd

import (
	"fmt"
	"log"

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
				log.Printf("Error creating engine: %v", err)
				return fmt.Errorf("failed to create database: %w", err)
			}
			log.Println("Database successfully initialized!")

			v := volume.NewVolumeManager(volumePath)
			log.Println("Volume successfully loaded!")

			it, err := db.LoadIndexTable()
			if err != nil {
				log.Printf("Error creating engine: %v", err)
				return fmt.Errorf("failed to load index table: %w", err)
			}
			log.Println("Index Table loaded!")

			eng, err = engine.NewEngine(v, db, it)
			if err != nil {
				log.Printf("Error creating engine: %v", err)
				return fmt.Errorf("failed to create engine: %w", err)
			}
			log.Println("Engine created!")

			return nil
		},
	}

	storeCmd = &cobra.Command{
		Use:   "Store command",
		Short: "This command is used to store a file in the block storage engine.",
		Run: func(cmd *cobra.Command, args []string) {
			file := args[0]
			log.Println("Store command was called...")
			log.Printf("File: %v\n", file)

			if err := eng.StoreFile(file); err != nil {
				log.Fatal("Error: Unable to store file: %w", err)
			}

			log.Println("File was stored successfully!")
		},
	}
)

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute: %w", err)
	}
	return nil
}
