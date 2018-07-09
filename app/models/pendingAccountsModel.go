package models

import (
	pgorm "../../db/gorm"
)

type PendingAccounts struct {
	Id               int64  `gorm:"primary_key" json:"id"`
	Number           string `gorm:"column:number" json:"number"`
	VerificationCode string `gorm:"column:verification_code" json:"verification_code"`
	Timestamp        int64  `gorm:"column:timestamp" json:"timestamp"`
}

func (PendingAccounts) TableName() string {
	return "public.pending_accounts"
}

// 根据ID获取用户信息
func (u *PendingAccounts) GetPendingAccounts(userId int64) {
	pgorm.AccountManager().First(u, userId)
	return
}

func (u *PendingAccounts) SaveOrUpdatePendingAccounts() {
	pgorm.AccountManager().Where(PendingAccounts{Number: u.Number}).Assign(PendingAccounts{VerificationCode: u.VerificationCode, Timestamp: u.Timestamp}).FirstOrCreate(u)
	return
}

func (u *PendingAccounts) RemovePendingAccount() {
	pgorm.AccountManager().Where("number = ?", u.Number).Delete(&PendingAccounts{})
	return
}
