package main

/*
#define QOI_IMPLEMENTATION
#include "./qoi/extra.h"
*/
import "C"

// Representation of the QOI header
type QOIDesc struct {
	// width of the image in pixels (BE)
	width uint32
	// height of the image in pixels (BE)
	height uint32
	// number of channels (3 = RGB, 4 = RGBA)
	channels uint8
	// colorspace of image (0 = sRGB with linear alpha, 1 = all channels are linear)
	colorspace uint8
}

func QoiEncode(data []byte, desc QOIDesc) ([]byte, int64) {
	cData := C.CBytes(data)
	defer C.free(cData)
	cDesc := C.qoi_desc{
		width:      C.uint(desc.width),
		height:     C.uint(desc.height),
		channels:   C.uchar(desc.channels),
		colorspace: C.uchar(desc.colorspace),
	}
	var cLen C.int
	cEncoded := C.qoi_encode(cData, &cDesc, &cLen)
	chunks := C.GoBytes(cEncoded, cLen)
	if cEncoded == nil {
		return nil, -1
	}
	return chunks, int64(cLen)
}
func QoiEncodeDiffLuma(data []byte, desc QOIDesc) ([]byte, int64) {
	cData := C.CBytes(data)
	defer C.free(cData)
	cDesc := C.qoi_desc{
		width:      C.uint(desc.width),
		height:     C.uint(desc.height),
		channels:   C.uchar(desc.channels),
		colorspace: C.uchar(desc.colorspace),
	}
	var cLen C.int

	cEncoded := C.qoi_encode_diff_luma(cData, &cDesc, &cLen)
	chunks := C.GoBytes(cEncoded, cLen)
	if cEncoded == nil {
		return nil, -1
	}
	return chunks, int64(cLen)
}
func QoiEncodeRun(data []byte, desc QOIDesc) ([]byte, int64) {
	cData := C.CBytes(data)
	defer C.free(cData)
	cDesc := C.qoi_desc{
		width:      C.uint(desc.width),
		height:     C.uint(desc.height),
		channels:   C.uchar(desc.channels),
		colorspace: C.uchar(desc.colorspace),
	}
	var cLen C.int

	cEncoded := C.qoi_encode_run(cData, &cDesc, &cLen)
	chunks := C.GoBytes(cEncoded, cLen)
	if cEncoded == nil {
		return nil, -1
	}
	return chunks, int64(cLen)
}