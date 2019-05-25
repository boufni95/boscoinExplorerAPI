package boscoin

import "strings"

func GetAccount(acc string) map[string]interface{} {
	var fullP strings.Builder
	fullP.WriteString("/api/v1/accounts/")
	fullP.WriteString(acc)
	fullP.WriteString("/frozen-accounts")
	bytes := GetBytes(fullP.String())
	ind := FmtIndent(bytes)
	SaveToFile("accountT.ggs", ind)
	return nil
}
