package kstar

import "testing"

func TestRemoveLoopPaths(t *testing.T) {
	e1 := &Edge{U: 0, V: 1, I: 0}
	e2 := &Edge{U: 1, V: 2, I: 0}
	e3 := &Edge{U: 1, V: 0, I: 0}
	pathWithoutLoops := []*Edge{
		e1,
		e2,
	}
	pathWithLoops := []*Edge{
		e1,
		e3,
	}
	paths := [][]*Edge{
		pathWithLoops,
		pathWithoutLoops,
	}
	modifiedPaths := RemoveLoopPaths(paths)
	expectedPaths := [][]*Edge{pathWithoutLoops}
	if len(modifiedPaths) != 1 || !modifiedPaths[0][0].equals(e1) || !modifiedPaths[0][1].equals(e2) {
		t.Errorf("RemoveLoopPaths failed.\nExpected\n%v\ngot\n%v", expectedPaths, modifiedPaths)
	}
}
