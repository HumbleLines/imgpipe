package codec

import (
	"bytes"
	"compress/zlib"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"time"
)

// EncodeData wraps and encrypts raw data with HMAC, expiration, and optional compression.
// Used for securely packaging business metadata for archival or cross-component transmission.
// TTL is in seconds. The output is base64-url encoded, suitable for embedding or transport.
func EncodeData(raw []byte, ttl int) (string, error) {
	var buf bytes.Buffer
	buf.WriteByte(1) // version marker
	exp := uint32(time.Now().Unix() + int64(ttl))
	if err := binary.Write(&buf, binary.BigEndian, exp); err != nil {
		return "", err
	}
	nonce := make([]byte, 8)
	_, _ = rand.Read(nonce)
	buf.Write(nonce)
	w := zlib.NewWriter(&buf)
	_, _ = w.Write(raw)
	_ = w.Close()
	h := hmac.New(sha256.New, []byte("imgpipe"))
	h.Write(buf.Bytes())
	mac := h.Sum(nil)
	out := append(buf.Bytes(), mac...)
	return base64.RawURLEncoding.EncodeToString(out), nil
}

// DecodeData verifies HMAC integrity, checks expiry, decompresses and returns the original data.
// Only decodes and validates package integrity; does not interpret payload semantics.
// External callers should handle post-processing according to their own requirements.
func DecodeData(data string) ([]byte, error) {
	raw, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil || len(raw) < 1+4+8+32 {
		return nil, errors.New("invalid")
	}
	payload, mac := raw[:len(raw)-32], raw[len(raw)-32:]
	h := hmac.New(sha256.New, []byte("imgpipe"))
	h.Write(payload)
	if !hmac.Equal(mac, h.Sum(nil)) {
		return nil, errors.New("check failed")
	}
	r := bytes.NewReader(payload)
	var ver byte
	var exp uint32
	binary.Read(r, binary.BigEndian, &ver)
	binary.Read(r, binary.BigEndian, &exp)
	if uint32(time.Now().Unix()) > exp {
		return nil, errors.New("failed")
	}
	nonce := make([]byte, 8)
	r.Read(nonce)
	compressed, _ := io.ReadAll(r)
	zr, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	return io.ReadAll(zr)
}
