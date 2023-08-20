package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testStore = New(Config{
	Database: os.Getenv("MYSQL_DATABASE"),
	Username: os.Getenv("MYSQL_USERNAME"),
	Password: os.Getenv("MYSQL_PASSWORD"),
	Reset:    true,
})

func Test_MYSQL_New(t *testing.T) {
	newConfigStore := New(Config{
		Database: os.Getenv("MYSQL_DATABASE"),
		Username: os.Getenv("MYSQL_USERNAME"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Reset:    true,
	})

	require.True(t, newConfigStore.db != nil)
	newConfigStore.Close()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", os.Getenv("MYSQL_USERNAME"), os.Getenv("MYSQL_PASSWORD"), "127.0.0.1", 3306, os.Getenv("MYSQL_DATABASE"))
	newConfigStore = New(Config{
		ConnectionURI: dsn,
		Reset:         true,
	})

	require.True(t, newConfigStore.db != nil)
	newConfigStore.Close()

	db, _ := sql.Open("mysql", dsn)
	newConfigStore = New(Config{
		Db:    db,
		Reset: true,
	})

	require.True(t, newConfigStore.db != nil)
	newConfigStore.Close()
}

func Test_MYSQL_Set(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	err := testStore.Set(key, val, 0)
	require.Nil(t, err)
}

func Test_MYSQL_Set_Override(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	err := testStore.Set(key, val, 0)
	require.Nil(t, err)

	err = testStore.Set(key, val, 0)
	require.Nil(t, err)
}

func Test_MYSQL_Get(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	err := testStore.Set(key, val, 0)
	require.Nil(t, err)

	result, err := testStore.Get(key)
	require.Nil(t, err)
	require.Equal(t, val, result)
}

func Test_MYSQL_Set_Expiration(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
		exp = 1 * time.Second
	)

	err := testStore.Set(key, val, exp)
	require.Nil(t, err)

	time.Sleep(1100 * time.Millisecond)
}

func Test_MYSQL_Get_Expired(t *testing.T) {
	key := "john"

	result, err := testStore.Get(key)
	require.Nil(t, err)
	require.Zero(t, len(result))
}

func Test_MYSQL_Get_NotExist(t *testing.T) {
	result, err := testStore.Get("notexist")
	require.Nil(t, err)
	require.Zero(t, len(result))
}

func Test_MYSQL_Delete(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	err := testStore.Set(key, val, 0)
	require.Nil(t, err)

	err = testStore.Delete(key)
	require.Nil(t, err)

	result, err := testStore.Get(key)
	require.Nil(t, err)
	require.Zero(t, len(result))
}

func Test_MYSQL_Reset(t *testing.T) {
	val := []byte("doe")

	err := testStore.Set("john1", val, 0)
	require.Nil(t, err)

	err = testStore.Set("john2", val, 0)
	require.Nil(t, err)

	err = testStore.Reset()
	require.Nil(t, err)

	result, err := testStore.Get("john1")
	require.Nil(t, err)
	require.Zero(t, len(result))

	result, err = testStore.Get("john2")
	require.Nil(t, err)
	require.Zero(t, len(result))
}

func Test_MYSQL_GC(t *testing.T) {
	testVal := []byte("doe")

	// This key should expire
	err := testStore.Set("john", testVal, time.Nanosecond)
	require.Nil(t, err)

	testStore.gc(time.Now())
	row := testStore.db.QueryRow(testStore.sqlSelect, "john")
	err = row.Scan(nil, nil)
	require.Equal(t, sql.ErrNoRows, err)

	// This key should not expire
	err = testStore.Set("john", testVal, 0)
	require.Nil(t, err)

	testStore.gc(time.Now())
	val, err := testStore.Get("john")
	require.Nil(t, err)
	require.Equal(t, testVal, val)
}

func Test_MYSQL_Non_UTF8(t *testing.T) {
	val := []byte("0xF5")

	err := testStore.Set("0xF6", val, 0)
	require.Nil(t, err)

	result, err := testStore.Get("0xF6")
	require.Nil(t, err)
	require.Equal(t, val, result)
}

func Test_MYSQL_Close(t *testing.T) {
	require.Nil(t, testStore.Close())
}

func Test_MYSQL_Conn(t *testing.T) {
	require.True(t, testStore.Conn() != nil)
}
