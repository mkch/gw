// Code generated by "stringer -output sys_string.go -type=HRESULT,RPC_STATUS"; DO NOT EDIT.

package sys

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[S_OK-0]
	_ = x[S_FALSE-1]
	_ = x[E_NOINTERFACE - -2147467262]
	_ = x[E_POINTER - -2147467261]
	_ = x[E_UNEXPECTED - -2147418113]
	_ = x[E_NOTIMPL - -2147467263]
	_ = x[E_OUTOFMEMORY - -2147024882]
	_ = x[E_INVALIDARG - -2147024809]
	_ = x[E_HANDLE - -2147024890]
	_ = x[E_ABORT - -2147467260]
	_ = x[E_FAIL - -2147467259]
	_ = x[E_ACCESSDENIED - -2147024891]
}

const (
	_HRESULT_name_0 = "E_NOTIMPLE_NOINTERFACEE_POINTERE_ABORTE_FAIL"
	_HRESULT_name_1 = "E_UNEXPECTED"
	_HRESULT_name_2 = "E_ACCESSDENIEDE_HANDLE"
	_HRESULT_name_3 = "E_OUTOFMEMORY"
	_HRESULT_name_4 = "E_INVALIDARG"
	_HRESULT_name_5 = "S_OKS_FALSE"
)

var (
	_HRESULT_index_0 = [...]uint8{0, 9, 22, 31, 38, 44}
	_HRESULT_index_2 = [...]uint8{0, 14, 22}
	_HRESULT_index_5 = [...]uint8{0, 4, 11}
)

func (i HRESULT) String() string {
	switch {
	case -2147467263 <= i && i <= -2147467259:
		i -= -2147467263
		return _HRESULT_name_0[_HRESULT_index_0[i]:_HRESULT_index_0[i+1]]
	case i == -2147418113:
		return _HRESULT_name_1
	case -2147024891 <= i && i <= -2147024890:
		i -= -2147024891
		return _HRESULT_name_2[_HRESULT_index_2[i]:_HRESULT_index_2[i+1]]
	case i == -2147024882:
		return _HRESULT_name_3
	case i == -2147024809:
		return _HRESULT_name_4
	case 0 <= i && i <= 1:
		return _HRESULT_name_5[_HRESULT_index_5[i]:_HRESULT_index_5[i+1]]
	default:
		return "HRESULT(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RPC_S_OK-0]
	_ = x[RPC_S_INVALID_STRING_UUID-1705]
}

const (
	_RPC_STATUS_name_0 = "RPC_S_OK"
	_RPC_STATUS_name_1 = "RPC_S_INVALID_STRING_UUID"
)

func (i RPC_STATUS) String() string {
	switch {
	case i == 0:
		return _RPC_STATUS_name_0
	case i == 1705:
		return _RPC_STATUS_name_1
	default:
		return "RPC_STATUS(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
