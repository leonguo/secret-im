package models

import (
	"encoding/json"
	redisClient "../../db/redis"
	pgorm "../../db/gorm"
)

type Keys struct {
	Id         int64  `gorm:"primary_key" json:"id"`
	Number     string `gorm:"column:number" json:"number"`
	KeyId      int64  `gorm:"column:key_id" json:"key_id"`
	PublicKey  string `gorm:"column:public_key" json:"public_key"`
	LastResort int    `gorm:"column:last_resort" json:"last_resort"`
	DeviceId   int64  `gorm:"column:device_id" json:"device_id"`
}

func (Keys) TableName() string {
	return "public.keys"
}

type PreKey struct {
	KeyId     int64  `gorm:"column:keyId" json:"keyId"`
	PublicKey string `gorm:"column:publicKey" json:"publicKey"`
}

type PreKeyState struct {
	PreKeys      []PreKey     `gorm:"column:preKeys" json:"preKeys"`
	IdentityKey  string       `gorm:"column:identityKey" json:"identityKey"`
	SignedPreKey SignedPreKey `gorm:"column:signedPreKey" json:"signedPreKey"`
}

type DirectoryStat struct {
	Token string `json:"token"`
	Relay string `json:"r"`
	Voice bool   `json:"v"`
	Video bool   `json:"w"`
}

type DirectoryResp struct {
	Token string `json:"token"`
	Relay string `json:"relay,omitempty"`
	Voice bool   `json:"voice"`
	Video bool   `json:"video"`
}

type DirectoryTokens struct {
	Contacts []string
}

type DirectoryStats struct {
	Contacts []DirectoryResp `json:"contacts"`
}

// redis token value
type tokenValue struct {
	Relay string `json:"r"`
	Voice bool   `json:"v"`
	Video bool   `json:"w"`
}

func UpdateDirectory(hs []byte, voice bool, video bool) (err error) {
	var ds tokenValue
	ds.Voice = voice
	ds.Video = video
	ds.Relay = ""
	v, _ := json.Marshal(ds)
	_, err = redisClient.RedisDirectoryManager().HSet("directory", string(hs[:]), string(v)).Result()
	return
}

func (k *Keys) UpdateKeys(number string, deviceId int64, keys []PreKey) {
	pgorm.AccountManager().Where("number = ? and device_id = ?", number, deviceId).Delete(&Keys{})
	for _, key := range keys {
		var ks Keys
		ks.Number = number
		ks.DeviceId = deviceId
		ks.PublicKey = key.PublicKey
		ks.KeyId = key.KeyId
		ks.LastResort = 0
		pgorm.AccountManager().Create(ks)
	}
	return
}

func (k *Keys) GetKeysCount(number string, deviceId int64) int {
	var count int = 0
	pgorm.AccountManager().Model(k).Where("number = ? and device_id = ?", number, deviceId).Count(&count)
	return count
}

func (k *Keys) GetKeysFirst(number string, deviceId int64) {
	pgorm.AccountManager().Where("number = ? and device_id = ?", number, deviceId).Order("key_id").First(k)
	return
}