package enum

import "fmt"

type pathFormatGen int

// PathFormatGen DO NOT DEFINE ANOTHER pathFormatGen
const PathFormatGen pathFormatGen = 1

func (_ pathFormatGen) SelectPathFormat() string {
	return "%v" + fmt.Sprintf("%v", SELECT) + "%v"
}
func (_ pathFormatGen) InsertPathFormat() string {
	return "%v" + fmt.Sprintf("%v", INSERT) + "%v"
}
func (_ pathFormatGen) UpdatePathFormat() string {
	return "%v" + fmt.Sprintf("%v", UPDATE) + "%v"
}
func (_ pathFormatGen) DeletePathFormat() string {
	return "%v" + fmt.Sprintf("%v", DELETE) + "%v"
}
