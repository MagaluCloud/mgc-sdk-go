BlockStorage
= audit availabilityzones blockstorage client cmd compute containerregistry cover.out dbaas docs go.mod go.sum helpers internal kubernetes lbaas LICENSE Makefile network README.md scripts sonar-project.properties sshkeys 12

Package Documentation
-------------------

.. code-block:: go

   package blockstorage // import "github.com/MagaluCloud/mgc-sdk-go/blockstorage"
   
   Package blockstorage provides functionality to interact with the MagaluCloud
   block storage service. This package allows managing volumes, volume types,
   and snapshots.
   
   Package blockstorage provides functionality to interact with the MagaluCloud
   block storage service. This package allows managing volumes, volume types,
   and snapshots.
   
   Package blockstorage provides functionality to interact with the MagaluCloud
   block storage service. This package allows managing volumes, volume types,
   and snapshots.
   
   Package blockstorage provides functionality to interact with the MagaluCloud
   block storage service. This package allows managing volumes, volume types,
   and snapshots.
   
   const VolumeTypeExpand = "volume_type" ...
   const DefaultBasePath = "/volume"


Types
-----

- :type:`AttachmentInstance`
- :type:`BlockStorageClient`
- :type:`ClientOption`
- :type:`CreateSnapshotRequest`
- :type:`CreateVolumeRequest`
- :type:`DiskType`
- :type:`ExtendVolumeRequest`
- :type:`IDOrName`
- :type:`Iops`
- :type:`ListOptions`
- :type:`ListSnapshotsResponse`
- :type:`ListVolumeTypesOptions`
- :type:`ListVolumeTypesResponse`
- :type:`ListVolumesResponse`
- :type:`RenameSnapshotRequest`
- :type:`RenameVolumeRequest`
- :type:`RetypeVolumeRequest`
- :type:`Snapshot`
- :type:`SnapshotError`
- :type:`SnapshotService`
- :type:`SnapshotStateV1`
- :type:`SnapshotStatusV1`
- :type:`Type`
- :type:`Volume`
- :type:`VolumeAttachment`
- :type:`VolumeError`
- :type:`VolumeService`
- :type:`VolumeStateV1`
- :type:`VolumeStatusV1`
- :type:`VolumeType`
- :type:`VolumeTypeIOPS`
- :type:`VolumeTypeService`

Constants
---------

- :const:`VolumeTypeExpand`
- :const:`DefaultBasePath`
- :const:`SnapshotVolumeExpand`

Example Usage
-------------

.. code-block:: go

   import "github.com/magalucloud/mgc-sdk-go/blockstorage"

   // Use the BlockStorage package
   // See the examples directory for complete examples

