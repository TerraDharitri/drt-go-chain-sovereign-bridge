package server

import "errors"

var errNilTxSender = errors.New("nil tx sender provided")

var errNilMarshaller = errors.New("nil marshaller provided")

var errNilGinHandler = errors.New("nil gin handler provided")

var errNilGRPCHandler = errors.New("nil grpc handler provided")
