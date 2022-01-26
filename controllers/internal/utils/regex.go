package utils

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"regexp"
	"strings"
)

// SplitAfter works as the same way as strings.SplitAfter() but the separator is regexp
func SplitAfter(s string, re *regexp.Regexp) (r []string) {
	if re == nil {
		return
	}
	re.ReplaceAllStringFunc(s, func(x string) string {
		s = strings.ReplaceAll(s, x, "::"+x)
		return s
	})
	for _, x := range strings.Split(s, "::") {
		if x != "" {
			r = append(r, x)
		}
	}
	return
}