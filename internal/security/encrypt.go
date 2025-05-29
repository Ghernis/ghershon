package encrypt

import (
	"fmt"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
	"errors"
	"os"
	"path/filepath"
    "io"
	"runtime"
)

const (
    configDirName = "ghershon"
    keyFileName   = "enck"
    keySize       = 32 // 256-bit AES key
)

func isWindows() bool {
    return runtime.GOOS == "windows"
}

func EnsureEncryptionKey() (string, error) {
    configDir, err := os.UserConfigDir()
    if err != nil {
        return "", fmt.Errorf("failed to find home dir: %w", err)
    }
    keyDir := filepath.Join(configDir, configDirName)
    keyPath := filepath.Join(keyDir, keyFileName)

    // Create directory if not exists
    if _, err := os.Stat(keyDir); os.IsNotExist(err) {
        if err := os.MkdirAll(keyDir, 0700); err != nil {
            return "", fmt.Errorf("failed to create config dir: %w", err)
        }
    }

    // Create key if not exists
    if _, err := os.Stat(keyPath); os.IsNotExist(err) {
        key := make([]byte, keySize)
        if _, err := rand.Read(key); err != nil {
            return "", fmt.Errorf("failed to generate key: %w", err)
        }

        encoded := base64.StdEncoding.EncodeToString(key)
        if err := os.WriteFile(keyPath, []byte(encoded), 0600); err != nil {
            return "", fmt.Errorf("failed to write key: %w", err)
        }
    }


    // Check permissions
	if !isWindows(){
		info, err := os.Stat(keyPath)
		if err != nil {
			return "", fmt.Errorf("failed to stat key file: %w", err)
		}

		if info.Mode().Perm() != 0600 {
			return "", errors.New("key file has incorrect permissions; expected 0600")
		}
	}

    return keyFileName, nil
}

func EncryptText(plaintext string, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, aesGCM.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptText(encodedCipherText string, key []byte) (string, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encodedCipherText)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := aesGCM.NonceSize()
    if len(ciphertext) < nonceSize {
        return "", fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}
