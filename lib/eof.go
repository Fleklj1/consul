package lib

import (
	"errors"
	"fmt"
	"io"
	"net/rpc"
	"strings"

	"github.com/hashicorp/yamux"
)

var yamuxStreamClosed = yamux.ErrStreamClosed.Error()
var yamuxSessionShutdown = yamux.ErrSessionShutdown.Error()

// IsErrEOF returns true if we get an EOF error from the socket itself, or
// an EOF equivalent error from yamux.
func IsErrEOF(err error) bool {
	if errors.Is(err, io.EOF) {
		return true
	}

	errStr := err.Error()
	if strings.Contains(errStr, yamuxStreamClosed) ||
		strings.Contains(errStr, yamuxSessionShutdown) {
		return true
	}

	if srvErr, ok := err.(rpc.ServerError); ok {
		return strings.HasSuffix(srvErr.Error(), fmt.Sprintf(": %s", io.EOF.Error()))
	}

	if srvErr, ok := errors.Unwrap(err).(rpc.ServerError); ok {
		return strings.HasSuffix(srvErr.Error(), fmt.Sprintf(": %s", io.EOF.Error()))
	}

	return false
}
