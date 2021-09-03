package flags

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func FileArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("one argument, as the full path to a file, is required")
	}

	fileInfo, err := os.Stat(args[0])
	if err != nil {
		return fmt.Errorf("could not get info about file %q: %s", args[0], err)
	}

	if fileInfo.IsDir() {
		// Only directories that are proper bundles are acceptable
		bundlePath := filepath.Join(args[0], "Contents", "Info.plist")
		if _, err := os.Stat(bundlePath); err == nil {
			return nil // Valid bundle, so try to handle it
		}
		return fmt.Errorf("%q is a directory, should be a file or bundled directory", args[0])
	}

	return nil
}
