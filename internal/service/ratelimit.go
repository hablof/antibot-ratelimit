package service

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	bucketSize        = 100
	tokenRecoveryTime = time.Minute
	banDuration       = 2 * time.Minute
	maskSize          = 24
)

type unaryLimiter struct {
	limitingMode atomic.Bool
	tokenPool    chan struct{}
}

func (ul *unaryLimiter) isLimitOK() bool {
	if ul.limitingMode.Load() {
		return false
	}

	select {
	// если имеется token взяли token
	// в горутине вернули токен через tokenRecoveryTime
	// isRateLimitOK вернёт true
	case <-ul.tokenPool:
		go func() {
			time.Sleep(tokenRecoveryTime)

			select {
			case ul.tokenPool <- struct{}{}:

			default: // ветка необходима на случай вызова ul.resetLimit() в момент ожидания
				return
			}
		}()

		return true

	// если токенов нет, значит ratelimit превышен
	// переходим в limiting mode на время banDuration
	default:
		ul.limitingMode.Store(true)

		go func() {
			time.Sleep(banDuration)
			ul.limitingMode.Store(false)
		}()
	}

	return false
}

// заполняет bucket token'ами до отказа
func (ul *unaryLimiter) resetLimit() {
	for {
		select {
		case ul.tokenPool <- struct{}{}:

		default:
			return
		}
	}
}

func newUnaryLimiter() *unaryLimiter {
	ul := unaryLimiter{
		tokenPool:    make(chan struct{}, bucketSize),
		limitingMode: atomic.Bool{},
	}

	ul.resetLimit()

	return &ul
}

// Ratelimiter -- структура верхнего уровня по отношнию к unaryLimiter
// создаёт unaryLimiter'ы на префиксы ip, вызывает их методы
type Ratelimiter struct {
	mapMutex sync.Mutex
	limiters map[string]*unaryLimiter
	mask     net.IPMask
}

func NewRatelimiter() *Ratelimiter {
	return &Ratelimiter{
		limiters: map[string]*unaryLimiter{},
		mask:     net.CIDRMask(maskSize, 32),
	}
}

func (rl *Ratelimiter) IsLimitOK(ip net.IP) bool {
	prefix := ip.Mask(rl.mask).String()

	rl.mapMutex.Lock()
	ul, ok := rl.limiters[prefix]
	if !ok {
		ul = newUnaryLimiter()
		rl.limiters[prefix] = ul
	}
	rl.mapMutex.Unlock()

	return ul.isLimitOK()
}

func (rl *Ratelimiter) ResetLimit(ipWithPrefix net.IP) {
	prefix := ipWithPrefix.Mask(rl.mask).String()

	rl.mapMutex.Lock()
	ul, ok := rl.limiters[prefix]
	rl.mapMutex.Unlock()

	if !ok {
		return
	}

	ul.resetLimit()
}
