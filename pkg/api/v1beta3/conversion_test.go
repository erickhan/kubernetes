/*
Copyright 2014 The Kubernetes Authors All rights reserved.

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

package v1beta3_test

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	versioned "github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta3"
)

func TestNodeConversion(t *testing.T) {
	obj, err := versioned.Codec.Decode([]byte(`{"kind":"Minion","apiVersion":"v1beta3"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := obj.(*api.Node); !ok {
		t.Errorf("unexpected type: %#v", obj)
	}

	obj, err = versioned.Codec.Decode([]byte(`{"kind":"MinionList","apiVersion":"v1beta3"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := obj.(*api.NodeList); !ok {
		t.Errorf("unexpected type: %#v", obj)
	}

	obj = &api.Node{}
	if err := versioned.Codec.DecodeInto([]byte(`{"kind":"Minion","apiVersion":"v1beta3"}`), obj); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBadSecurityContextConversion(t *testing.T) {
	priv := false
	testCases := map[string]struct {
		c   *versioned.Container
		err string
	}{
		// this use case must use true for the container and false for the sc.  Otherwise the defaulter
		// will assume privileged was left undefined (since it is the default value) and copy the
		// sc setting upwards
		"mismatched privileged": {
			c: &versioned.Container{
				Privileged: true,
				SecurityContext: &versioned.SecurityContext{
					Privileged: &priv,
				},
			},
			err: "container privileged settings do not match security context settings, cannot convert",
		},
		"mismatched caps add": {
			c: &versioned.Container{
				Capabilities: versioned.Capabilities{
					Add: []versioned.Capability{"foo"},
				},
				SecurityContext: &versioned.SecurityContext{
					Capabilities: &versioned.Capabilities{
						Add: []versioned.Capability{"bar"},
					},
				},
			},
			err: "container capability settings do not match security context settings, cannot convert",
		},
		"mismatched caps drop": {
			c: &versioned.Container{
				Capabilities: versioned.Capabilities{
					Drop: []versioned.Capability{"foo"},
				},
				SecurityContext: &versioned.SecurityContext{
					Capabilities: &versioned.Capabilities{
						Drop: []versioned.Capability{"bar"},
					},
				},
			},
			err: "container capability settings do not match security context settings, cannot convert",
		},
	}

	for k, v := range testCases {
		got := api.Container{}
		err := api.Scheme.Convert(v.c, &got)
		if err == nil {
			t.Errorf("expected error for case %s but got none", k)
		} else {
			if err.Error() != v.err {
				t.Errorf("unexpected error for case %s.  Expected: %s but got: %s", k, v.err, err.Error())
			}
		}
	}

}
