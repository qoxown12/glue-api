package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"sort"
	"strings"
)

type SortMoldParams []MoldParams
type MoldParams map[string]string

func (s SortMoldParams) Len() int {
	return len(s)
}
func (s SortMoldParams) Less(i, j int) bool {
	for keyi := range s[i] {
		for keyj := range s[j] {
			return keyi < keyj
		}
	}
	return false
}

func (s SortMoldParams) Swap(i, j int) {
	for keyi, valuei := range s[i] {
		for keyj, valuej := range s[j] {
			s[i][keyj] = valuej
			s[j][keyi] = valuei
			delete(s[i], keyi)
			delete(s[j], keyj)
			return
		}
	}
}
func makeStringParams(params []MoldParams) string {
	var result string
	mold, _ := ReadMoldFile()

	params1 := []MoldParams{
		{"apikey": mold.MoldApiKey},
		{"response": "json"},
	}
	params = append(params, params1...)
	sort.Sort(SortMoldParams(params))

	for _, tuple := range params {
		for key, value := range tuple {
			result = result + key + "=" + value + "&"
		}
	}
	result = strings.TrimRight(result, "&")
	return result
}

func makeSignature(payload string) string {
	mold, _ := ReadMoldFile()
	secretkey := mold.MoldSecretKey
	strurl := strings.Replace(strings.ToLower(payload), "+", "%20", -1)
	secret := []byte(secretkey)
	message := []byte(strurl)
	hash := hmac.New(sha1.New, secret)
	hash.Write(message)
	strHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	log.Fatal("makeSignature payload [%v]\n", payload)
	println(strHash)
	returnString := strings.Replace(strHash, "+", "%2B", -1)
	println(returnString)
	return returnString
}
