package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

type CryptoStream struct {
	cipher        cipher.AEAD
	chunkSize     int
	readBuffer    []byte
	decryptBuffer []byte
	nonceBuffer   []byte
	writeBuffer   []byte
	wrapped       io.ReadWriteCloser
}

func NewDefaultCipherStream(stream io.ReadWriteCloser, key []byte) (*CryptoStream, error) {
	return NewChaCha20CryptoStream(stream, key)
}

func NewChaCha20CryptoStream(stream io.ReadWriteCloser, key []byte) (*CryptoStream, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("wrong key length")
	}

	cipher, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed creating cipher: %w", err)
	}

	return newCryptoStream(stream, cipher)
}

func NewAesCryptoStream(stream io.ReadWriteCloser, key []byte) (*CryptoStream, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("wrong key length")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed creating cipher: %w", err)
	}

	aesGcm, err := cipher.NewGCMWithNonceSize(block, 32)
	if err != nil {
		return nil, fmt.Errorf("failed creating GCM: %w", err)
	}

	return newCryptoStream(stream, aesGcm)
}

func newCryptoStream(stream io.ReadWriteCloser, cipher cipher.AEAD) (*CryptoStream, error) {

	bufferSize := 65535 + 2 + cipher.NonceSize()

	chunkSize := 65535 - cipher.Overhead()

	readBuffer := make([]byte, bufferSize)
	writeBuffer := make([]byte, bufferSize)

	nonceBuffer := make([]byte, cipher.NonceSize())

	return &CryptoStream{
		cipher:        cipher,
		chunkSize:     chunkSize,
		readBuffer:    readBuffer,
		decryptBuffer: readBuffer[:0],
		nonceBuffer:   nonceBuffer,
		writeBuffer:   writeBuffer,
		wrapped:       stream,
	}, nil
}

func (c *CryptoStream) Read(b []byte) (int, error) {

	if len(c.decryptBuffer) == 0 {
		var err error
		c.decryptBuffer, err = c.readChunk()
		if err != nil {
			return 0, fmt.Errorf("failed to read encrypted chunk: %w", err)
		}
	}

	n := copy(b, c.decryptBuffer)
	c.decryptBuffer = c.decryptBuffer[n:]
	return n, nil
}

func (c *CryptoStream) readChunk() ([]byte, error) {

	_, err := io.ReadFull(c.wrapped, c.readBuffer[:2])
	if err != nil {
		return nil, fmt.Errorf("failed to read first two bytes: %w", err)
	}
	var length int = int(c.readBuffer[0])<<8 | int(c.readBuffer[1])
	if length == 0 {
		return nil, fmt.Errorf("zero length chunk")
	}
	length += c.cipher.NonceSize()

	_, err = io.ReadFull(c.wrapped, c.readBuffer[:length])
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	nonce := c.readBuffer[:c.cipher.NonceSize()]
	ciphertext := c.readBuffer[c.cipher.NonceSize():int(length)]

	plaintext, err := c.cipher.Open(c.readBuffer[c.cipher.NonceSize():c.cipher.NonceSize()], nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

func (c *CryptoStream) Write(b []byte) (int, error) {
	offset := 0
	for offset < len(b) {
		size := len(b) - offset

		if size > c.chunkSize {
			size = c.chunkSize
		}

		n, err := c.writeChunk(b[offset : offset+size])
		if err != nil {
			return offset, fmt.Errorf("failed to write encrypted chunk: %w", err)
		}

		offset += n
	}

	return offset, nil
}

func (c *CryptoStream) writeChunk(chunk []byte) (int, error) {
	_, err := io.ReadFull(rand.Reader, c.nonceBuffer)
	if err != nil {
		return 0, fmt.Errorf("failed to generate nonce: %w", err)
	}

	copy(c.writeBuffer[2:], c.nonceBuffer)

	ciphertext := c.cipher.Seal(c.writeBuffer[:len(c.nonceBuffer)+2], c.nonceBuffer, chunk, nil)

	size := len(ciphertext) - 2 - c.cipher.NonceSize()

	ciphertext[0] = byte(size >> 8)
	ciphertext[1] = byte(size)

	n, err := c.wrapped.Write(ciphertext)
	if err != nil {
		return n, fmt.Errorf("failed to write encrypted chunk: %w", err)
	}

	return len(chunk), nil
}

func (c *CryptoStream) Close() error {
	return c.wrapped.Close()
}
