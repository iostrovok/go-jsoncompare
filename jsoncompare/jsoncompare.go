package jsoncompare

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type PathDiff struct {
	PathLeft  string
	PathRight string
	IsEqual   bool
}

type path struct {
	path   string      // from top to elemnt likes: "/user/email/0/login"
	mytype string      // "map", "slice", "value"
	value  interface{} // Vvalue of element. There is len(value) if mytype is maps or slices.
}

func SplitBySide(list []*PathDiff) (leftOnly, rightOnly, noEqual, goodList []*PathDiff) {
	for _, v := range list {
		if v.PathLeft == "" {
			rightOnly = append(rightOnly, v)
		} else if v.PathRight == "" {
			leftOnly = append(leftOnly, v)
		} else if !v.IsEqual {
			noEqual = append(noEqual, v)
		} else {
			goodList = append(goodList, v)
		}
	}

	return
}

// Compare returns list of diffs
func Compare(left, rigth []byte) ([]*PathDiff, error) {

	var bodyLeft, bodyRigth map[string]interface{}
	var err error

	if bodyLeft, err = getJson(left); err != nil {
		return nil, err
	}
	if bodyRigth, err = getJson(rigth); err != nil {
		return nil, err
	}

	pathsLeft := allPaths(bodyLeft, "")
	pathsRigth := allPaths(bodyRigth, "")

	checkList := comparePaths(pathsLeft, pathsRigth)

	return checkList, nil
}

func getJson(data []byte) (map[string]interface{}, error) {

	list := []byte(`{"top":`)
	list = append(list, data...)
	list = append(list, []byte(`}`)...)

	var out map[string]interface{}
	err := json.Unmarshal(list, &out)
	return out, err
}

func allPaths(body interface{}, way string) []path {

	list := []path{}

	switch body.(type) {
	case map[string]interface{}:
		list = mAppend(list, way, "map", len(body.(map[string]interface{})))
		for k, v := range body.(map[string]interface{}) {
			next := way + "/" + k
			p := allPaths(v, next)
			list = append(list, p...)
		}
	case []interface{}:
		list = mAppend(list, way, "slice", len(body.([]interface{})))
		for k, v := range body.([]interface{}) {
			next := way + "/" + strconv.Itoa(k)
			p := allPaths(v, next)
			list = append(list, p...)
		}
	default:
		list = mAppend(list, way, fmt.Sprintf("%T", body), body)
	}

	out := []path{}
	for _, v := range list {
		v.path = strings.TrimPrefix(v.path, "/top")
		if v.path != "" {
			out = append(out, v)
		}
	}

	return out
}

func mAppend(list []path, way, mytype string, body interface{}) []path {

	if way == "" {
		return list
	}

	p := path{
		path:   way,
		mytype: mytype,
		value:  body,
	}
	list = append(list, p)
	return list
}

func comparePaths(left, rigth []path) []*PathDiff {

	out := []*PathDiff{}
	rightCheck := map[string]bool{}

	for _, vL := range left {
		diff := fundInPath(true, vL, rigth)
		if diff.PathRight != "" {
			rightCheck[diff.PathRight] = true
		}
		out = append(out, diff)
	}

	for _, vR := range rigth {
		if !rightCheck[vR.path] {
			diff := fundInPath(false, vR, left)
			out = append(out, diff)
		}
	}

	return out
}

func fundInPath(leftOrRight bool, from path, p []path) *PathDiff {

	diff := &PathDiff{
		PathLeft:  "",
		PathRight: "",
		IsEqual:   false,
	}

	if leftOrRight {
		diff.PathLeft = from.path
	} else {
		diff.PathRight = from.path
	}

	for _, to := range p {
		if from.path == to.path {

			if diff.PathLeft == "" {
				diff.PathLeft = to.path
			} else {
				diff.PathRight = to.path
			}

			diff.IsEqual = getEqual(from, to)
			break
		}
	}

	return diff
}

func getEqual(vL, vR path) bool {

	if vL.mytype != vR.mytype {
		return false
	}

	return vL.value == vR.value
}
