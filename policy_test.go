package fastimage

import (
	"testing"
)

func Test_NewPolicy(t *testing.T) {
	policy := NewPolicy()
	if policy == nil {
		t.Error("NewPolicy() returned nil")
	}
}

func Test_Policy_WithImageTypes(t *testing.T) {
	policy := NewPolicy().WithImageTypes(JPEG, PNG)
	if len(policy.imageTypes) != 2 {
		t.Errorf("Expected 2 image types, got %d", len(policy.imageTypes))
	}
}

func Test_Policy_WithMaxImageFileSizeInBytes(t *testing.T) {
	policy := NewPolicy().WithMaxImageFileSizeInBytes(10000)
	if policy.maxImageFileBytesSize != 10000 {
		t.Errorf("Expected maxImageFileBytesSize to be 10000, got %d", policy.maxImageFileBytesSize)
	}
}

func Test_Policy_WithMinImageFileSizeInBytes(t *testing.T) {
	policy := NewPolicy().WithMinImageFileSizeInBytes(100)
	if policy.minImageFileBytesSize != 100 {
		t.Errorf("Expected minImageFileBytesSize to be 100, got %d", policy.minImageFileBytesSize)
	}
}

func Test_Policy_WithFetchFullImage(t *testing.T) {
	policy := NewPolicy().WithFetchFullImage(true)
	if !policy.fetchFullImage {
		t.Error("Expected fetchFullImage to be true, got false")
	}
}
