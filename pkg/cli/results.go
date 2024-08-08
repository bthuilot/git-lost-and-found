package cli

import (
	"encoding/json"
	"io"
)

func WriteResults[T any](output io.Writer, results T) error {
	enc := json.NewEncoder(output)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
