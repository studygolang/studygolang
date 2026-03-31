# Password Hashing Migration Guide

## Overview

This document describes the migration from MD5 to bcrypt password hashing in the StudyGolang project.

## Background

The original system used MD5 hashing with a random salt (passcode):
```
hash = md5(password + passcode)
```

This approach is no longer considered secure. We are migrating to bcrypt, which provides:
- Adaptive hashing (configurable cost)
- Built-in salt management
- Resistance to rainbow table attacks

## Migration Strategy

The migration is **gradual** and **backward compatible**:

1. **Dual-mode verification**: The system supports both MD5 and bcrypt passwords
2. **Auto-upgrade on login**: When users with MD5 passwords log in, their passwords are automatically upgraded to bcrypt
3. **New users use bcrypt**: All new registrations use bcrypt immediately

## Database Changes

### Schema Migration

Run the following SQL to add the `passwd_type` column:

```bash
mysql -u username -p database_name < migrations/20260329_add_passwd_type_column.sql
```

### Schema Definition

```sql
ALTER TABLE user_login
ADD COLUMN passwd_type VARCHAR(10) DEFAULT 'md5'
COMMENT 'password type: md5(legacy) or bcrypt(new)';
```

## Code Changes

### Model Layer (`internal/model/user.go`)

1. Added `PasswdType` field to `UserLogin` struct
2. Added `GenHashedPasswd()` - generates bcrypt hash for new users
3. Added `VerifyPasswd()` - unified verification supporting both MD5 and bcrypt
4. Kept `GenMd5Passwd()` for backward compatibility

### Logic Layer (`internal/logic/user.go`)

1. **Login**: Auto-upgrades MD5 passwords to bcrypt on successful login
2. **Register**: Uses bcrypt for new users (`doCreateUser`)
3. **UpdatePasswd**: Uses bcrypt when users change passwords
4. **ResetPasswd**: Uses bcrypt when passwords are reset

## Verification Steps

### 1. Run Tests

```bash
go test -v ./internal/model/
```

Expected output:
```
=== RUN   TestGenMd5Passwd
--- PASS: TestGenMd5Passwd
=== RUN   TestGenHashedPasswd
--- PASS: TestGenHashedPasswd
=== RUN   TestVerifyPasswd_MD5
--- PASS: TestVerifyPasswd_MD5
=== RUN   TestVerifyPasswd_Bcrypt
--- PASS: TestVerifyPasswd_Bcrypt
=== RUN   TestPasswordMigration
--- PASS: TestPasswordMigration
PASS
```

### 2. Manual Testing

1. **Test existing user login (MD5)**
   - Login with existing account
   - Check logs for "auto-upgrading user password to bcrypt"
   - Verify database: `passwd_type` should be 'bcrypt' after login

2. **Test new user registration**
   - Register a new account
   - Verify database: `passwd_type` should be 'bcrypt'
   - `passcode` should be empty

3. **Test password change**
   - Change password for any user
   - Verify database: `passwd_type` should be 'bcrypt'

## Monitoring

### Key Metrics to Monitor

1. **Password type distribution**:
   ```sql
   SELECT passwd_type, COUNT(*) as count
   FROM user_login
   GROUP BY passwd_type;
   ```

2. **Failed password upgrades**: Check logs for "failed to upgrade password to bcrypt"

3. **Login failures**: Monitor for any increase in login errors

## Rollback Plan

If issues arise:

1. **Revert code changes**:
   ```bash
   git revert <commit-hash>
   ```

2. **Database remains compatible** (MD5 passwords still work)

3. **No data loss**: Users with bcrypt passwords can still login (code will check both types)

## Security Considerations

- **No password exposure**: The migration happens in-memory; plain passwords are never logged
- **Graceful degradation**: If bcrypt upgrade fails, login still succeeds (MD5 verified)
- **Audit trail**: All upgrades are logged for security review

## Timeline

- **2026-03-29**: Migration deployed
- **Ongoing**: Monitor password type distribution
- **Target**: 90%+ users migrated within 30 days (through natural login)

## References

- [bcrypt package documentation](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
