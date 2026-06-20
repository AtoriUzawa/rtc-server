package xlog

import "go.uber.org/zap"

// SafeString returns a generic safe field
func SafeString(key, val string) zap.Field {
	return zap.String(key, maskDefault(val))
}

// Token returns a masked "token" field for safe logging.
func Token(val string) zap.Field {
	return zap.String("token", maskDefault(val))
}

// Password returns a masked "password" field for safe logging.
func Password(val string) zap.Field {
	return zap.String("password", maskPassword(val))
}

// Email returns a masked "email" field for safe logging.
func Email(val string) zap.Field {
	return zap.String("email", maskEmail(val))
}

// Phone returns a masked "phone" field for safe logging.
func Phone(val string) zap.Field {
	return zap.String("phone", maskPhone(val))
}
