/*
Copyright 2026.

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

package controller

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	imagesv1alpha1 "github.com/saphire/image-updater-operator/api/v1alpha1"
)

func TestScanDue(t *testing.T) {
	const interval = time.Minute

	newIP := func(gen, observed int64, last *time.Time) *imagesv1alpha1.ImagePolicy {
		ip := &imagesv1alpha1.ImagePolicy{}
		ip.Generation = gen
		ip.Status.ObservedGeneration = observed
		if last != nil {
			t := metav1.NewTime(*last)
			ip.Status.LastScanTime = &t
		}
		return ip
	}
	ago := func(d time.Duration) *time.Time { t := time.Now().Add(-d); return &t }

	t.Run("first reconcile (never scanned) is due", func(t *testing.T) {
		due, _ := scanDue(newIP(1, 0, nil), interval)
		if !due {
			t.Fatal("want due on first reconcile")
		}
	})

	t.Run("spec change (generation ahead of observed) is due", func(t *testing.T) {
		due, _ := scanDue(newIP(2, 1, ago(2*time.Second)), interval)
		if !due {
			t.Fatal("want due after a spec change even within the interval")
		}
	})

	t.Run("scanned recently is not due and reports remaining", func(t *testing.T) {
		due, after := scanDue(newIP(1, 1, ago(20*time.Second)), interval)
		if due {
			t.Fatal("want not due 20s into a 60s interval")
		}
		if after <= 0 || after > 41*time.Second {
			t.Fatalf("remaining = %v, want roughly 40s", after)
		}
	})

	t.Run("interval elapsed is due", func(t *testing.T) {
		due, _ := scanDue(newIP(1, 1, ago(90*time.Second)), interval)
		if !due {
			t.Fatal("want due after the interval has elapsed")
		}
	})
}
