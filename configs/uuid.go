package configs

import "math/rand"

type UID string

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Generate(identifier string, talename string) string {
	var length int = 10

	for {
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}

		uid := identifier + string(b)
		if uidExists(uid, talename) {
			return uid
		}
	}

}

func uidExists(uid string, tablename string) bool {
	var count int64
	err := DataBase.Table(tablename).Where("uid = ?", uid).Count(&count)

	if err.Error != nil {
		Logger.Error("error", err.Error)
		return false
	}

	Logger.Info(count)
	return count == 0
}
