// Buggeon - SelfHosted service for bug and task tracking
// Copyright (C) 2026 DEVE corp.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	SaltLength = 16
	KeyLength  = 32
	Time       = 3
	Memory     = 64 * 1024
	Threads    = 4
)

type Argon2Hash struct {
	Hash   string
	Salt   string
	Params string
}

type Argon2Params struct {
	Version int
	Memory  uint32
	Time    uint32
	Threads uint8
}

func HashPassword(password string) (string, error) {

	salt := make([]byte, SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		Time,
		Memory,
		Threads,
		KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		Memory,
		Time,
		Threads,
		b64Salt,
		b64Hash,
	)

	return encoded, nil

}

func VerifyPassword(password, encodedHash string) (bool, error) {

	params, salt, hash, err := decodeHash(encodedHash)

	if err != nil {
		return false, fmt.Errorf("Failed to decode hash: %w", err)
	}

	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)

	if err != nil {
		return false, fmt.Errorf("Failed to decode salt: %w", err)
	}

	newHash := argon2.IDKey(
		[]byte(password),
		saltBytes,
		params.Time,
		params.Memory,
		params.Threads,
		uint32(len(hash)),
	)

	return compareHash(hash, newHash), nil

}

func decodeHash(encodedHash string) (Argon2Params, string, []byte, error) {

	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return Argon2Params{}, "", nil, fmt.Errorf("Invalid hash format")
	}

	var version int

	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return Argon2Params{}, "", nil, err
	}

	var memory, time uint32
	var threads uint8

	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return Argon2Params{}, "", nil, err
	}

	salt := parts[4]
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])

	if err != nil {
		return Argon2Params{}, "", nil, err
	}

	return Argon2Params{
		Version: version,
		Memory:  memory,
		Threads: threads,
		Time:    time,
	}, salt, hash, nil

}

func compareHash(a, b []byte) bool {

	if len(a) != len(b) {
		return false
	}

	result := 0

	for i := 0; i < len(a); i++ {
		result |= int(a[i] ^ b[i])
	}

	return result == 0
}
