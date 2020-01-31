/*
Copyright 2020 Cornelius Weig.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"context"
	"testing"
)

func TestHomebrew_Fetch(t *testing.T) {
	homebrew := NewHomebrew("https://formulae.brew.sh/api/formula/krew.json")

	stats, err := homebrew.FetchAnalytics(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if stats.Installs365d == 0 {
		t.Errorf("install count over 1y should not be zero")
	}
}
