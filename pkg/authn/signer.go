/*
 * wraith - a game engine
 * Copyright (c) 2022 Michael D. Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package authn

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"time"
)

// SigningKey is the data required for an HMAC256 signer.
type SigningKey struct {
	// KeyId is the unique identifier for the key.
	keyId string
	// Secret is the plain-text key value using for signing.
	secret []byte
	// ExpiresAt is the moment the key is no longer valid.
	// It is always stored as UTC.
	expiresAt time.Time
}

func newKey(id, secret string, expiresAt time.Time) *SigningKey {
	return &SigningKey{
		keyId:     id,
		secret:    []byte(secret),
		expiresAt: expiresAt.UTC(),
	}
}

func (k *SigningKey) Sign(msg string) ([]byte, error) {
	if k == nil || len(k.secret) == 0 {
		return nil, fmt.Errorf("invalid key")
	} else if !time.Now().Before(k.expiresAt) {
		return nil, fmt.Errorf("expired key")
	}

	// sign the message
	hm := hmac.New(sha256.New, k.secret)
	if _, err := hm.Write([]byte(msg)); err != nil {
		return nil, err
	}
	digest := hm.Sum(nil)

	// return a copy of the digest
	s := make([]byte, len(digest))
	copy(s, digest)
	return s, nil
}
