//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2025 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/weaviate/weaviate/entities/backup"
	"github.com/weaviate/weaviate/entities/moduletools"
	mod "github.com/weaviate/weaviate/modules/backup-gcs"
	"github.com/weaviate/weaviate/test/docker"
	moduleshelper "github.com/weaviate/weaviate/test/helper/modules"
	ubak "github.com/weaviate/weaviate/usecases/backup"
)

func Test_GcsBackend_Start(t *testing.T) {
	gCSBackend_Backup(t, "", "")

	gCSBackend_Backup(t, "gcsbetestbucketoverride", "gcsbetestBucketPathOverride")
}

func gCSBackend_Backup(t *testing.T, overrideBucket, overridePath string) {
	ctx := context.Background()
	compose, err := docker.New().WithGCS().Start(ctx)
	if err != nil {
		t.Fatal(errors.Wrapf(err, "cannot start"))
	}

	t.Setenv(envGCSEndpoint, compose.GetGCS().URI())

	t.Run("store backup meta", func(t *testing.T) { moduleLevelStoreBackupMeta(t, overrideBucket, overridePath) })
	t.Run("copy objects", func(t *testing.T) { moduleLevelCopyObjects(t, overrideBucket, overridePath) })
	t.Run("copy files", func(t *testing.T) { moduleLevelCopyFiles(t, overrideBucket, overridePath) })

	if err := compose.Terminate(ctx); err != nil {
		t.Fatal(errors.Wrapf(err, "failed to terminate test containers"))
	}
}

func moduleLevelStoreBackupMeta(t *testing.T, overrideBucket, overridePath string) {
	testCtx := context.Background()

	className := "BackupClass"
	backupID := "backup_id"
	bucketName := "bucket-level-store-backup-meta"
	if overrideBucket != "" {
		bucketName = overrideBucket
	}
	projectID := "project-id"
	endpoint := os.Getenv(envGCSEndpoint)
	metadataFilename := "backup.json"
	gcsUseAuth := "false"

	t.Log("setup env")
	t.Setenv(envGCSEndpoint, endpoint)
	t.Setenv(envGCSStorageEmulatorHost, endpoint)
	t.Setenv(envGCSCredentials, "")
	t.Setenv(envGCSProjectID, projectID)
	t.Setenv(envGCSBucket, bucketName)
	t.Setenv(envGCSUseAuth, gcsUseAuth)
	moduleshelper.CreateGCSBucket(testCtx, t, projectID, bucketName)
	defer moduleshelper.DeleteGCSBucket(testCtx, t, bucketName)

	t.Run("store backup meta in gcs", func(t *testing.T) {
		t.Setenv("BACKUP_GCS_BUCKET", bucketName)
		gcs := mod.New()
		params := moduletools.NewMockModuleInitParams(t)
		params.EXPECT().GetLogger().Return(logrus.New())
		params.EXPECT().GetStorageProvider().Return(&fakeStorageProvider{dataPath: t.TempDir()})
		err := gcs.Init(testCtx, params)
		require.Nil(t, err)

		t.Run("access permissions", func(t *testing.T) {
			err := gcs.Initialize(testCtx, backupID, overrideBucket, overridePath)
			assert.Nil(t, err)
		})

		t.Run("backup meta does not exist yet", func(t *testing.T) {
			meta, err := gcs.GetObject(testCtx, backupID, metadataFilename, overrideBucket, overridePath)
			assert.Nil(t, meta)
			assert.NotNil(t, err)
			assert.IsType(t, backup.ErrNotFound{}, err)
		})

		t.Run("put backup meta on backend", func(t *testing.T) {
			desc := &backup.BackupDescriptor{
				StartedAt:   time.Now(),
				CompletedAt: time.Time{},
				ID:          backupID,
				Classes: []backup.ClassDescriptor{
					{
						Name: className,
					},
				},
				Status:  string(backup.Started),
				Version: ubak.Version,
			}

			b, err := json.Marshal(desc)
			require.Nil(t, err)

			err = gcs.PutObject(testCtx, backupID, metadataFilename, overrideBucket, overridePath, b)
			require.Nil(t, err)

			dest := gcs.HomeDir(backupID, overrideBucket, overridePath)
			if overridePath == "" {
				expected := fmt.Sprintf("gs://%s/%s", bucketName, backupID)
				assert.Equal(t, expected, dest)
			} else {
				expected := fmt.Sprintf("gs://%s/%s/%s", bucketName, overridePath, backupID)
				assert.Equal(t, expected, dest)
			}
		})

		t.Run("assert backup meta contents", func(t *testing.T) {
			obj, err := gcs.GetObject(testCtx, backupID, metadataFilename, overrideBucket, overridePath)
			require.Nil(t, err)

			var meta backup.BackupDescriptor
			err = json.Unmarshal(obj, &meta)
			require.Nil(t, err)
			assert.NotEmpty(t, meta.StartedAt)
			assert.Empty(t, meta.CompletedAt)
			assert.Equal(t, meta.Status, string(backup.Started))
			assert.Empty(t, meta.Error)
			assert.Len(t, meta.Classes, 1)
			assert.Equal(t, meta.Classes[0].Name, className)
			assert.Equal(t, meta.Version, ubak.Version)
			assert.Nil(t, meta.Classes[0].Error)
		})
	})
}

func moduleLevelCopyObjects(t *testing.T, overrideBucket, overridePath string) {
	testCtx := context.Background()

	key := "moduleLevelCopyObjects"
	backupID := "backup_id"
	bucketName := "bucket-level-copy-objects"
	if overrideBucket != "" {
		bucketName = overrideBucket
	}
	projectID := "project-id"
	endpoint := os.Getenv(envGCSEndpoint)
	gcsUseAuth := "false"

	t.Log("setup env")
	t.Setenv(envGCSEndpoint, endpoint)
	t.Setenv(envGCSStorageEmulatorHost, endpoint)
	t.Setenv(envGCSCredentials, "")
	t.Setenv(envGCSProjectID, projectID)
	t.Setenv(envGCSBucket, bucketName)
	t.Setenv(envGCSUseAuth, gcsUseAuth)
	moduleshelper.CreateGCSBucket(testCtx, t, projectID, bucketName)
	defer moduleshelper.DeleteGCSBucket(testCtx, t, bucketName)

	t.Run("copy objects", func(t *testing.T) {
		t.Setenv("BACKUP_GCS_BUCKET", bucketName)
		gcs := mod.New()
		params := moduletools.NewMockModuleInitParams(t)
		params.EXPECT().GetLogger().Return(logrus.New())
		params.EXPECT().GetStorageProvider().Return(&fakeStorageProvider{dataPath: t.TempDir()})
		err := gcs.Init(testCtx, params)
		require.Nil(t, err)

		t.Run("put object to bucket", func(t *testing.T) {
			err := gcs.PutObject(testCtx, backupID, key, overrideBucket, overridePath, []byte("hello"))
			assert.Nil(t, err)
		})

		t.Run("get object from bucket", func(t *testing.T) {
			meta, err := gcs.GetObject(testCtx, backupID, key, overrideBucket, overridePath)
			assert.Nil(t, err)
			assert.Equal(t, []byte("hello"), meta)
		})
	})
}

func moduleLevelCopyFiles(t *testing.T, overrideBucket, overridePath string) {
	testCtx := context.Background()

	dataDir := t.TempDir()
	key := "moduleLevelCopyFiles"
	backupID := "backup_id"
	bucketName := "bucket-level-copy-files"
	if overrideBucket != "" {
		bucketName = overrideBucket
	}
	projectID := "project-id"
	endpoint := os.Getenv(envGCSEndpoint)
	gcsUseAuth := "false"

	t.Log("setup env")
	t.Setenv(envGCSEndpoint, endpoint)
	t.Setenv(envGCSStorageEmulatorHost, endpoint)
	t.Setenv(envGCSCredentials, "")
	t.Setenv(envGCSProjectID, projectID)
	t.Setenv(envGCSBucket, bucketName)
	t.Setenv(envGCSUseAuth, gcsUseAuth)
	moduleshelper.CreateGCSBucket(testCtx, t, projectID, bucketName)
	defer moduleshelper.DeleteGCSBucket(testCtx, t, bucketName)

	t.Run("copy files", func(t *testing.T) {
		fpaths := moduleshelper.CreateTestFiles(t, dataDir)
		fpath := fpaths[0]
		expectedContents, err := os.ReadFile(fpath)
		require.Nil(t, err)
		require.NotNil(t, expectedContents)

		t.Setenv("BACKUP_GCS_BUCKET", bucketName)
		gcs := mod.New()
		params := moduletools.NewMockModuleInitParams(t)
		params.EXPECT().GetLogger().Return(logrus.New())
		params.EXPECT().GetStorageProvider().Return(&fakeStorageProvider{dataPath: dataDir})
		err = gcs.Init(testCtx, params)
		require.Nil(t, err)

		t.Run("verify source data path", func(t *testing.T) {
			assert.Equal(t, dataDir, gcs.SourceDataPath())
		})

		t.Run("copy file to backend", func(t *testing.T) {
			err := gcs.PutObject(testCtx, backupID, key, overrideBucket, overridePath, expectedContents)
			require.Nil(t, err)

			contents, err := gcs.GetObject(testCtx, backupID, key, overrideBucket, overridePath)
			require.Nil(t, err)
			assert.Equal(t, expectedContents, contents)
		})

		t.Run("fetch file from backend", func(t *testing.T) {
			destPath := dataDir + "/file_0.copy.db"

			err := gcs.WriteToFile(testCtx, backupID, key, destPath, overrideBucket, overridePath)
			require.Nil(t, err)

			contents, err := os.ReadFile(destPath)
			require.Nil(t, err)
			assert.Equal(t, expectedContents, contents)
		})
	})
}

type fakeStorageProvider struct {
	dataPath string
}

func (f *fakeStorageProvider) Storage(name string) (moduletools.Storage, error) {
	return nil, nil
}

func (f *fakeStorageProvider) DataPath() string {
	return f.dataPath
}
