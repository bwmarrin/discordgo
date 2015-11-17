// this file has small funcs used without the pacakge
// or, one.. util, maybe I'll have more later :)

package discordgo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// printJSON is a helper function to display JSON data in a easy to read format.
func printJSON(body []byte) {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		fmt.Print("JSON parse error: ", error)
	}
	fmt.Println(string(prettyJSON.Bytes()))
}
