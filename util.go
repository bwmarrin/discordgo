package discordgo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func printJSON(body []byte) {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		fmt.Print("JSON parse error: ", error)
	}
	fmt.Println("RESPONSE ::\n" + string(prettyJSON.Bytes()))
}
