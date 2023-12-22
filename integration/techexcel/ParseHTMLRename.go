package techexcel

import "strings"

func ParseHTMLError(htmltext string) string {
	returnVal := ""
	newstr := strings.Split(htmltext, "Messages:</strong>")
	if len(newstr) >= 2 {
		newstr = strings.Split(newstr[1], "</div>")
		if len(newstr) > 0 {
			// returnVal = common.RemoveSpace(newstr[0])
		}
	}
	return returnVal
}
