package ristretto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testStore = New()

func Test_Ristretto_Set(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	err := testStore.Set(key, val, 0)
	require.NoError(t, err)
}

func Test_Ristretto_Set_Override(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	err := testStore.Set(key, val, 0)
	require.NoError(t, err)

	err = testStore.Set(key, val, 0)
	require.NoError(t, err)
}

func Test_Ristretto_Get(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	// Create a channel to synchronize the test
	done := make(chan struct{})

	// Set the value in a goroutine
	go func() {
		err := testStore.Set(key, val, 0)
		require.NoError(t, err)

		// Send a value on the channel to signal that the value has been set
		done <- struct{}{}
	}()

	// Wait for the value to be set or a timeout to occur
	select {
	case <-done:
		// The value has been set, proceed with the test
		result, err := testStore.Get(key)
		require.NoError(t, err)
		require.Equal(t, val, result)
	case <-time.After(1 * time.Second):
		// The value was not set within 1 second, fail the test
		t.Errorf("timed out waiting for value to be set")
	}
}

func Test_Ristretto_Set_Expiration(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
		exp = 1 * time.Second
	)

	err := testStore.Set(key, val, exp)
	require.NoError(t, err)

	testStore.Reset()
}

func Test_Ristretto_Get_Expired(t *testing.T) {
	key := "john"

	result, err := testStore.Get(key)
	require.NoError(t, err)
	require.Zero(t, len(result))
}

func Test_Ristretto_Get_NotExist(t *testing.T) {
	result, err := testStore.Get("notexist")
	require.NoError(t, err)
	require.Zero(t, len(result))
}

func Test_Ristretto_Delete(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	// Create a channel to synchronize the test
	done := make(chan struct{})

	// Set the value in a goroutine
	go func() {
		err := testStore.Set(key, val, 0)
		require.NoError(t, err)

		// Send a value on the channel to signal that the value has been set
		done <- struct{}{}
	}()

	// Wait for the value to be set or a timeout to occur
	select {
	case <-done:
		// The value has been set, proceed with the test
		err := testStore.Delete(key)
		require.NoError(t, err)

		result, err := testStore.Get(key)
		require.NoError(t, err)
		require.Zero(t, len(result))
	case <-time.After(1 * time.Second):
		// The value was not set within 1 second, fail the test
		t.Errorf("timed out waiting for value to be set")
	}
}

func Test_Ristretto_Reset(t *testing.T) {
	val := []byte("doe")

	err := testStore.Set("john1", val, 0)
	require.NoError(t, err)

	err = testStore.Set("john2", val, 0)
	require.NoError(t, err)

	err = testStore.Reset()
	require.NoError(t, err)

	result, err := testStore.Get("john1")
	require.NoError(t, err)
	require.Zero(t, len(result))

	result, err = testStore.Get("john2")
	require.NoError(t, err)
	require.Zero(t, len(result))
}

func Test_Ristretto_Close(t *testing.T) {
	require.Nil(t, testStore.Close())
}

func Test_Ristretto_Conn(t *testing.T) {
	require.True(t, testStore.Conn() != nil)
}
