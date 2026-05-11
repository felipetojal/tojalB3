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
	eng *engine.Engine

	dbDirPath  = "./badger-data"
	volumePath = "./volume.dat"

	rootCmd = &cobra.Command{
		Use:   "tojalB3",
		Short: "Local Block Storage Engine",
		Long:  "TojalB3 is a Block Storage Engine implementation",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			db, err := metadata.NewDatabase(dbDirPath)
			if err != nil {
				fmt.Printf("Error creating engine: %v", err)
				return fmt.Errorf("failed to create database: %w", err)
			}
			log.Println("Database successfully initialized!")

			v := volume.NewVolumeManager(volumePath)
			fmt.Println("Volume successfully loaded!")

			it, err := db.LoadIndexTable()
			if err != nil {
				fmt.Printf("Error creating engine: %v", err)
				return fmt.Errorf("failed to load index table: %w", err)
			}
			fmt.Println("Index Table loaded!")

			eng, err = engine.NewEngine(v, db, it)
			if err != nil {
				fmt.Printf("Error creating engine: %v", err)
				return fmt.Errorf("failed to create engine: %w", err)
			}
			fmt.Println("Engine created!")

			return nil
		},
	}

	storeCmd = &cobra.Command{
		Use:   "store",
		Short: "This command is used to store a file in the block storage engine.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Usage: store <file>")
				return
			}
			file := args[0]
			fmt.Println("Store command was called...")
			fmt.Printf("File: %v\n", file)

			if err := eng.StoreFile(file); err != nil {
				log.Fatal("Error: Unable to store file: %w", err)
			}

			fmt.Println("File was stored successfully!")
		},
	}

	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "This command is used to delete a file from the block storage engine.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Usage: delete <file>")
				return
			}
			file := args[0]
			fmt.Println("Delete command was called.")

			if err := eng.DeleteFile(file); err != nil {
				fmt.Println("Error: Unable to delete file: %w", err)
				return
			}
			fmt.Println("File was deleted successfully!")
		},
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "This command is used to retrieve a file from the block storage engine.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				fmt.Println("Usage: get <file> <output_path>")
				return
			}
			file := args[0]
			outputPath := args[1]
			fmt.Println("Get command was called.")

			if err := eng.GetFile(file, outputPath); err != nil {
				fmt.Println("Error: Unable to get file: %w", err)
				return
			}
			fmt.Println("File was retrieved successfully!")
		},
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "This command is used to list all the files stored in the block engine.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("List command was called.")

			files := eng.ListFiles()

			if len(files) == 0 {
				fmt.Println("No files found.")
				return
			}

			for _, file := range files {
				fmt.Println(file)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(storeCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute: %w", err)
	}
	return nil
}
