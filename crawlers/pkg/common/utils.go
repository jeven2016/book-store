package common

import (
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

var InvalidFileErr = errors.New("the path isn't a valid file")

// PrintCmdErr print error in console
func PrintCmdErr(err error) {
	_, err = fmt.Fprintf(os.Stderr, "Error: '%s' \n", err)
	if err != nil {
		panic(err)
	}
}

func IsFileExists(file string) (bool, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return false, err
	}
	if stat.IsDir() {
		return false, InvalidFileErr
	}
	return true, err
}

func GenExpireTime() time.Duration {
	min := 2
	max := 5
	return time.Duration(rand.Intn(max-min)+min) * time.Minute
}

func ParseBaseUri(url string) string {
	reg, err := regexp.Compile("(https?://[^/]*)")
	if err != nil {
		// print log
		return ""
	}
	subs := reg.FindStringSubmatch(url)

	if len(subs) > 1 {
		return subs[1]
	}
	return ""
}

func BuildUrl(baseUri string, path string) string {
	return strings.TrimSuffix(ParseBaseUri(baseUri), "/") + "/" + strings.TrimPrefix(path, "/")
}

func GetSiteConfig(siteKey string) *SiteConfig {
	cfg, ok := slice.FindBy(GetConfig().WebSites, func(index int, item SiteConfig) bool {
		return item.Name == siteKey
	})
	if !ok {
		return nil
	}
	return &cfg
}
