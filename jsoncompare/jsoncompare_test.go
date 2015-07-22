package jsoncompare

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestJsonCompare(t *testing.T) {
	TestingT(t)
}

const (
	// 10000 microsecons = 10 milliseconds
	TTL       int = 10000
	ChLen     int = 10
	ThreadLen int = 5
)

var JSON []byte = []byte(`{"family":{"papa":"yes","mama":"too","children":{"count":2,"list":[{"type":"boy","name":"John"},{"type":"girl","name":"Linda"}]}}}`)
var JSON3 []byte = []byte(`{"family":{"papa":"yes","mama":"too","children":{"number":3,"list":[{"type":"boy","name":"John"},{"type":"girl","name":"Linda"},{"type":"boy","name":"Mike"}]}}}`)

type JsonCompareTestsSuite struct{}

var _ = Suite(&JsonCompareTestsSuite{})

func (s *JsonCompareTestsSuite) TestJsonCompareGetEqual(c *C) {
	//c.Skip("Not now")
	p1 := []path{}
	p1 = append(p1, path{path: "1", mytype: "map", value: 11})
	p1 = append(p1, path{path: "2", mytype: "slice", value: 22})
	p1 = append(p1, path{path: "3", mytype: "int64", value: 45})
	p1 = append(p1, path{path: "4", mytype: "string", value: "www"})
	p1 = append(p1, path{path: "5", mytype: "<nil>", value: nil})

	for i, v := range p1 {
		for j, w := range p1 {
			eq := getEqual(v, w)
			if i == j {
				c.Assert(eq, Equals, true)
			} else {
				c.Assert(eq, Equals, false)
			}
		}
	}

}

func (s *JsonCompareTestsSuite) TestJsonCompareGetJson(c *C) {
	//c.Skip("Not now")
	obj, err := getJson(JSON)
	c.Assert(err, IsNil)
	c.Assert(obj, NotNil)
}

func (s *JsonCompareTestsSuite) TestJsonCompareGetJsonEmpty(c *C) {
	//c.Skip("Not now")
	json := []byte(`{}`)
	obj, err := getJson(json)
	c.Assert(err, IsNil)
	c.Assert(obj, NotNil)
}

func (s *JsonCompareTestsSuite) TestJsonCompareAllPaths(c *C) {
	//c.Skip("Not now")
	obj, err := getJson(JSON)
	c.Assert(err, IsNil)
	c.Assert(obj, NotNil)
	paths := allPaths(obj, "")
	c.Assert(paths, NotNil)
	c.Assert(len(paths), Equals, 12)
}

func (s *JsonCompareTestsSuite) TestJsonCompareAllPathsEmpty(c *C) {
	//c.Skip("Not now")
	obj, err := getJson([]byte(`{}`))
	c.Assert(err, IsNil)
	c.Assert(obj, NotNil)
	paths := allPaths(obj, "")
	c.Assert(paths, NotNil)
	c.Assert(len(paths), Equals, 0)
}

func (s *JsonCompareTestsSuite) TestJsonCompareComparePaths(c *C) {

	//c.Skip("Not now")
	obj_l, err_l := getJson(JSON)
	c.Assert(err_l, IsNil)
	c.Assert(obj_l, NotNil)

	obj_r, err_r := getJson(JSON3)
	c.Assert(err_r, IsNil)
	c.Assert(obj_r, NotNil)

	pathsLeft := allPaths(obj_l, "")
	pathsRigth := allPaths(obj_r, "")

	checkList := comparePaths(pathsLeft, pathsRigth)
	c.Assert(len(checkList), Equals, 16)
}

func (s *JsonCompareTestsSuite) TestJsonCompareCompare(c *C) {
	//c.Skip("Not now")
	checkList, err := Compare(JSON, JSON3)
	c.Assert(err, IsNil)
	c.Assert(len(checkList), Equals, 16)
}

func (s *JsonCompareTestsSuite) TestJsonCompareSplitBySide(c *C) {

	//c.Skip("Not now")

	checkList, err := Compare(JSON, JSON3)
	c.Assert(err, IsNil)

	leftOnly, rightOnly, noEqual, goodList := SplitBySide(checkList)

	c.Assert(len(leftOnly), Equals, 1)
	c.Assert(len(rightOnly), Equals, 4)
	c.Assert(len(noEqual), Equals, 1)
	c.Assert(len(goodList), Equals, 10)
}
