// Copyright 2026 Byterio
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memkey

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func checkTimeStamp(t testing.TB, expectedCurrent, actualCurrent uint32) {
	require.True(t, actualCurrent >= expectedCurrent-1 || actualCurrent <= expectedCurrent+1)
}

func Test_TimeStampUpdater(t *testing.T) {
	t.Parallel()
	startTimeStampUpdater()
	now := uint32(time.Now().Unix())
	checkTimeStamp(t, now, atomic.LoadUint32(&timestamp))
	time.Sleep(1 * time.Second)
	checkTimeStamp(t, now+1, atomic.LoadUint32(&timestamp))
	time.Sleep(1 * time.Second)
	checkTimeStamp(t, now+2, atomic.LoadUint32(&timestamp))
}

func Test_Memkey_Set(t *testing.T) {
	var (
		memkey = New()
		key    = "ping"
		val    = []byte("pong")
	)
	err := memkey.Set(key, val, 0)
	require.NoError(t, err)
	keys, err := memkey.Keys()
	size := memkey.Size()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, size, 1)
}

func Test_Memkey_Set_Override(t *testing.T) {
	var (
		memkey = New()
		key    = "ping"
		val    = []byte("pong")
	)
	err := memkey.Set(key, val, 0)
	require.NoError(t, err)
	err = memkey.Set(key, val, 0)
	require.NoError(t, err)
	keys, err := memkey.Keys()
	size := memkey.Size()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, size, 1)

}

func Test_Memkey_Get(t *testing.T) {
	var (
		memkey = New()
		key    = "ping"
		val    = []byte("pong")
	)
	err := memkey.Set(key, val, 0)
	require.NoError(t, err)
	result, err := memkey.Get(key)
	require.NoError(t, err)
	require.Equal(t, val, result)
	keys, err := memkey.Keys()
	size := memkey.Size()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, size, 1)
}

func Test_Memkey_Set_Expiration(t *testing.T) {
	var (
		memkey = New()
		key    = "ping"
		val    = []byte("pong")
		exp    = 1 * time.Second
	)
	err := memkey.Set(key, val, exp)
	require.NoError(t, err)
	time.Sleep(1100 * time.Millisecond)
	result, err := memkey.Get(key)
	require.NoError(t, err)
	require.Zero(t, len(result))
	keys, err := memkey.Keys()
	size := memkey.Size()
	require.NoError(t, err)
	require.Nil(t, keys)
	require.Equal(t, size, 0)
}

func Test_Memkey_Set_Long_Expiration_with_Keys(t *testing.T) {
	var (
		memkey = New()
		key    = "ping"
		val    = []byte("pong")
		exp    = 5 * time.Second
	)
	keys, err := memkey.Keys()
	size := memkey.Size()
	require.NoError(t, err)
	require.Nil(t, keys)
	require.Equal(t, size, 0)
	err = memkey.Set(key, val, exp)
	require.NoError(t, err)
	time.Sleep(1100 * time.Millisecond)
	keys, err = memkey.Keys()
	size = memkey.Size()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, size, 1)
	time.Sleep(4000 * time.Millisecond)
	result, err := memkey.Get(key)
	require.NoError(t, err)
	require.Zero(t, len(result))
	keys, err = memkey.Keys()
	size = memkey.Size()
	require.NoError(t, err)
	require.Nil(t, keys)
	require.Equal(t, size, 0)
}

func Test_Memkey_Get_NotExist(t *testing.T) {
	memkey := New()
	result, err := memkey.Get("notexist")
	require.NoError(t, err)
	require.Zero(t, len(result))
	keys, err := memkey.Keys()
	size := memkey.Size()
	require.NoError(t, err)
	require.Nil(t, keys)
	require.Equal(t, size, 0)
}

func Test_Memkey_Delete(t *testing.T) {
	var (
		memkey = New()
		key    = "ping"
		val    = []byte("pong")
	)
	err := memkey.Set(key, val, 0)
	require.NoError(t, err)
	keys, err := memkey.Keys()
	size := memkey.Size()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, size, 1)
	err = memkey.Delete(key)
	require.NoError(t, err)
	result, err := memkey.Get(key)
	require.NoError(t, err)
	require.Zero(t, len(result))
	keys, err = memkey.Keys()
	size = memkey.Size()
	require.NoError(t, err)
	require.Nil(t, keys)
	require.Equal(t, size, 0)
}

func Test_Memkey_Reset(t *testing.T) {
	memkey := New()
	val := []byte("pong")
	err := memkey.Set("ping1", val, 0)
	require.NoError(t, err)
	err = memkey.Set("ping2", val, 0)
	require.NoError(t, err)
	keys, err := memkey.Keys()
	size := memkey.Size()
	require.NoError(t, err)
	require.Len(t, keys, 2)
	require.Equal(t, size, 2)
	err = memkey.Reset()
	require.NoError(t, err)
	result, err := memkey.Get("ping1")
	require.NoError(t, err)
	require.Zero(t, len(result))
	result, err = memkey.Get("ping2")
	require.NoError(t, err)
	require.Zero(t, len(result))
	keys, err = memkey.Keys()
	size = memkey.Size()
	require.NoError(t, err)
	require.Nil(t, keys)
	require.Equal(t, size, 0)
}

func Test_Memkey_Close(t *testing.T) {
	memkey := New()
	require.Nil(t, memkey.Close())
}

func Test_Memkey_Has(t *testing.T) {
	var (
		memkey = New()
		key    = "ping"
		val    = []byte("pong")
	)
	err := memkey.Set(key, val, 0)
	require.NoError(t, err)
	result := memkey.Has(key)
	require.NoError(t, err)
	require.Equal(t, true, result)
	keys, err := memkey.Keys()
	size := memkey.Size()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, size, 1)
}

func Benchmark_Memkey_Set(b *testing.B) {
	memkey := New()
	b.ReportAllocs()
	b.ResetTimer()
	var err error
	for i := 0; i < b.N; i++ {
		err = memkey.Set("ping", []byte("pong"), 0)
	}
	require.NoError(b, err)
}

func Benchmark_Memkey_Get(b *testing.B) {
	memkey := New()
	err := memkey.Set("ping", []byte("pong"), 0)
	require.NoError(b, err)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = memkey.Get("ping")
	}
	require.NoError(b, err)
}

func Benchmark_Memkey_SetAndDelete(b *testing.B) {
	memkey := New()
	b.ReportAllocs()
	b.ResetTimer()
	var err error
	for i := 0; i < b.N; i++ {
		_ = memkey.Set("ping", []byte("pong"), 0)
		err = memkey.Delete("ping")
	}
	require.NoError(b, err)
}
