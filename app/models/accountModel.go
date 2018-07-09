package models

import (
	pgorm "../../db/gorm"
)

type Accounts struct {
	Id     int64  `gorm:"primary_key" json:"id"`
	Number string `gorm:"column:number" json:"number"`
	Data   string `gorm:"column:data" json:"data"`
}

type SignedPreKey struct {
	KeyId     int64  `gorm:"column:keyId" json:"keyId"`
	PublicKey string `gorm:"column:publicKey" json:"publicKey"`
	Signature string `gorm:"column:signature" json:"signature"`
}

type Device struct {
	Id              int64        `gorm:"primary_key" json:"id"`
	Salt            string       `gorm:"column:salt" json:"salt,omitempty"`
	AuthToken       string       `gorm:"column:authToken" json:"authToken,omitempty"`
	Name            string       `gorm:"column:name" json:"name"`
	FetchesMessages bool         `gorm:"column:fetchesMessages" json:"fetchesMessages,omitempty"`
	SignalingKey    string       `gorm:"column:signalingKey" json:"signalingKey,omitempty"`
	GcmId           string       `gorm:"column:gcmId" json:"gcmId,omitempty"`
	ApnId           string       `gorm:"column:apnId" json:"apnId,omitempty"`
	RegistrationId int          `gorm:"column:registrationId" json:"registrationId"`
	SignedPreKey   SignedPreKey `gorm:"column:signedPreKey" json:"signedPreKey"`
	Voice          bool         `gorm:"column:voice" json:"voice,omitempty"`
	Video          bool         `gorm:"column:video" json:"video,omitempty"`
	LastSeen       int64        `gorm:"column:lastSeen" json:"lastSeen,omitempty"`
	VoipApnId      string       `gorm:"column:voipApnId" json:"voipApnId,omitempty"`
	UerAgent       string       `gorm:"column:userAgent" json:"userAgent,omitempty"`
	Created        int64        `gorm:"column:created" json:"created,omitempty"`
}

type Devices struct {
	Devices []Device
}

type Account struct {
	Number       string   `gorm:"column:number" json:"number"`
	Devices      []Device `gorm:"type:json;column:devices" json:"devices"`
	IdentityKey  string   `gorm:"column:identityKey" json:"identityKey"`
	Name         string   `gorm:"column:name" json:"name"`
	Avatar       string   `gorm:"column:avatar" json:"avatar"`
	AvatarDigest string   `gorm:"column:avatarDigest" json:"avatarDigest"`
	Pin          string   `gorm:"column:pin" json:"pin"`
}

type AccountAttributes struct {
	SignalingKey    string `gorm:"column:signalingKey" json:"signalingKey"`
	FetchesMessages bool   `gorm:"column:fetchesMessages;default:true" json:"fetchesMessages"`
	RegistrationId  int    `gorm:"column:registrationId" json:"registrationId"`
	Name            string `gorm:"column:name" json:"name"`
	Voice           bool   `gorm:"column:voice" json:"voice"`
	Video           bool   `gorm:"column:video" json:"video"`
}

func (Accounts) TableName() string {
	return "public.accounts"
}

// 根据ID获取用户信息
func (u *Accounts) GetAccountByNumber(number string) {
	pgorm.AccountManager().Where("number = ?", number).First(u)
	return
}

func (u *Accounts) CreateAccount() {
	pgorm.AccountManager().Create(u)
	return
}

func (u *Accounts) DeleteAccount() {
	pgorm.AccountManager().Where("number = ?", u.Number).Delete(&Account{})
	return
}

func (u *Accounts) UpdateAccount() (err error) {
	// TODO REDIS Cache
	err = pgorm.AccountManager().Save(u).Error
	return
}
