package client

import "errors"

var errNilClientConnection = errors.New("nil grpc client connection provided")
