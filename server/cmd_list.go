package server


func cmdLpush(c *client, args [][]byte) error {
	return nil
}

func cmdRpush(c *client, args [][]byte) error {
	return nil
}

func init() {
	register("lpush", cmdLpush)
	register("rpush", cmdRpush)
}