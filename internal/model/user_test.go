// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"testing"
)

func TestGenMd5Passwd(t *testing.T) {
	userLogin := &UserLogin{
		Passwd: "password123",
	}

	err := userLogin.GenMd5Passwd()
	if err != nil {
		t.Fatalf("GenMd5Passwd failed: %v", err)
	}

	if userLogin.Passcode == "" {
		t.Error("Passcode should not be empty for MD5 password")
	}

	if userLogin.PasswdType != "md5" {
		t.Errorf("Expected passwd_type='md5', got '%s'", userLogin.PasswdType)
	}

	t.Logf("MD5 password generated: passcode=%s, passwd=%s", userLogin.Passcode, userLogin.Passwd)
}

func TestGenHashedPasswd(t *testing.T) {
	userLogin := &UserLogin{
		Passwd: "password123",
	}

	err := userLogin.GenHashedPasswd()
	if err != nil {
		t.Fatalf("GenHashedPasswd failed: %v", err)
	}

	if userLogin.Passcode != "" {
		t.Error("Passcode should be empty for bcrypt password")
	}

	if userLogin.PasswdType != "bcrypt" {
		t.Errorf("Expected passwd_type='bcrypt', got '%s'", userLogin.PasswdType)
	}

	t.Logf("Bcrypt password generated: passwd=%s", userLogin.Passwd)
}

func TestVerifyPasswd_MD5(t *testing.T) {
	// Test MD5 password verification (legacy)
	userLogin := &UserLogin{
		Passwd: "password123",
	}

	err := userLogin.GenMd5Passwd()
	if err != nil {
		t.Fatalf("GenMd5Passwd failed: %v", err)
	}

	// Correct password
	if !userLogin.VerifyPasswd("password123") {
		t.Error("MD5 password verification failed for correct password")
	}

	// Wrong password
	if userLogin.VerifyPasswd("wrongpassword") {
		t.Error("MD5 password verification should fail for wrong password")
	}
}

func TestVerifyPasswd_Bcrypt(t *testing.T) {
	// Test bcrypt password verification
	userLogin := &UserLogin{
		Passwd: "password123",
	}

	err := userLogin.GenHashedPasswd()
	if err != nil {
		t.Fatalf("GenHashedPasswd failed: %v", err)
	}

	// Correct password
	if !userLogin.VerifyPasswd("password123") {
		t.Error("Bcrypt password verification failed for correct password")
	}

	// Wrong password
	if userLogin.VerifyPasswd("wrongpassword") {
		t.Error("Bcrypt password verification should fail for wrong password")
	}
}

func TestPasswordMigration(t *testing.T) {
	// Simulate migration from MD5 to bcrypt
	oldUser := &UserLogin{
		Passwd: "testpassword",
	}

	// Step 1: Generate MD5 password (old user)
	err := oldUser.GenMd5Passwd()
	if err != nil {
		t.Fatalf("GenMd5Passwd failed: %v", err)
	}

	// Step 2: Verify MD5 password works
	if !oldUser.VerifyPasswd("testpassword") {
		t.Error("MD5 password verification failed")
	}

	// Step 3: Simulate auto-upgrade to bcrypt (on login)
	oldUser.Passwd = "testpassword" // Restore plain password for upgrade
	err = oldUser.GenHashedPasswd()
	if err != nil {
		t.Fatalf("GenHashedPasswd failed: %v", err)
	}

	// Step 4: Verify bcrypt password works
	if !oldUser.VerifyPasswd("testpassword") {
		t.Error("Bcrypt password verification failed after migration")
	}

	t.Log("Password migration from MD5 to bcrypt successful")
}
