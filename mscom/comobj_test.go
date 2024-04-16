package mscom_test

import (
	"fmt"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/mscom"
	"github.com/mkch/gw/mscom/sys"
)

type INetworkListManagerEventsImplVMT struct {
	mscom.IUnknownVMT

	connectivityChanged mscom.MethodPtr
}

type INetworkListManagerEventsImpl struct {
	vt *INetworkListManagerEventsImplVMT
}

func (s *INetworkListManagerEventsImpl) IUnknown() *mscom.IUnknown {
	return (*mscom.IUnknown)(unsafe.Pointer(s))
}

var IID_INetworkListManagerEvents = gg.Must(sys.UuidFromStringW("dcb00001-570f-4a9b-8d69-199fdba5723b"))

func ExampleInitIUnknownImpl() {
	var i *INetworkListManagerEventsImpl
	if err := CreateINetworkListManagerEventsImpl(&i); err != nil {
		panic(err)
	}
	defer i.IUnknown().Release()
}

func CreateINetworkListManagerEventsImpl(ppObject **INetworkListManagerEventsImpl) error {
	if ppObject == nil {
		return sys.HResultError(sys.E_POINTER)
	}
	// Alloc the interface and v-table in one block of memory.
	mem := mscom.Alloc[struct {
		INetworkListManagerEventsImpl
		INetworkListManagerEventsImplVMT
	}]()
	// Setup v-table.
	mem.INetworkListManagerEventsImpl.vt = &mem.INetworkListManagerEventsImplVMT
	// Call InitIUnknownImpl to initialize the IUnknownVMT part.
	mscom.InitIUnknownImpl(&mem.INetworkListManagerEventsImpl, &mem.IUnknownVMT, func(id sys.REFIID, p *unsafe.Pointer) sys.HRESULT {
		// This object can be converted to INetworkListManagerEvents interface.
		if *id == *IID_INetworkListManagerEvents {
			*p = unsafe.Pointer(&mem.INetworkListManagerEventsImpl)
			return sys.S_OK
		}
		return sys.E_NOINTERFACE
	}, func() {
		mscom.Free(mem)
	}).
		// Create the ConnectivityChanged method of INetworkListManagerEvents.
		Create(&mem.INetworkListManagerEventsImplVMT.connectivityChanged, func(connectivity uintptr) uintptr {
			fmt.Printf("connectivityChanged: %v\n", connectivity)
			return uintptr(sys.S_OK)
		})

	*ppObject = &mem.INetworkListManagerEventsImpl
	return nil
}
