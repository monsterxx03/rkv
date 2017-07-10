package server


func cmdLpush(c *client) ([]byte, error) {
	return []byte("+lpush\r\n"), nil
}

func cmdRpush(c *client) ([]byte, error) {
	return []byte("+rpush\r\n"), nil
}

func init() {
	register("lpush", cmdLpush)
	register("rpush", cmdRpush)
}