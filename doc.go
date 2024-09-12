// # cgo, C, go type convert
//
//	|C Language Type        | CGO Type      |Go Language Type|     SDK        |
//	|---------------------- | ------------- | -------------- | -------------- |
//	|char                   | C.char        | byte           |                |
//	|singed char            | C.schar       | int8           |                |
//	|unsigned char          | C.uchar       | uint8          | BYTE           |
//	|short                  | C.short       | int16          |                |
//	|unsigned short         | C.ushort      | uint16         | WORD,USHORT    |
//	|int                    | C.int         | int32          | BOOL,LONG      |
//	|unsigned int           | C.uint        | uint32         | DWORD,UINT     |
//	|long                   | C.long        | int32          |                |
//	|unsigned long          | C.ulong       | uint32         |                |
//	|long long int          | C.longlong    | int64          |                |
//	|unsigned long long int | C.ulonglong   | uint64         |                |
//	|float                  | C.float       | float32        |                |
//	|double                 | C.double      | float64        |                |
//	|size_t                 | C.size_t      | uint           |                |
//	|void*                  | unsafe.Pointer| unsafe.Pointer | LPVOID.HANDLE  |
//
// # cgo, C, go string convert
//
//   - func C.CString(string) *C.char
//
//     Go string to C string.
//     The C string is allocated in the C heap using malloc.
//     It is the caller's responsibility to arrange for it to be
//     freed, such as by calling C.free (be sure to include stdlib.h
//     if C.free is needed).
//
//   - func C.CBytes([]byte) unsafe.Pointer
//
//     Go []byte slice to C array.
//     The C array is allocated in the C heap using malloc.
//     It is the caller's responsibility to arrange for it to be
//     freed, such as by calling C.free (be sure to include stdlib.h
//     if C.free is needed).
//
//   - func C.GoString(*C.char) string
//
//     C string to Go string
//
//   - func C.GoStringN(*C.char, C.int) string
//
//     C data with explicit length to Go string
//
//   - func C.GoBytes(unsafe.Pointer, C.int) []byte
//
//     C data with explicit length to Go []byte
//
// # C header -> CGO header
//
//  1. delete __stdcall
//
//  2. delete CALLBACK
//
//  3. enum XXX{}; -> type enum _XXX {}XXX;
//
//  4. delete function with init parameters, pUser = NULL
//
//  5. struct must have tag name
package libmodbusgo
