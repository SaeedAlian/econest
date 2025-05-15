package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/SaeedAlian/econest/api/types"
)

var ctx = context.Background()

type KeyPair struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type KeyServer struct {
	cache *redis.Client
}

func NewKeyServer(cache *redis.Client) *KeyServer {
	return &KeyServer{
		cache: cache,
	}
}

func (ks *KeyServer) RotateKeys(newKID string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicKeyBytes})

	if err := ks.cache.Set(ctx, "jwt:public:"+newKID, pubPEM, 7*24*time.Hour).Err(); err != nil {
		return err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes})

	if err := ks.cache.Set(ctx, "jwt:private:"+newKID, privPEM, 7*24*time.Hour).Err(); err != nil {
		return err
	}

	if err := ks.cache.Set(ctx, "jwt:current_kid", newKID, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (ks *KeyServer) GetCurrentKID() (string, error) {
	val, err := ks.cache.Get(ctx, "jwt:current_kid").Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (ks *KeyServer) GetPrivateKey(kid string) (*rsa.PrivateKey, error) {
	privPEM, err := ks.cache.Get(ctx, "jwt:private:"+kid).Bytes()
	if err != nil {
		return nil, err
	}

	b, _ := pem.Decode(privPEM)
	if b == nil || b.Type != "RSA PRIVATE KEY" {
		return nil, types.ErrInvalidPEMBlockForPrivateKey
	}

	privKey, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func (ks *KeyServer) GetPublicKey(kid string) (*rsa.PublicKey, error) {
	pubPEM, err := ks.cache.Get(ctx, "jwt:public:"+kid).Bytes()
	if err != nil {
		return nil, err
	}

	b, _ := pem.Decode(pubPEM)
	if b == nil || b.Type != "RSA PUBLIC KEY" {
		return nil, types.ErrInvalidPEMBlockForPublicKey
	}

	pubKey, err := x509.ParsePKCS1PublicKey(b.Bytes)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}
