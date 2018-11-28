// Copyright 2018 Jake van der Putten.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"time"

	"github.com/gbrlsnchs/jwt"
)

// Function to sign a JWT token
// INPUT: secret, user, audience as string
// INPUT: the validation time of the token in time.Duration
// OUTPUT: token in slice of bytes if successful and an error if applicable
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

// Function to validate a JWT token
// INPUT: token, secret, user, audience as string
// OUTPUT: whether the token is valid or not as a bool and an error if applicable
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
