package lambda

import "path"

/**
*@Author lyer
*@Date 4/8/21 11:18
*@Describe
**/

func cleanPath(fullPath string) string {
	if len(fullPath) == 0 || fullPath[0] != '/' {
		panic(InvalidFullPath)
	}
	return path.Clean(fullPath)
}
