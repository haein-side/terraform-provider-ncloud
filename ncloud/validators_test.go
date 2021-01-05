package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func Test_validateInstanceName(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "tEsting123",
			ErrCount: 1,
		},
		{
			Value:    "testing123!",
			ErrCount: 1,
		},
		{
			Value:    "1testing123",
			ErrCount: 1,
		},
		{
			Value:    "한글-123",
			ErrCount: 1,
		},
		{
			Value:    "te",
			ErrCount: 1,
		},
		{
			Value:    "testing",
			ErrCount: 0,
		},
		{
			Value:    "testing-123",
			ErrCount: 0,
		},
		{
			Value:    "testing--123",
			ErrCount: 0,
		},
		{
			Value:    "testing_123",
			ErrCount: 1,
		},
		{
			Value:    "testing123-",
			ErrCount: 1,
		},
		{
			Value:    "testing123*",
			ErrCount: 1,
		},
		{
			Value:    acctest.RandStringFromCharSet(256, acctest.CharSetAlpha),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateInstanceName(tc.Value, "name")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Instance Name to trigger a validation error for %q", tc.Value)
		}
	}
}

func Test_validatePortRange(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "22",
			ErrCount: 0,
		},
		{
			Value:    "65535",
			ErrCount: 0,
		},
		{
			Value:    "1-65535",
			ErrCount: 0,
		},
		{
			Value:    "1-65536",
			ErrCount: 1,
		},
		{
			Value:    "0-65535",
			ErrCount: 1,
		},
		{
			Value:    "8081-22",
			ErrCount: 1,
		},
		{
			Value:    "65536",
			ErrCount: 1,
		},
		{
			Value:    "a22",
			ErrCount: 1,
		},
		{
			Value:    "a22-8081",
			ErrCount: 1,
		},
		{
			Value:    "22-33-567",
			ErrCount: 1,
		},
		{
			Value:    "22-",
			ErrCount: 1,
		},
		{
			Value:    "-22",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validatePortRange(tc.Value, "portRange")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the portRange to trigger a validation error for %q", tc.Value)
		}
	}
}
