package common

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"sync"
)

var restyInstanceMap = make(map[string]*resty.Client)
var restyLock sync.Mutex

// GetRestyClient 一个域名对应一个resty.Client
// https://github.com/go-resty/resty/issues/612
func GetRestyClient(url string) (*resty.Client, error) {
	base := ParseBaseUri(url)
	if base == "" {
		return nil, fmt.Errorf("invalid url: %s", url)
	}

	if client, ok := restyInstanceMap[base]; !ok {
		restyLock.Lock()
		defer restyLock.Unlock()

		if client, ok = restyInstanceMap[base]; ok {
			return client, nil
		}

		newClient := resty.New()
		restyInstanceMap[base] = newClient
		return newClient, nil
	} else {
		return client, nil
	}

}
