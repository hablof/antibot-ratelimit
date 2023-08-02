package service

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/hablof/antibot-ratelimit/internal/config"
	"github.com/stretchr/testify/assert"
)

var (
	testConfig config.Config = config.Config{
		BucketSize:  20,
		RPMLimit:    100,
		PrefixSize:  24,
		BanDuration: 2 * time.Minute,
	}
)

func Test_unaryLimiter_isLimitOK(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Run("fast valid 20 requests", func(t *testing.T) {
		t.Parallel()

		ul := newUnaryLimiter(testConfig.BucketSize, 12*time.Second, testConfig.BanDuration)

		for i := 0; i < 20; i++ {
			ok := ul.isLimitOK()
			assert.Equal(t, true, ok)
		}
	})

	t.Run("fast valid 20 requests and 21st blocked", func(t *testing.T) {
		t.Parallel()

		ul := newUnaryLimiter(testConfig.BucketSize, 12*time.Second, testConfig.BanDuration)

		for i := 0; i < 20; i++ {
			ok := ul.isLimitOK()
			fmt.Println("case 2", time.Now(), ok, len(ul.tokenPool))
			assert.Equal(t, true, ok)
		}

		ok := ul.isLimitOK()
		fmt.Println("case 2", time.Now(), ok, len(ul.tokenPool))
		assert.Equal(t, false, ok)
	})

	t.Run("slow 20 requests, 21st but blocked anyway", func(t *testing.T) {
		t.Parallel()

		ul := newUnaryLimiter(testConfig.BucketSize, 12*time.Second, testConfig.BanDuration)

		for i := 0; i < 20; i++ {
			ok := ul.isLimitOK()
			fmt.Println("case 3", time.Now(), ok, len(ul.tokenPool))
			assert.Equal(t, true, ok)

			time.Sleep((12*time.Second - 50*time.Millisecond) / 21)
		}

		ok := ul.isLimitOK()
		fmt.Println("case 3", time.Now(), ok, len(ul.tokenPool))
		assert.Equal(t, false, ok)
	})

	t.Run("slow enough 25 requests, is not blocked", func(t *testing.T) {
		t.Parallel()

		ul := newUnaryLimiter(testConfig.BucketSize, 12*time.Second, testConfig.BanDuration)

		for i := 0; i < 25; i++ {
			ok := ul.isLimitOK()
			t.Log("case 4", time.Now(), ok, len(ul.tokenPool))
			assert.Equal(t, true, ok)

			time.Sleep((12*time.Second + 50*time.Millisecond) / 20)
		}
	})

	t.Run("fast valid 20 requests, 21st blocked, reset limiter and 22nd is valid again", func(t *testing.T) {
		t.Parallel()

		ul := newUnaryLimiter(testConfig.BucketSize, 12*time.Second, testConfig.BanDuration)

		for i := 0; i < 20; i++ {
			ok := ul.isLimitOK()
			t.Log("case 5", time.Now(), ok, len(ul.tokenPool))
			assert.Equal(t, true, ok)
		}

		ok := ul.isLimitOK()
		t.Log("case 5", time.Now(), ok, len(ul.tokenPool))
		assert.Equal(t, false, ok)

		ul.resetLimit()

		ok = ul.isLimitOK()
		t.Log("case 5", time.Now(), ok, len(ul.tokenPool))
		assert.Equal(t, true, ok)
	})

	t.Run("fast valid 20 requests, reset limiter and 20 more valid requests again", func(t *testing.T) {
		t.Parallel()

		ul := newUnaryLimiter(testConfig.BucketSize, 12*time.Second, testConfig.BanDuration)

		for i := 0; i < 20; i++ {
			ok := ul.isLimitOK()
			t.Log("case 6", time.Now(), ok, len(ul.tokenPool))
			assert.Equal(t, true, ok)
		}

		ul.resetLimit()

		for i := 0; i < 20; i++ {
			ok := ul.isLimitOK()
			t.Log("case 6", time.Now(), ok, len(ul.tokenPool))
			assert.Equal(t, true, ok)
		}
	})
}

func Test_Ratelimiter(t *testing.T) {

	ips := []net.IP{net.ParseIP("110.212.84.67"), net.ParseIP("31.153.109.117"), net.ParseIP("175.54.199.111")}

	t.Run("60 requests between 3 ip", func(t *testing.T) {

		r := NewRatelimiter(testConfig)
		for i := 0; i < 60; i++ {
			assert.Equal(t, true, r.IsLimitOK(ips[i%3]))
		}
	})

	t.Run("one ip blocked, but other two is working", func(t *testing.T) {

		r := NewRatelimiter(testConfig)
		for i := 0; i < 20; i++ {
			assert.Equal(t, true, r.IsLimitOK(ips[0]))
		}

		assert.Equal(t, false, r.IsLimitOK(ips[0]))
		assert.Equal(t, true, r.IsLimitOK(ips[1]))
		assert.Equal(t, true, r.IsLimitOK(ips[2]))
	})

	t.Run("all three ip blocked, but one reseted and working again", func(t *testing.T) {

		r := NewRatelimiter(testConfig)
		for i := 0; i < 60; i++ {
			assert.Equal(t, true, r.IsLimitOK(ips[i%3]))

		}

		assert.Equal(t, false, r.IsLimitOK(ips[0]))
		assert.Equal(t, false, r.IsLimitOK(ips[1]))
		assert.Equal(t, false, r.IsLimitOK(ips[2]))

		r.ResetLimit(ips[0].Mask(net.CIDRMask(testConfig.PrefixSize, 32)).String())

		assert.Equal(t, true, r.IsLimitOK(ips[0]))
		assert.Equal(t, true, r.IsLimitOK(ips[0]))
	})
}
