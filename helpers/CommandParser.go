package helpers

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/trembachLeonid/lest-memory-storage/logging"
)

type CommandParser struct {
	reader *bufio.Reader
}

func NewCommandParser(reader *bufio.Reader) *CommandParser {
	return &CommandParser{reader: reader}
}

func (cp *CommandParser) Parse(ctx *context.Context) ([][]byte, error) {
	logger := logging.FromContext(*ctx)

	message, err := cp.reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		logger.Error("Error reading command", "error", err)
		return nil, err
	}

	var entryType = message[0]
	if entryType == '+' && string(message[1:2]) == "QUIT" {
		return [][]byte{QUIT}, nil
	}
	if entryType != '*' {
		logger.Error("Unsupported entry type", "type", entryType)
		return nil, errors.New("unsupported entry type")
	}

	paramCount, err := strconv.Atoi(string(message[1:2])) // What if param count >= 10?
	if err != nil {
		logger.Error("Error parsing param count", "error", err)
	}

	params := make([][]byte, paramCount)

	for i := 0; i < paramCount; i++ {
		line, err := cp.reader.ReadString('\n')

		if err != nil && err != io.EOF {
			logger.Error("Error reading parameter", "error", err)
			return nil, err
		}

		if line[0] == '$' {
			paramLength, err := strconv.Atoi(strings.TrimSpace(line[1:]))

			if err != nil {
				logger.Error("Error parsing param length", "error", err)
				return nil, err
			}

			param := make([]byte, paramLength)
			_, err = io.ReadFull(cp.reader, param)

			params[i] = param

			crlf := make([]byte, 2)
			_, err = io.ReadFull(cp.reader, crlf)
			if !bytes.Equal(crlf, CRLF) {
				logger.Error("Error CRLF must be read after each parameter", "error", err)
				return nil, err
			}
		} else {
			logger.Error("Unsupported param type", "type", line[0])
		}

	}

	logger.Info("COMMAND PARSED", "params", params)

	return params, nil
}
