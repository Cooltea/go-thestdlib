package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
)

const (
    KeyFile       = "aes.%d.key"
    EncryptedFile = "aes.go.%d.enc"
)

var (
    IV      = []byte("batman and robin")
    keySize = flag.Int("keysize", 32, "The keysize in bytes to use: 16, 24, or 32 (default)")
    do      = flag.String("do", "enc", "Which operation to perform: enc (encryption, default) or dec (decryption)")
)

func MakeKey() []byte {
    key := make([]byte, *keySize)
    n, err := rand.Read(key)
    if err != nil {
        log.Fatalf("Failed to read new random key: %s", err)
    }
    if n < *keySize {
        log.Fatalf("Failed to read entire key, only read %d out of %d", n, *keySize)
    }
    return key
}

func Key() []byte {
    file := fmt.Sprintf(KeyFile, *keySize)
    key, err := ioutil.ReadFile(file)
    if err != nil {
        log.Println("Failed reading keyfile, making a new one...")
        key = MakeKey()
        err = ioutil.WriteFile(file, key, 0644)
        if err != nil {
            log.Fatalf("Failed saving key to %s: %s", file, err)
        }
        return key
    }
    return key
}

func MakeCipher() cipher.Block {
    cipher, err := aes.NewCipher(Key())
    if err != nil {
        log.Fatalf("Failed making the AES cipher: %s", err)
    }
    return cipher
}

func Crypt(input, output string) {
    blockCipher := MakeCipher()
    stream := cipher.NewCTR(blockCipher, IV)
    bytes, err := ioutil.ReadFile(input)
    if err != nil {
        log.Fatalf("Failed reading input file: %s", err)
    }
    stream.XORKeyStream(bytes, bytes)
    err = ioutil.WriteFile(output, bytes, 0644)
    if err != nil {
        log.Fatalf("Failed writing output file: %s", err)
    }
}

func Encrypt() {
    Crypt("aes.go", fmt.Sprintf(EncryptedFile, *keySize))
}

func Decrypt() {
    Crypt(fmt.Sprintf(EncryptedFile, *keySize), "aes.go.dec")
}

func main() {
    flag.Parse()

    switch *keySize {
    case 16, 24, 32:
        // Keep calm and carry on...
    default:
        log.Fatalf("%d is not a valid keysize. Must be one of 16, 24, 32", *keySize)
    }

    switch *do {
    case "enc":
        Encrypt()
    case "dec":
        Decrypt()
    default:
        log.Fatalf("%s is not a valid operation. Must be one of enc or dec", *do)
    }
}
