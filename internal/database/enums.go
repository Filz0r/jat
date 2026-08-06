package database

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

type JatMode string

const (
	StandaloneMode JatMode = "standalone"
	ServerMode     JatMode = "server"
	ClientMode     JatMode = "client"
)

func (m JatMode) Valid() bool {
	switch m {
	case StandaloneMode, ServerMode, ClientMode:
		return true
	}
	return false
}

func (m JatMode) String() string {
	return string(m)
}

type ApplicationStatusKind string

const (
	Applied     ApplicationStatusKind = "applied"
	Rejected                          = "rejected"
	Ghosted                           = "ghosted"
	Interviewed                       = "interviewed"
	Irrelevant                        = "irrelevant"
	Accepted                          = "accepted"
)

func (k ApplicationStatusKind) String() string {
	return string(k)
}

func (k ApplicationStatusKind) Valid() bool {
	switch k {
	case Applied, Rejected, Ghosted, Interviewed, Irrelevant, Accepted:
		return true
	default:
		return false
	}
}

type ConfigType string

const (
	StringConfig  ConfigType = "string"
	IntegerConfig            = "integer"
	IDConfig                 = "uuid"
	BooleanConfig            = "boolean"
	JatModeConfig            = "jat_mode"
)

func (t ConfigType) String() string {
	return string(t)
}

func (t ConfigType) Valid() bool {
	switch t {
	case StringConfig, IntegerConfig, IDConfig, BooleanConfig, JatModeConfig:
		return true
	default:
		return false
	}
}

func (t ConfigType) ToType(data string) (any, error) {
	switch t {
	case StringConfig:
		return data, nil
	case IntegerConfig:
		return strconv.Atoi(data)
	case IDConfig:
		return uuid.Parse(data)
	case BooleanConfig:
		return strconv.ParseBool(data)
	case JatModeConfig:
		mode := JatMode(data)
		if !mode.Valid() {
			return nil, errors.New("invalid jat_mode")
		}
		return mode, nil
	default:
		return nil, fmt.Errorf("unknown config type: %s", t)
	}
}

func (t ConfigType) ToString(data any) (string, error) {
	switch t {
	case StringConfig:
		return data.(string), nil
	case IntegerConfig:
		return strconv.Itoa(data.(int)), nil
	case IDConfig:
		return data.(uuid.UUID).String(), nil
	case BooleanConfig:
		return strconv.FormatBool(data.(bool)), nil
	case JatModeConfig:
		var mode JatMode
		switch v := data.(type) {
		case JatMode:
			mode = v
		case string:
			mode = JatMode(v)
		default:
			return "", fmt.Errorf("invalid jat mode value type")
		}
		if !mode.Valid() {
			return "", fmt.Errorf("invalid jat mode: %s", mode)
		}
		return mode.String(), nil
	default:
		return "", fmt.Errorf("unknown config type: %s", t)
	}
}
