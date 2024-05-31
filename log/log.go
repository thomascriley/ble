package log

import (
	"fmt"
	"log/slog"
)

func Stringer(key string, value fmt.Stringer) slog.Attr {
	return slog.String(key, value.String())
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func Bytes(key string, bs []byte) slog.Attr {
	return slog.String(key, fmt.Sprintf("%02X", bs))
}

func Uint8(key string, value uint8) slog.Attr {
	return slog.Int(key, int(value))
}

func Uint16(key string, value uint16) slog.Attr {
	return slog.Int(key, int(value))
}
