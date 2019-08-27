package discordgo

import (
	"sort"
	"strconv"
	"time"
)

// Contains checks if a slice of strings contains the string to search for
// haystack      : slice of strings to search in
// needle        : string to search for
func Contains(haystack []string, needle string) bool {
	sort.Strings(haystack)
	pos := sort.SearchStrings(haystack, needle)

	if pos == len(haystack) {
		return false
	}
	return haystack[pos] == needle
}

// ContainsIDObject checks if the haystack IDGettable contains the needle IDGettable
// haystack      : slice of IDGettables to search in
// needle        : IDGettable to search for
func ContainsIDObject(haystack []IDGettable, needle IDGettable) (contains bool) {
	if len(haystack) < 1 {
		return false
	}

	for _, item := range haystack {
		if item.GetID() == needle.GetID() {
			return true
		}
	}

	return false
}

// SnowflakeToTime converts a snowflake ID to a Time object
// snowflake      : the snowflake ID to convert
func SnowflakeToTime(snowflake string) (returnTime time.Time, err error) {
	n, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return
	}

	timestamp := ((n >> 22) + 1420070400000) * 1000000
	returnTime = time.Unix(timestamp, 0).UTC()
	return
}
