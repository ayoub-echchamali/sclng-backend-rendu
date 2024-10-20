package util

func Substring(str string, start int, length int) string {
    return string([]rune(str)[start:length+start])
}