package pgtype_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/testutil"
)

func TestInetArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "inet[]", []interface{}{
		&pgtype.InetArray{
			Elements:   nil,
			Dimensions: nil,
			Valid:      true,
		},
		&pgtype.InetArray{
			Elements: []pgtype.Inet{
				{IPNet: mustParseCIDR(t, "12.34.56.0/32"), Valid: true},
				{},
			},
			Dimensions: []pgtype.ArrayDimension{{Length: 2, LowerBound: 1}},
			Valid:      true,
		},
		&pgtype.InetArray{},
		&pgtype.InetArray{
			Elements: []pgtype.Inet{
				{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Valid: true},
				{IPNet: mustParseCIDR(t, "12.34.56.0/32"), Valid: true},
				{IPNet: mustParseCIDR(t, "192.168.0.1/32"), Valid: true},
				{IPNet: mustParseCIDR(t, "2607:f8b0:4009:80b::200e/128"), Valid: true},
				{},
				{IPNet: mustParseCIDR(t, "255.0.0.0/8"), Valid: true},
			},
			Dimensions: []pgtype.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Valid:      true,
		},
		&pgtype.InetArray{
			Elements: []pgtype.Inet{
				{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Valid: true},
				{IPNet: mustParseCIDR(t, "12.34.56.0/32"), Valid: true},
				{IPNet: mustParseCIDR(t, "192.168.0.1/32"), Valid: true},
				{IPNet: mustParseCIDR(t, "2607:f8b0:4009:80b::200e/128"), Valid: true},
			},
			Dimensions: []pgtype.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Valid: true,
		},
	})
}

func TestInetArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result pgtype.InetArray
	}{
		{
			source: []*net.IPNet{mustParseCIDR(t, "127.0.0.1/32")},
			result: pgtype.InetArray{
				Elements:   []pgtype.Inet{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Valid: true}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Valid:      true},
		},
		{
			source: (([]*net.IPNet)(nil)),
			result: pgtype.InetArray{},
		},
		{
			source: []net.IP{mustParseCIDR(t, "127.0.0.1/32").IP},
			result: pgtype.InetArray{
				Elements:   []pgtype.Inet{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Valid: true}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Valid:      true},
		},
		{
			source: (([]net.IP)(nil)),
			result: pgtype.InetArray{},
		},
		{
			source: [][]net.IP{{mustParseCIDR(t, "127.0.0.1/32").IP}, {mustParseCIDR(t, "10.0.0.1/32").IP}},
			result: pgtype.InetArray{
				Elements: []pgtype.Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Valid: true},
					{IPNet: mustParseCIDR(t, "10.0.0.1/32"), Valid: true}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Valid:      true},
		},
		{
			source: [][][][]*net.IPNet{
				{{{
					mustParseCIDR(t, "127.0.0.1/24"),
					mustParseCIDR(t, "10.0.0.1/24"),
					mustParseCIDR(t, "172.16.0.1/16")}}},
				{{{
					mustParseCIDR(t, "192.168.0.1/16"),
					mustParseCIDR(t, "224.0.0.1/24"),
					mustParseCIDR(t, "169.168.0.1/16")}}}},
			result: pgtype.InetArray{
				Elements: []pgtype.Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/24"), Valid: true},
					{IPNet: mustParseCIDR(t, "10.0.0.1/24"), Valid: true},
					{IPNet: mustParseCIDR(t, "172.16.0.1/16"), Valid: true},
					{IPNet: mustParseCIDR(t, "192.168.0.1/16"), Valid: true},
					{IPNet: mustParseCIDR(t, "224.0.0.1/24"), Valid: true},
					{IPNet: mustParseCIDR(t, "169.168.0.1/16"), Valid: true}},
				Dimensions: []pgtype.ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Valid: true},
		},
		{
			source: [2][1]net.IP{{mustParseCIDR(t, "127.0.0.1/32").IP}, {mustParseCIDR(t, "10.0.0.1/32").IP}},
			result: pgtype.InetArray{
				Elements: []pgtype.Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Valid: true},
					{IPNet: mustParseCIDR(t, "10.0.0.1/32"), Valid: true}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Valid:      true},
		},
		{
			source: [2][1][1][3]*net.IPNet{
				{{{
					mustParseCIDR(t, "127.0.0.1/24"),
					mustParseCIDR(t, "10.0.0.1/24"),
					mustParseCIDR(t, "172.16.0.1/16")}}},
				{{{
					mustParseCIDR(t, "192.168.0.1/16"),
					mustParseCIDR(t, "224.0.0.1/24"),
					mustParseCIDR(t, "169.168.0.1/16")}}}},
			result: pgtype.InetArray{
				Elements: []pgtype.Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/24"), Valid: true},
					{IPNet: mustParseCIDR(t, "10.0.0.1/24"), Valid: true},
					{IPNet: mustParseCIDR(t, "172.16.0.1/16"), Valid: true},
					{IPNet: mustParseCIDR(t, "192.168.0.1/16"), Valid: true},
					{IPNet: mustParseCIDR(t, "224.0.0.1/24"), Valid: true},
					{IPNet: mustParseCIDR(t, "169.168.0.1/16"), Valid: true}},
				Dimensions: []pgtype.ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Valid: true},
		},
	}

	for i, tt := range successfulTests {
		var r pgtype.InetArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestInetArrayAssignTo(t *testing.T) {
	var ipnetSlice []*net.IPNet
	var ipSlice []net.IP
	var ipSliceDim2 [][]net.IP
	var ipnetSliceDim4 [][][][]*net.IPNet
	var ipArrayDim2 [2][1]net.IP
	var ipnetArrayDim4 [2][1][1][3]*net.IPNet

	simpleTests := []struct {
		src      pgtype.InetArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: pgtype.InetArray{
				Elements:   []pgtype.Inet{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Valid: true}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Valid:      true,
			},
			dst:      &ipnetSlice,
			expected: []*net.IPNet{mustParseCIDR(t, "127.0.0.1/32")},
		},
		{
			src: pgtype.InetArray{
				Elements:   []pgtype.Inet{{}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Valid:      true,
			},
			dst:      &ipnetSlice,
			expected: []*net.IPNet{nil},
		},
		{
			src: pgtype.InetArray{
				Elements:   []pgtype.Inet{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Valid: true}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Valid:      true,
			},
			dst:      &ipSlice,
			expected: []net.IP{mustParseCIDR(t, "127.0.0.1/32").IP},
		},
		{
			src: pgtype.InetArray{
				Elements:   []pgtype.Inet{{}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Valid:      true,
			},
			dst:      &ipSlice,
			expected: []net.IP{nil},
		},
		{
			src:      pgtype.InetArray{},
			dst:      &ipnetSlice,
			expected: (([]*net.IPNet)(nil)),
		},
		{
			src:      pgtype.InetArray{Valid: true},
			dst:      &ipnetSlice,
			expected: []*net.IPNet{},
		},
		{
			src:      pgtype.InetArray{},
			dst:      &ipSlice,
			expected: (([]net.IP)(nil)),
		},
		{
			src:      pgtype.InetArray{Valid: true},
			dst:      &ipSlice,
			expected: []net.IP{},
		},
		{
			src: pgtype.InetArray{
				Elements: []pgtype.Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Valid: true},
					{IPNet: mustParseCIDR(t, "10.0.0.1/32"), Valid: true}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Valid:      true},
			dst:      &ipSliceDim2,
			expected: [][]net.IP{{mustParseCIDR(t, "127.0.0.1/32").IP}, {mustParseCIDR(t, "10.0.0.1/32").IP}},
		},
		{
			src: pgtype.InetArray{
				Elements: []pgtype.Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/24"), Valid: true},
					{IPNet: mustParseCIDR(t, "10.0.0.1/24"), Valid: true},
					{IPNet: mustParseCIDR(t, "172.16.0.1/16"), Valid: true},
					{IPNet: mustParseCIDR(t, "192.168.0.1/16"), Valid: true},
					{IPNet: mustParseCIDR(t, "224.0.0.1/24"), Valid: true},
					{IPNet: mustParseCIDR(t, "169.168.0.1/16"), Valid: true}},
				Dimensions: []pgtype.ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Valid: true},
			dst: &ipnetSliceDim4,
			expected: [][][][]*net.IPNet{
				{{{
					mustParseCIDR(t, "127.0.0.1/24"),
					mustParseCIDR(t, "10.0.0.1/24"),
					mustParseCIDR(t, "172.16.0.1/16")}}},
				{{{
					mustParseCIDR(t, "192.168.0.1/16"),
					mustParseCIDR(t, "224.0.0.1/24"),
					mustParseCIDR(t, "169.168.0.1/16")}}}},
		},
		{
			src: pgtype.InetArray{
				Elements: []pgtype.Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Valid: true},
					{IPNet: mustParseCIDR(t, "10.0.0.1/32"), Valid: true}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Valid:      true},
			dst:      &ipArrayDim2,
			expected: [2][1]net.IP{{mustParseCIDR(t, "127.0.0.1/32").IP}, {mustParseCIDR(t, "10.0.0.1/32").IP}},
		},
		{
			src: pgtype.InetArray{
				Elements: []pgtype.Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/24"), Valid: true},
					{IPNet: mustParseCIDR(t, "10.0.0.1/24"), Valid: true},
					{IPNet: mustParseCIDR(t, "172.16.0.1/16"), Valid: true},
					{IPNet: mustParseCIDR(t, "192.168.0.1/16"), Valid: true},
					{IPNet: mustParseCIDR(t, "224.0.0.1/24"), Valid: true},
					{IPNet: mustParseCIDR(t, "169.168.0.1/16"), Valid: true}},
				Dimensions: []pgtype.ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Valid: true},
			dst: &ipnetArrayDim4,
			expected: [2][1][1][3]*net.IPNet{
				{{{
					mustParseCIDR(t, "127.0.0.1/24"),
					mustParseCIDR(t, "10.0.0.1/24"),
					mustParseCIDR(t, "172.16.0.1/16")}}},
				{{{
					mustParseCIDR(t, "192.168.0.1/16"),
					mustParseCIDR(t, "224.0.0.1/24"),
					mustParseCIDR(t, "169.168.0.1/16")}}}},
		},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}
}
