package common

import "github.com/bytom/vapor/errors"

var (
	ErrBadLenXPubStr      = errors.New("bad length of pubkey key string")
	ErrInvalidXPub        = errors.New("invalid xPub")
	ErrInvalidAccessToken = errors.New("invalid access token")
)
