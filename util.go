package severe

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

//huh....
func copyLine(line []rune) []rune {
	line_ := make([]rune, len(line))
	copy(line_, line)
	return line_
}

func copyLines(lines [][]rune) [][]rune {
	lines_ := make([][]rune, len(lines))
	copy(lines_, lines)
	return lines_
}
