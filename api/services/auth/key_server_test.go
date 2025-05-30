package auth

import (
	"log"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"

	"github.com/SaeedAlian/econest/api/config"
)

func TestKeyServer(t *testing.T) {
	if config.Env.Env != "test" {
		log.Panic("environment is not on test!!")
		os.Exit(1)
	}

	ksCache := redis.NewClient(&redis.Options{
		Addr: config.Env.KeyServerRedisAddr,
	})

	ks := NewKeyServer(ksCache)

	kid1 := "kid1"

	t.Run("should rotate keys successfully", func(t *testing.T) {
		err := ks.RotateKeys(kid1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should get current kid without issue", func(t *testing.T) {
		currKID, err := ks.GetCurrentKID()
		if err != nil {
			t.Fatal(err)
		}

		if currKID != kid1 {
			t.Errorf("expected value %s for kid but got %s\n", kid1, currKID)
		}
	})

	t.Run("should get current private key without issue", func(t *testing.T) {
		currKID, err := ks.GetCurrentKID()
		if err != nil {
			t.Fatal(err)
		}

		_, err = ks.GetPrivateKey(currKID)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should get current public key without issue", func(t *testing.T) {
		currKID, err := ks.GetCurrentKID()
		if err != nil {
			t.Fatal(err)
		}

		_, err = ks.GetPublicKey(currKID)
		if err != nil {
			t.Fatal(err)
		}
	})
}
