/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package log

type Level int

const (
	FATAL Level = 1
	ERROR Level = 2 | FATAL
	WARN  Level = 4 | ERROR
	INFO  Level = 8 | WARN
	DEBUG Level = 16 | INFO
	STACK Level = 32 | DEBUG
	ALL   Level = 64 | FATAL | ERROR | WARN | INFO | DEBUG | STACK
)

func (that Level) Is(level int) bool {
	return int(that) == (int(that) & level)
}

func (that Level) String() string {
	if that == FATAL {
		return "fatal"
	} else if that == ERROR {
		return "error"
	} else if that == WARN {
		return "warn"
	} else if that == INFO {
		return "info"
	} else if that == DEBUG {
		return "debug"
	} else if that == STACK {
		return "stack"
	}
	return "all"
}
