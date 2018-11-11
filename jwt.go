package main

import (
	"fmt"
	"time"

	"github.com/gbrlsnchs/jwt"
)

func jwtSignToken(s, u, a string, v time.Duration) ([]byte, error) {
	// Timestamp the beginning.
	now := time.Now()
	// Define a signer.
	hs256 := jwt.NewHS256(s)
	jot := &jwt.JWT{
		Issuer:         "MicroHTTP",
		Subject:        u,
		Audience:       a,
		ExpirationTime: now.Add(24 * 30 * 12 * time.Hour).Unix(),
		NotBefore:      now.Add(v).Unix(),
		IssuedAt:       now.Unix(),
	}
	jot.SetAlgorithm(hs256)
	payload, err := jwt.Marshal(jot)
	if err != nil {
		return []byte(""), fmt.Errorf("jwt: marshall error: %s", err)
	}
	token, err := hs256.Sign(payload)
	if err != nil {
		return []byte(""), fmt.Errorf("jwt: sign error: %s", err)
	}
	return token, nil
}

func jwtValidateToken(t, s, u, a string) (bool, error) {
	// Timestamp the beginning.
	now := time.Now()
	// Define a signer.
	hs256 := jwt.NewHS256(s)

	// First, extract the payload and signature.
	// This enables unmarshaling the JWT first and
	// verifying it later or vice versa.
	payload, sig, err := jwt.Parse(t)
	if err != nil {
		return false, fmt.Errorf("jwt: parse error: %s", err)
	}
	if err = hs256.Verify(payload, sig); err != nil {
		return false, fmt.Errorf("jwt: verify error: %s", err)
	}
	var jot jwt.JWT
	if err = jwt.Unmarshal(payload, &jot); err != nil {
		return false, fmt.Errorf("jwt: unmarshall error: %s", err)
	}

	// Validate fields.
	iatValidator := jwt.IssuedAtValidator(now)
	expValidator := jwt.ExpirationTimeValidator(now)
	audValidator := jwt.AudienceValidator(a)
	subValidator := jwt.SubjectValidator(u)
	if err = jot.Validate(iatValidator, expValidator, audValidator, subValidator); err != nil {
		switch err {
		case jwt.ErrIatValidation:
			// handle "iat" validation error
			return false, fmt.Errorf("jwt: validation error: %s", err)
		case jwt.ErrExpValidation:
			// handle "exp" validation error
			return false, fmt.Errorf("jwt: validation error: %s", err)
		case jwt.ErrAudValidation:
			// handle "aud" validation error
			return false, fmt.Errorf("jwt: validation error: %s", err)
		case jwt.ErrSubValidation:
			// handle "sub" validation error
			return false, fmt.Errorf("jwt: validation error: %s", err)
		}
	}
	return true, nil
}
