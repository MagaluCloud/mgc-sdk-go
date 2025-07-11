Helpers
= audit availabilityzones blockstorage client cmd compute containerregistry cover.out dbaas docs go.mod go.sum helpers internal kubernetes lbaas LICENSE Makefile network README.md scripts sonar-project.properties sshkeys 7

Package Documentation
-------------------

.. code-block:: go

   package helpers // import "github.com/MagaluCloud/mgc-sdk-go/helpers"
   
   func BoolPtr(b bool) *bool
   func Float32Ptr(f float32) *float32
   func Float64Ptr(f float64) *float64
   func Int16Ptr(i int16) *int16
   func Int32Ptr(i int32) *int32
   func Int64Ptr(i int64) *int64
   func Int8Ptr(i int8) *int8
   func IntPtr(i int) *int
   func StrPtr(s string) *string
   func Uint16Ptr(u uint16) *uint16
   func Uint32Ptr(u uint32) *uint32
   func Uint64Ptr(u uint64) *uint64
   func Uint8Ptr(u uint8) *uint8
   func UintPtr(u uint) *uint
   type QueryParams interface{ ... }
       func NewQueryParams(httpReq *http.Request) QueryParams


Functions
---------

- :func:`BoolPtr`
- :func:`Float32Ptr`
- :func:`Float64Ptr`
- :func:`Int16Ptr`
- :func:`Int32Ptr`
- :func:`Int64Ptr`
- :func:`Int8Ptr`
- :func:`IntPtr`
- :func:`StrPtr`
- :func:`Uint16Ptr`
- :func:`Uint32Ptr`
- :func:`Uint64Ptr`
- :func:`Uint8Ptr`
- :func:`UintPtr`

Types
-----

- :type:`QueryParams`

Example Usage
-------------

.. code-block:: go

   import "github.com/magalucloud/mgc-sdk-go/helpers"

   // Use the Helpers package
   // See the examples directory for complete examples

