package api

import (
	"testing"

	"github.com/studygolang/studygolang/internal/model"
)

func TestCanViewTopic(t *testing.T) {
	owner := &model.Me{Uid: 10}
	member := &model.Me{Uid: 20}
	vip := &model.Me{Uid: 30, IsVip: true}
	root := &model.Me{Uid: 40, IsRoot: true}

	tests := []struct {
		name       string
		permission int
		me         *model.Me
		want       bool
	}{
		{name: "public anonymous", permission: model.PermissionPublic, want: true},
		{name: "login anonymous", permission: model.PermissionLogin, want: false},
		{name: "login member", permission: model.PermissionLogin, me: member, want: true},
		{name: "pay anonymous", permission: model.PermissionPay, want: false},
		{name: "pay member", permission: model.PermissionPay, me: member, want: false},
		{name: "pay vip", permission: model.PermissionPay, me: vip, want: true},
		{name: "pay owner", permission: model.PermissionPay, me: owner, want: true},
		{name: "pay root", permission: model.PermissionPay, me: root, want: true},
		{name: "only me member", permission: model.PermissionOnlyMe, me: member, want: false},
		{name: "only me owner", permission: model.PermissionOnlyMe, me: owner, want: true},
		{name: "only me root", permission: model.PermissionOnlyMe, me: root, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canViewTopic(tt.permission, 10, tt.me); got != tt.want {
				t.Fatalf("canViewTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}
