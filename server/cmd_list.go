package server


func cmdLpush(c *client) error {
	return nil
}

func cmdRpush(c *client) error {
	return nil
}

func init() {
	register("lpush", cmdLpush)
	register("rpush", cmdRpush)
}