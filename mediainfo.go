package mediainfo

/*
#cgo LDFLAGS: -ldl
#include <stdlib.h>
#include "MediaInfoDLL.h"

static inline void* g_MediaInfo_New() {
    void *mi = MediaInfo_New();
    if (mi)
        MediaInfo_Option(mi, "CharSet", "UTF-8");
    return mi;
}

static inline void g_MediaInfo_Delete(void *mi) {
    return MediaInfo_Delete(mi);
}

static inline size_t g_MediaInfo_Open(void *mi, const char *name) {
    return MediaInfo_Open(mi, name);
}

static inline const char* g_MediaInfo_Option(void *mi, const char *option, const char *value) {
    return MediaInfo_Option(mi, option, value);
}

static inline const char* g_MediaInfo_Inform(void *mi) {
    return MediaInfo_Inform(mi, 0);
}

static inline const char* g_MediaInfo_Get(void *mi, MediaInfo_stream_C streamKind, size_t streamNumber,
    const char *parameter, MediaInfo_info_C kindOfInfo, MediaInfo_info_C kindOfSearch) {
    return MediaInfo_Get(mi, streamKind, streamNumber, parameter, kindOfInfo, kindOfSearch);
}

static inline const char* g_MediaInfo_GetI(void *mi, MediaInfo_stream_C streamKind, size_t streamNumber,
    size_t parameter, MediaInfo_info_C kindOfInfo) {
    return MediaInfo_GetI(mi, streamKind, streamNumber, parameter, kindOfInfo);
}

static inline size_t g_MediaInfo_Count_Get(void *mi, MediaInfo_stream_C streamKind) {
    return MediaInfo_Count_Get(mi, streamKind, -1);
}

static inline void g_MediaInfo_Close(void *mi) {
    MediaInfo_Close(mi);
}

*/
import "C"

import (
	"unsafe"
	"runtime"
	"errors"
)

// MediaInfo is an instance of a mediainfo accessor.
type MediaInfo struct {
	cptr unsafe.Pointer
}

// StreamKind is used to specify the type of stream (audio, video, chapters, etc) when getting information.
type StreamKind int
// InfoKind is used to specify the aspect of information (name, value, unit of measure) when retrieving information.
type InfoKind int

const (
	// StreamGeneral is a for general container stream.
	StreamGeneral StreamKind = 0
	// StreamVideo is a video stream.
	StreamVideo = 1
	// StreamAudio is an audio stream.
	StreamAudio = 2
	// StreamText is embedded text (subtitles) stream.
	StreamText = 3
	// StreamOther is for chapters.
	StreamOther = 4
	// StreamImage is for embedded images.
	StreamImage = 5
	// StreamMenu is for dynamic menus.
	StreamMenu = 6
)

const (
	// InfoName is the unique name of parameter.
	InfoName InfoKind = 0
	// InfoText is value of parameter.
	InfoText = 1
	// InfoMeasure is the unique name of measure unit of parameter.
	InfoMeasure = 2
	InfoOptions = 3
	// InfoNameText is translated name of parameter.
	InfoNameText = 4
	// InfoMeasureText is translated name of measure unit.
	InfoMeasureText = 5
	// InfoInfo is more information about the parameter.
	InfoInfo = 6
	// InfoHowTo is how this parameter is supported, could be N (No), B (Beta), R (Read only), W (Read/write).
	InfoHowTo = 7
)


func toCInfo(i InfoKind) C.MediaInfo_info_C {
	return C.MediaInfo_info_C(i)
}

func toCStream(s StreamKind) C.MediaInfo_stream_C {
	return C.MediaInfo_stream_C(s)
}

var cEmptyString int = 0

func emptyCString() *_Ctype_char {
	return (*_Ctype_char)(unsafe.Pointer(&cEmptyString))
}

// ErrOpenFailed is returned by Open when mediainfo cannot open the file.
var ErrOpenFailed = errors.New("file open failed")

func init() {
	C.MediaInfoDLL_Load()
}

// New initializes a MediaInfo handle.
func New() *MediaInfo {
	cmi := C.g_MediaInfo_New()
	mi := &MediaInfo{cmi}
	runtime.SetFinalizer(mi, func (mi *MediaInfo) {
		if mi.cptr != nil {
			C.g_MediaInfo_Delete(mi.cptr)
		}
	})
	return mi
}

// Open opens the file at path with the mediainfo library.
func (mi *MediaInfo) Open(path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	rc := C.g_MediaInfo_Open(mi.cptr, cpath)

	if rc != 1 {
		return ErrOpenFailed
	}
	return nil
}

// Option gets a MediaInfo handle option
func (mi *MediaInfo) Option(option string) string {
	coption := C.CString(option)
	defer C.free(unsafe.Pointer(coption))

	return C.GoString(C.g_MediaInfo_Option(mi.cptr, coption, emptyCString()))
}

// SetOption sets a MediaInfo handle option
func (mi *MediaInfo) SetOption(option, value string) string {
	coption, cvalue := C.CString(option), C.CString(value)
	defer C.free(unsafe.Pointer(coption))
	defer C.free(unsafe.Pointer(cvalue))

	return C.GoString(C.g_MediaInfo_Option(mi.cptr, coption, cvalue))
}

// Inform gets the file info (if available) according to previous options set by Option
func (mi *MediaInfo) Inform() string {
	return C.GoString(C.g_MediaInfo_Inform(mi.cptr))
}

// Get gets the file info according to the parameters.
func (mi *MediaInfo) Get(streamKind StreamKind, streamNumber int, parameter string) string {
	return mi.GetKind(streamKind, streamNumber, parameter, InfoText)
}

// GetKind gets the file info according to the parameters.
func (mi *MediaInfo) GetKind(streamKind StreamKind, streamNumber int, parameter string, kindOfInfo InfoKind) string {
	cparameter := C.CString(parameter)
	defer C.free(unsafe.Pointer(cparameter))
	return C.GoString(C.g_MediaInfo_Get(mi.cptr, toCStream(streamKind), C.size_t(streamNumber), cparameter,
		toCInfo(kindOfInfo), toCInfo(InfoName)))
}

// GetAt gets the file info at a particular parameter index.
func (mi *MediaInfo) GetAt(streamKind StreamKind, streamNumber int, parameter int, kindOfInfo InfoKind) string {
	return C.GoString(C.g_MediaInfo_GetI(mi.cptr, toCStream(streamKind), C.size_t(streamNumber), C.size_t(parameter), toCInfo(kindOfInfo)))
}

// Count gets the count of stream.
func (mi *MediaInfo) Count(streamKind StreamKind) int {
	return int(C.g_MediaInfo_Count_Get(mi.cptr, toCStream(streamKind)))
}

// Close closes the handle, releasing internal resources.
func (mi *MediaInfo) Close() {
	C.g_MediaInfo_Close(mi.cptr)
}
