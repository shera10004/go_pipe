package testcrypto_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"math/big"

	"fmt"
	"testing"
)

func TestSha256(t *testing.T) {

	msg := []byte("abc def")
	hash := sha256.New()
	code := hash.Sum(msg)

	fmt.Printf("code: %x\n", code)
	fmt.Println("block size:", hash.BlockSize())
	fmt.Println("size:", hash.Size())

	code2 := sha256.Sum256(msg)
	fmt.Printf("code2: %x\n", code2)

}

func TestAES(t *testing.T) {

	//encoding
	var ed []byte
	var edstr string
	enkey := "sakdfjkalsdjfoiwjeofrjsdfjskadfk"
	cekey := "nonce6789123askdfksdfksjdfksjkdf"
	{
		fmt.Println("------ aes encoding ----------")
		key := []byte(enkey)
		block, err := aes.NewCipher(key)

		if err != nil {
			fmt.Println("err", err)
			return
		}
		fmt.Println("block size:", block.BlockSize())

		//aead, err := cipher.NewGCM(block)
		aead, err := cipher.NewGCMWithNonceSize(block, 32)
		fmt.Println("nonce size:", aead.NonceSize())
		if err != nil {
			fmt.Println("err", err)
			return
		}

		nonce := []byte(cekey)
		msg := "hi! golang"

		encdata := aead.Seal(nil, nonce, []byte(msg), nil)

		ed = make([]byte, len(encdata))
		copy(ed, encdata)
		fmt.Printf("encdata : %x \n", ed)
		edstr = fmt.Sprintf("%x", ed)
		fmt.Printf("edstr : %s \n", edstr)

	}

	//decoding
	{
		fmt.Println("\n------ aes decoding ----------")
		key := []byte(enkey)
		nonce := []byte(cekey)

		block, _ := aes.NewCipher(key)
		//aead, _ := cipher.NewGCM(block)
		aead, _ := cipher.NewGCMWithNonceSize(block, 32)

		cloneed, _ := hex.DecodeString(edstr)

		result, err := aead.Open(nil, nonce, cloneed, nil)
		if err != nil {
			fmt.Println("err", err)
			return
		}
		fmt.Println(string(result))

	}

}

func TestRSA(t *testing.T) {

	privatekey, err := rsa.GenerateKey(rand.Reader, 1024)
	if checkError(err) {
		return
	}

	_ = privatekey

	publickey := privatekey.Public()
	_ = publickey

	fmt.Println("private", privatekey)
	fmt.Println("public", publickey)

	PublickeyToBytes := func(pk *rsa.PublicKey) []byte {
		re := make([]byte, 4)
		binary.BigEndian.PutUint32(re, uint32(pk.E))

		re = append(re, pk.N.Bytes()...)

		return re
	}
	BytesToPublickey := func(data []byte) *rsa.PublicKey {
		pk := &rsa.PublicKey{}
		e := binary.BigEndian.Uint32(data)
		pk.E = int(e)
		data = data[4:]

		// bint := &big.Int{}
		// bint.SetBytes(data)
		// pk.N = bint

		pk.N = new(big.Int).SetBytes(data)

		return pk
	}

	bytesKey := PublickeyToBytes(publickey.(*rsa.PublicKey))
	fmt.Println("PublickeyToBytes", bytesKey)
	pk := BytesToPublickey(bytesKey)
	fmt.Println("BytesToPublickey", pk)

	fmt.Println("===== PKCS1v15 =====")
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pk, []byte("go"))
	if checkError(err) {
		return
	}
	fmt.Println("ciphertext", ciphertext)

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privatekey, ciphertext)
	if checkError(err) {
		return
	}
	fmt.Println("plaintext", string(plaintext))

	fmt.Println("===== OAEP =====")
	levelsalt := []byte{1}
	_ = levelsalt
	ciphertext, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, pk, []byte("go lang"), levelsalt)
	if checkError(err) {
		return
	}
	fmt.Println("ciphertext", ciphertext)

	plaintext, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, privatekey, ciphertext, levelsalt)
	if checkError(err) {
		return
	}
	fmt.Println("plaintext", string(plaintext))

}

func checkError(err error) bool {
	if err != nil {
		fmt.Println("err", err)
		return true
	}
	return false
}
