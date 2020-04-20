package tools

func GenerateUID(username string) string {
	uid := Sha1(username)
	return uid
}
