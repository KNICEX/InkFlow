package backoff

import "time"

type Policy struct {
	MaxRetries         int
	InitialInterval    time.Duration
	MaxRetryInterval   time.Duration
	BackoffCoefficient float64
	OnRetry            func(int, error)
}

var DefaultPolicy = Policy{
	MaxRetries:         5,
	InitialInterval:    1 * time.Second,
	MaxRetryInterval:   time.Minute,
	BackoffCoefficient: 2.0,
	OnRetry:            nil,
}

func setDefaultPolicy(policy Policy) Policy {
	if policy.MaxRetries == 0 {
		policy.MaxRetries = 5
	}
	if policy.InitialInterval == 0 {
		policy.InitialInterval = 1 * time.Second
	}
	if policy.MaxRetryInterval == 0 {
		policy.MaxRetryInterval = time.Minute
	}
	if policy.BackoffCoefficient == 0 {
		policy.BackoffCoefficient = 2.0
	}
	return policy
}

func Wrap(fn func() error, policy Policy) func() error {
	policy = setDefaultPolicy(policy)
	return func() error {
		var err error
		for i := 0; i < policy.MaxRetries; i++ {
			err = fn()
			if err == nil {
				return nil
			}
			if policy.OnRetry != nil {
				policy.OnRetry(i, err)
			}
			time.Sleep(min(
				time.Duration(int64(policy.BackoffCoefficient*float64(i)*float64(policy.InitialInterval))),
				policy.MaxRetryInterval),
			)
		}
		return err
	}
}
