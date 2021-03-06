// Copyright (c) 2015-2018, NVIDIA CORPORATION. All rights reserved.

// +build linux darwin

package nvml

import (
	"unsafe"
)

/*
#include <dlfcn.h>
#include "nvml.h"

// We wrap the call to nvmlInit() here to ensure that we pick up the correct
// version of this call. The macro magic in nvml.h that #defines the symbol
// 'nvmlInit' to 'nvmlInit_v2' is unfortunately lost on cgo.
static nvmlReturn_t nvmlInit_dl(void) {
	return nvmlInit();
}
*/
import "C"

type dlhandles struct{ handles []unsafe.Pointer }

var dl dlhandles

// Initialize NVML, opening a dynamic reference to the NVML library in the process.
func (dl *dlhandles) nvmlInit() C.nvmlReturn_t {
	handle := C.dlopen(C.CString("libnvidia-ml.so.1"), C.RTLD_LAZY|C.RTLD_GLOBAL)
	if handle == C.NULL {
		return C.NVML_ERROR_LIBRARY_NOT_FOUND
	}
	dl.handles = append(dl.handles, handle)
	return C.nvmlInit_dl()
}

// Shutdown NVML, closing our dynamic reference to the NVML library in the process.
func (dl *dlhandles) nvmlShutdown() C.nvmlReturn_t {
	ret := C.nvmlShutdown()
	if ret != C.NVML_SUCCESS {
		return ret
	}

	for _, handle := range dl.handles {
		err := C.dlclose(handle)
		if err != 0 {
			return C.NVML_ERROR_UNKNOWN
		}
	}

	return C.NVML_SUCCESS
}
