package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var simpleRegex, keyRegex, setRegex, setExRegex, zAddRegex, zRankRegex, zRangeRegex *regexp.Regexp

func init() {
	simpleRegex = regexp.MustCompile("^DBSIZE$")
	keyRegex = regexp.MustCompile("^(?P<cmd>GET|DEL|INCR|ZCARD) (?P<key>[a-zA-Z0-9-_]+)$")
	setRegex = regexp.MustCompile("^SET (?P<key>[a-zA-Z0-9-_]+) (?P<value>[a-zA-Z0-9-_]+)$")
	setExRegex = regexp.MustCompile("^SET (?P<key>[a-zA-Z0-9-_]+) (?P<value>[a-zA-Z0-9-_]+) EX (?P<seconds>[0-9]+)$")
	zAddRegex = regexp.MustCompile("^ZADD (?P<key>[a-zA-Z0-9-_]+) (?P<score>[0-9]+) (?P<member>[a-zA-Z0-9-_]+)$")
	zRankRegex = regexp.MustCompile("^ZRANK (?P<key>[a-zA-Z0-9-_]+) (?P<member>[a-zA-Z0-9-_]+)$")
	zRangeRegex = regexp.MustCompile("^ZRANGE (?P<key>[a-zA-Z0-9-_]+) (?P<start>[0-9]+) (?P<stop>[0-9]+)$")
}

type Interpreter struct {
	*Store
}

func (intr Interpreter) Exec(cmd string) (interface{}, error) {
	switch {
	case simpleRegex.MatchString(cmd):
		return intr.handleSimpleRegex(cmd)
	case keyRegex.MatchString(cmd):
		return intr.handleKeyRegex(cmd)
	case setRegex.MatchString(cmd):
		return intr.handleSetRegex(cmd)
	case setExRegex.MatchString(cmd):
		return intr.handleSetExRegex(cmd)
	case zAddRegex.MatchString(cmd):
		return intr.handleZAddRegex(cmd)
	case zRankRegex.MatchString(cmd):
		return intr.handleZRankRegex(cmd)
	case zRangeRegex.MatchString(cmd):
		return intr.handleZRangeRegex(cmd)
	}

	return errorReturn(cmd)
}

func (intr Interpreter) handleSimpleRegex(cmd string) (interface{}, error) {
	switch cmd {
	case "DBSIZE":
		return intr.DbSize(), nil
	}

	return errorReturn(cmd)
}

func (intr Interpreter) handleKeyRegex(str string) (interface{}, error) {
	values := scanVars(keyRegex, str, "cmd", "key")
	cmd, key := values[0], values[1]

	switch cmd {
	case "GET":
		if v, ok, err := intr.Get(key); err == nil {
			if ok {
				return v, nil
			}
		} else {
			return nil, err
		}

		return nil, nil
	case "DEL":
		return intr.Del(key), nil
	case "INCR":
		return intr.Incr(key)
	case "ZCARD":
		return intr.ZCard(key)
	}

	return errorReturn(cmd)
}

func (intr Interpreter) handleSetRegex(str string) (interface{}, error) {
	values := scanVars(setRegex, str, "key", "value")
	key, value := values[0], values[1]

	return intr.Set(key, value)
}

func (intr Interpreter) handleSetExRegex(str string) (interface{}, error) {
	values := scanVars(setExRegex, str, "key", "value", "seconds")
	key, value, secondsStr := values[0], values[1], values[2]

	seconds, _ := strconv.Atoi(secondsStr)
	return intr.SetEx(key, value, seconds)
}

func (intr Interpreter) handleZAddRegex(str string) (interface{}, error) {
	values := scanVars(zAddRegex, str, "key", "score", "member")
	key, scoreStr, member := values[0], values[1], values[2]

	score, _ := strconv.Atoi(scoreStr)
	item := SortedSetItem{float64(score), member}

	return intr.ZAdd(key, item)
}

func (intr Interpreter) handleZRankRegex(str string) (interface{}, error) {
	values := scanVars(zRankRegex, str, "key", "member")
	key, member := values[0], values[1]

	index, ok, err := intr.ZRank(key, member)

	switch {
	case err != nil:
		return nil, err
	case ok:
		return index, nil
	default:
		return nil, nil
	}
}

func (intr Interpreter) handleZRangeRegex(str string) (interface{}, error) {
	values := scanVars(zRangeRegex, str, "key", "start", "stop")
	key, startStr, stopStr := values[0], values[1], values[2]

	start, _ := strconv.Atoi(startStr)
	stop, _ := strconv.Atoi(stopStr)

	if items, err := intr.ZRange(key, start, stop); err == nil {
		members := make([]string, len(items))
		for index, item := range items {
			members[index] = item.Member
		}

		return members, nil
	} else {
		return nil, err
	}
}

func scanVars(regex *regexp.Regexp, str string, keys ...string) []string {
	groupNames := regex.SubexpNames()
	matches := regex.FindStringSubmatch(str)

	values := make([]string, 0)
	for index, value := range matches {
		for _, key := range keys {
			if groupNames[index] == key {
				values = append(values, value)
			}
		}
	}

	return values
}

func errorReturn(cmd string) (interface{}, error) {
	return nil, fmt.Errorf("miniredis: invalid command %q", cmd)
}
