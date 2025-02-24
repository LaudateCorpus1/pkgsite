// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fetch

import (
	"context"
	"testing"

	"golang.org/x/pkgsite/internal/derrors"
	"golang.org/x/pkgsite/internal/proxy/proxytest"
	"golang.org/x/pkgsite/internal/stdlib"
)

func TestLatestModuleVersions(t *testing.T) {
	// latestVersion is tested above.
	// Contents of the go.mod file are tested in proxydatasource.
	// Here, test retractions and presence of a go.mod file.
	prox, teardown := proxytest.SetupTestClient(t, testModules)
	defer teardown()

	stdlib.UseTestData = true
	defer func() { stdlib.UseTestData = false }()

	// These tests (except for std) depend on the test modules, which are taken from the contents
	// of internal/proxy/testdata/*.txtar.
	for _, test := range []struct {
		modulePath          string
		wantRaw, wantCooked string
	}{
		{"example.com/basic", "v1.1.0", "v1.1.0"},
		{"example.com/retractions", "v1.2.0", "v1.0.0"},
		{"std", "v1.14.6", "v1.14.6"},
	} {
		got, err := LatestModuleVersions(context.Background(), test.modulePath, prox, nil)
		if err != nil {
			t.Fatal(err)
		}
		if got.GoModFile == nil {
			t.Errorf("%s: no go.mod file", test.modulePath)
		}
		if got.RawVersion != test.wantRaw {
			t.Errorf("%s, raw: got %q, want %q", test.modulePath, got.RawVersion, test.wantRaw)
		}
		if got.CookedVersion != test.wantCooked {
			t.Errorf("%s, cooked: got %q, want %q", test.modulePath, got.CookedVersion, test.wantCooked)
		}
	}
}

func TestLatestModuleVersionsNotFound(t *testing.T) {
	// Verify that we get (nil, nil) if there is no version information.
	const modulePath = "example.com/no-versions"
	server := proxytest.NewServer(testModules)
	server.AddModuleNoVersions(&proxytest.Module{
		ModulePath: modulePath,
		Version:    "v0.0.0-20181107005212-dafb9c8d8707",
	})
	client, teardown, err := proxytest.NewClientForServer(server)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	got, err := LatestModuleVersions(context.Background(), modulePath, client, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestLatestModuleVersionsBadGoMod(t *testing.T) {
	// Verify that we get a BadModule error if the go.mod file is bad.
	const modulePath = "example.com/bad-go-mod"
	server := proxytest.NewServer([]*proxytest.Module{
		{
			ModulePath: modulePath,
			Version:    "v1.0.0",
			Files: map[string]string{
				"go.mod": "module example.com/bad-go-mod\ngo bad",
			},
		},
	})
	client, teardown, err := proxytest.NewClientForServer(server)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	_, err = LatestModuleVersions(context.Background(), modulePath, client, nil)
	if got, want := derrors.ToStatus(err), 490; got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}
