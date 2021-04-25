package tsdb

import (
	"os"

	"github.com/go-kit/kit/log"
)

var logger = log.NewLogfmtLogger(os.Stderr)
