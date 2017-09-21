package storage

func GetVerifyCode(pkg, mobile string) string {
	key := pkg + ":vcode"
	vcd, err := rdsclt.HGet(key, mobile).Result()
	if err != nil {
		LogErr("1st Err from storage.vcode.GetVerifyCode(): ", err)
		return ""
	}
	go func() {
		if err := rdsclt.HDel(key, mobile).Err(); err != nil {
			LogErr("2nd Err from storage.vcode.GetVerifyCode(): ", err)
		}
	}()
	return vcd
}
