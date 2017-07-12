package server


func cmdLpush(c *client, args Args) error {
	return nil
}

func cmdRpush(c *client, args Args) error {
	return nil
}

func init() {
	register("lpush", cmdLpush)
	register("rpush", cmdRpush)
}