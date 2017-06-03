package libkademlia

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	mathrand "math/rand"
	"time"
	"sss"
	"fmt"
	"log"
)

type VanashingDataObject struct {
	AccessKey  int64
	Ciphertext []byte
	NumberKeys byte
	Threshold  byte
}

func GenerateRandomCryptoKey() (ret []byte) {
	for i := 0; i < 32; i++ {
		ret = append(ret, uint8(mathrand.Intn(256)))
	}
	return
}

func GenerateRandomAccessKey() (accessKey int64) {
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	accessKey = r.Int63()
	return
}

func CalculateSharedKeyLocations(accessKey int64, count int64) (ids []ID) {
	r := mathrand.New(mathrand.NewSource(accessKey))
	ids = make([]ID, count)
	for i := int64(0); i < count; i++ {
		for j := 0; j < IDBytes; j++ {
			ids[i][j] = uint8(r.Intn(256))
		}
	}
	return
}

func encrypt(key []byte, text []byte) (ciphertext []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphertext = make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], text)
	return
}

func decrypt(key []byte, ciphertext []byte) (text []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext is not long enough")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}

func (k *Kademlia) VanishData(data []byte, numberKeys byte, threshold byte, timeoutSeconds int) (vdo VanashingDataObject) {
	cryptographicKey := GenerateRandomCryptoKey()
	encryptedText := encrypt(cryptographicKey, data)
	sssKeys, err_0 := sss.Split(numberKeys, threshold, cryptographicKey)

	if err_0 != nil {
		fmt.Println("Split cryptographicKey failed")
		return
	}

	accessKey := GenerateRandomAccessKey()
	IDs := CalculateSharedKeyLocations(accessKey, int64(numberKeys))

	fail_count := 0
	for i := 0; i < int(numberKeys); i++ {
		_ , err_1 := k.DoIterativeStore(IDs[i], append([]byte{byte(i)}, sssKeys[byte(i)]...))
		if err_1 != nil {
			fail_count += 1
		}
	}

	fmt.Println(fail_count, "out of", int64(numberKeys), "failed")

	if fail_count >= int(numberKeys) - int(threshold) {
		fmt.Println("Less than threshold sssKeys are stored")
		return
	}

	vdo.AccessKey = accessKey
	vdo.Ciphertext = encryptedText
	vdo.NumberKeys = numberKeys
	vdo.Threshold = threshold
	return
}

func (k *Kademlia) UnvanishData(vdo VanashingDataObject) (data []byte) {
	secret := make(map[byte][]byte)

	IDs := CalculateSharedKeyLocations(vdo.AccessKey, int64(vdo.NumberKeys))

	find_count := 0
	for i := 0; i < int(vdo.NumberKeys); i++ {
		value , err := k.DoIterativeFindValue(IDs[i])
		if err == nil {

			if (len(value[1:]) != 0) {
				find_count += 1
				fmt.Println("value[1:]", value[1:])
				secret[value[0]] = value[1:]
			} else {
				continue
			}

			if find_count == int(vdo.Threshold) {
				break
			}
		}
	}

	if find_count < int(vdo.Threshold) {
		fmt.Println("Less than threshold pieces are found")
		return nil
	}

	log.Println("length is:", len(secret))
	cryptographicKey := sss.Combine(secret)

	return decrypt(cryptographicKey, vdo.Ciphertext)
}
