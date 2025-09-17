// 代码生成时间: 2025-09-17 19:31:08
package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// Config represents the configuration for the backup and sync tool
type Config struct {
    Source   string        "json:"source""
    Destination string      "json:"destination""
    Excludes  []string      "json:"excludes""
    Interval  time.Duration `json:"interval"`
}

// BackupSync is the main structure that holds configuration and performs backup and sync
type BackupSync struct {
    config Config
}

// NewBackupSync creates a new instance of BackupSync with the given config
func NewBackupSync(config Config) *BackupSync {
    return &BackupSync{config: config}
}

// Sync performs the file backup and synchronization
func (b *BackupSync) Sync() error {
    // Check if source and destination directories exist
    info, err := os.Stat(b.config.Source)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("source directory does not exist: %s", b.config.Source)
        }
        return fmt.Errorf("error checking source directory: %w", err)
    }
    if !info.IsDir() {
        return fmt.Errorf("source is not a directory: %s", b.config.Source)
    }

    info, err = os.Stat(b.config.Destination)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("destination directory does not exist: %s", b.config.Destination)
        }
        return fmt.Errorf("error checking destination directory: %w", err)
    }
    if !info.IsDir() {
        return fmt.Errorf("destination is not a directory: %s", b.config.Destination)
    }

    // Walk through the source directory
    err = filepath.WalkDir(b.config.Source, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }
        // Skip directories and files that are in the excludes list
        if d.IsDir() {
            for _, exclude := range b.config.Excludes {
                if strings.HasPrefix(path, filepath.Join(b.config.Source, exclude)) {
                    return filepath.SkipDir
                }
            }
        } else {
            for _, exclude := range b.config.Excludes {
                if strings.HasPrefix(path, filepath.Join(b.config.Source, exclude)) {
                    return nil
                }
            }
            // Copy file to destination
            return b.copyFile(path, b.config.Destination)
        }
        return nil
    })
    if err != nil {
        return fmt.Errorf("error during sync: %w", err)
    }
    return nil
}

// copyFile copies a single file from source to destination
func (b *BackupSync) copyFile(sourcePath, destinationPath string) error {
    srcFile, err := os.Open(sourcePath)
    if err != nil {
        return fmt.Errorf("error opening source file: %w", err)
    }
    defer srcFile.Close()

    destFile, err := os.OpenFile(filepath.Join(destinationPath, filepath.Base(sourcePath)), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        return fmt.Errorf("error creating destination file: %w", err)
    }
    defer destFile.Close()

    _, err = bufio.NewReader(srcFile).WriteTo(destFile)
    if err != nil {
        return fmt.Errorf("error writing to destination file: %w", err)
    }
    return nil
}

func main() {
    // Define the configuration
    config := Config{
        Source:   "/path/to/source",
        Destination: "/path/to/destination",
        Excludes:  []string{
            ".git",
            "node_modules",
            "tmp",
        },
        Interval: 10 * time.Minute,
    }

    // Create a new BackupSync instance
    sync := NewBackupSync(config)

    // Perform the initial sync
    if err := sync.Sync(); err != nil {
        log.Fatalf("initial sync failed: %v", err)
    }

    // Continue syncing at the specified interval
    for range time.NewTicker(config.Interval).C {
        if err := sync.Sync(); err != nil {
            log.Printf("sync failed: %v", err)
        }
    }
}
