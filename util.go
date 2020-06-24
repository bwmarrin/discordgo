package discordgo

import (
	"strconv"
	"time"
)

// SnowflakeTimestamp returns the creation time of a Snowflake ID relative to the creation of Discord.
func SnowflakeTimestamp(ID string) (t time.Time, err error) {
	i, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return
	}
	timestamp := (i >> 22) + 1420070400000
	t = time.Unix(0, timestamp*1000000)
	return
}

//PermissionsSynced Checks if channel has permissions synced with another channel
func PermissionsSynced(c1 *Channel, c2 *Channel) bool {
	for index, po := range c1.PermissionOverwrites {
		if po != c2.PermissionOverwrites[index] {
			return false
		}
	}
	return true
}
