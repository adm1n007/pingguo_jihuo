package rsa

import (
    . "active_apple/ml/trace"
    . "active_apple/ml/strings"
    . "fmt"

    "io"
    "math/big"
    "crypto/rand"
    "crypto/x509"
    "crypto/rsa"
    "crypto/subtle"
    "encoding/pem"
)

type RSACipher struct {
    privateKey  *rsa.PrivateKey
    publicKey   *rsa.PublicKey
}

func NewRSACipher(key interface{}) *RSACipher {
    cipher := &RSACipher{}

    switch v := key.(type) {
        case []byte:
            cipher.loadKey(v)

        case string:
            cipher.loadKey([]byte(v))

        case String:
            cipher.loadKey([]byte(v))

        default:
            Raise(NewRSAInvalidKeyError(Sprintf("unsupported key type: %T", key)))
    }

    return cipher
}

func (self *RSACipher) loadKey(key []byte) {
    var err error

    block, _ := pem.Decode(key)
    if block != nil {
        key = block.Bytes
    }

    privKey, err := x509.ParsePKCS1PrivateKey(key)
    if err == nil {
        self.loadPrivateKey(privKey)
        return
    }

    pubkey, err := x509.ParsePKIXPublicKey(key)
    RaiseIf(err)

    if _, ok := pubkey.(*rsa.PublicKey); ok == false {
        Raise(NewRSAInvalidKeyError("invalid rsa public key"))
    }

    self.loadPublicKey(pubkey.(*rsa.PublicKey))
}

func (self *RSACipher) loadPrivateKey(privateKey *rsa.PrivateKey) {
    self.privateKey = privateKey
}

func (self *RSACipher) loadPublicKey(publicKey *rsa.PublicKey) {
    self.publicKey = publicKey
}

func (self *RSACipher) PublicKey() rsa.PublicKey {
    return *self.publicKey
}

func (self *RSACipher) PrivateKey() rsa.PrivateKey {
    return *self.privateKey
}

func (self *RSACipher) EncryptPKCS1v15(plain []byte) (crypto []byte) {
    var err error

    if self.publicKey != nil {
        crypto = encryptPKCS1v15(plain, big.NewInt(int64(self.publicKey.E)), self.publicKey.N)

    } else if self.privateKey != nil {
        crypto = encryptPKCS1v15(plain, self.privateKey.D, self.privateKey.N)

    } else {
        err = NewRSAInvalidKeyError("key has not been set")
    }

    RaiseIf(err)
    return
}

func (self *RSACipher) DecryptPKCS1v15(crypto []byte) (plain []byte) {
    var n *big.Int

    if self.privateKey != nil {
        plain = decrypt(new(big.Int).SetBytes(crypto), self.privateKey.D, self.privateKey.N).Bytes()
        n = self.privateKey.N

    } else if self.publicKey != nil {
        plain = decrypt(new(big.Int).SetBytes(crypto), big.NewInt(int64(self.publicKey.E)), self.publicKey.N).Bytes()
        n = self.publicKey.N
    }

    em, index := decryptPKCS1v15(n, plain)
    plain = em[index:]

    return
}

func encrypt(msg *big.Int, e *big.Int, n *big.Int) *big.Int {
    return new(big.Int).Exp(msg, e, n)
}

func decrypt(msg *big.Int, d *big.Int, n *big.Int) *big.Int {
    return new(big.Int).Exp(msg, d, n)
}

func leftPad(input []byte, size int) (out []byte) {
    n := len(input)
    if n > size {
        n = size
    }
    out = make([]byte, size)
    copy(out[len(out)-n:], input)
    return
}

func copyWithLeftPad(dest, src []byte) {
    numPaddingBytes := len(dest) - len(src)
    for i := 0; i < numPaddingBytes; i++ {
        dest[i] = 0
    }
    copy(dest[numPaddingBytes:], src)
}

func decryptPKCS1v15(n *big.Int, plain []byte) (em []byte, index int) {
    k := (n.BitLen() + 7) / 8
    if k < 11 {
        Raise(rsa.ErrDecryption)
    }

    m := new(big.Int).SetBytes(plain)

    em = leftPad(m.Bytes(), k)
    firstByteIsZero := subtle.ConstantTimeByteEq(em[0], 0)
    secondByteIsTwo := subtle.ConstantTimeByteEq(em[1], 2)

    lookingForIndex := 1
    for i := 2; i < len(em); i++ {
        equals0 := subtle.ConstantTimeByteEq(em[i], 0)
        index = subtle.ConstantTimeSelect(lookingForIndex & equals0, i, index)
        lookingForIndex = subtle.ConstantTimeSelect(equals0, 0, lookingForIndex)
    }

    validPS := subtle.ConstantTimeLessOrEq(2+8, index)

    valid := firstByteIsZero & secondByteIsTwo & (^lookingForIndex & 1) & validPS
    if valid == 0 {
        Raise(rsa.ErrDecryption)
    }

    index = subtle.ConstantTimeSelect(valid, index+1, 0)
    return em, index
}

func encryptPKCS1v15(msg []byte, e *big.Int, n *big.Int) (out []byte) {
    k := (n.BitLen() + 7) / 8
    if len(msg) > k - 11 {
        Raise(rsa.ErrMessageTooLong)
    }

    // EM = 0x00 || 0x02 || PS || 0x00 || M
    em := make([]byte, k)
    em[1] = 2
    ps, mm := em[2:len(em)-len(msg)-1], em[len(em)-len(msg):]
    nonZeroRandomBytes(ps, rand.Reader)

    em[len(em)-len(msg)-1] = 0
    copy(mm, msg)

    m := new(big.Int).SetBytes(em)
    c := encrypt(m, e, n)

    copyWithLeftPad(em, c.Bytes())
    out = em
    return
}

func nonZeroRandomBytes(s []byte, rand io.Reader) {
    var err error

    _, err = io.ReadFull(rand, s)
    RaiseIf(err)

    for i := 0; i < len(s); i++ {
        for s[i] == 0 {
            _, err = io.ReadFull(rand, s[i:i+1])
            RaiseIf(err)

            s[i] ^= 0x42
        }
    }

    return
}
