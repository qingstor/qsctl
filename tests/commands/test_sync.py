import os
import time
import shutil
import unittest

from mock_options import MockOptions

from tests.test_data import zone, test_bucket1, test_bucket2

from qingstor.qsctl.commands.sync import SyncCommand
from qingstor.qsctl.utils import load_conf


class TestSyncCommand(unittest.TestCase):
    Sync = SyncCommand

    def setUp(self):
        # Set the http connection
        conf = load_conf("~/.qingstor/config.yaml")
        options = MockOptions()
        self.Sync.client = self.Sync.get_client(conf)

        self.test_bucket = self.Sync.client.Bucket(test_bucket1, zone)
        self.test_bucket.put()
        resp = self.test_bucket.head()
        if resp.status_code != 200:
            self.fail("setUp failed: please use another bucket name")

        # Temp directory for testing
        if not os.path.exists("tmp/"):
            os.mkdir("tmp/")

    def test_confirm_key_remove(self):
        # Test case 1: key exists in bucket but not exsits in source
        # local directory will be deleted.
        with open("tmp/test.file2", 'w') as f:
            f.write("just for testing")
        options = MockOptions(
            source_path="tmp/",
            exclude=None,
            include=None,
        )
        key1 = "test.file1"
        key2 = "test.file2"
        self.assertTrue(self.Sync.confirm_key_remove(key1, options))
        self.assertFalse(self.Sync.confirm_key_remove(key2, options))

        # Test case 2: key not match the pattern will be deleted
        with open("tmp/test.txt", 'w') as f:
            f.write("just for testing")
        with open("tmp/test.jpg", 'w') as f:
            f.write("just for testing")
        options = MockOptions(
            source_path="tmp/",
            exclude="*",
            include="*.jpg",
        )
        key1 = "test.txt"
        key2 = "test.jpg"
        self.assertTrue(self.Sync.confirm_key_remove(key1, options))
        self.assertFalse(self.Sync.confirm_key_remove(key2, options))

    def test_confirm_key_download(self):
        # Prepare some files and time stamps
        time_stamp_1 = time.time()
        time.sleep(3)
        local_path = "tmp/test.file"
        with open(local_path, 'w') as f:
            f.write("just for testing")
        time.sleep(3)
        time_stamp_2 = time.time()

        # Testing begin
        options = MockOptions()
        confirm_1 = self.Sync.confirm_key_download(
            options, local_path, time_stamp_1
        )
        confirm_2 = self.Sync.confirm_key_download(
            options, local_path, time_stamp_2
        )
        confirm_3 = self.Sync.confirm_key_download(
            options, "tmp/noneexistsfile", time_stamp_1
        )

        self.assertFalse(confirm_1)
        self.assertTrue(confirm_2)
        self.assertTrue(confirm_3)

    def test_confirm_key_upload(self):
        # Prepare local file
        local_path = "tmp/test.file"
        with open(local_path, 'w') as f:
            f.write("just for testing")

        time.sleep(3)

        # Prepare key
        key = "testkey"
        self.test_bucket.put_object(key)

        # Testing begin
        options = MockOptions()
        confirm_1 = self.Sync.confirm_key_upload(
            options, local_path, test_bucket1, key
        )

        time.sleep(3)

        # Change the modified time of local file
        with open(local_path, 'w') as f:
            f.write("write something new")

        confirm_2 = self.Sync.confirm_key_upload(
            options, local_path, test_bucket1, key
        )

        confirm_3 = self.Sync.confirm_key_upload(
            options, local_path, test_bucket1, "noneexistskey"
        )

        self.assertFalse(confirm_1)
        self.assertTrue(confirm_2)
        self.assertTrue(confirm_3)

    def test_is_local_file_modified(self):
        time_stamp_1 = time.mktime(time.gmtime(time.time()))
        time.sleep(3)
        local_path = "tmp/test.file"
        with open(local_path, 'w') as f:
            f.write("just for testing")
        time.sleep(3)
        time_stamp_2 = time.mktime(time.gmtime(time.time()))
        is_modified_1 = self.Sync.is_local_file_modified(
            local_path, time_stamp_1
        )
        is_modified_2 = self.Sync.is_local_file_modified(
            local_path, time_stamp_2
        )
        self.assertTrue(is_modified_1)
        self.assertFalse(is_modified_2)

    def test_upload_files(self):
        # Create some local files for testing.
        for i in range(1, 10):
            with open("tmp/file" + str(i), 'w') as f:
                f.write("just for testing")
        options = MockOptions(
            source_path="tmp/",
            dest_path="qs://" + test_bucket1,
            exclude="*",
            include="*",
            delete=True
        )
        print(options.exclude)
        self.Sync.upload_files(options)

    def test_download_files(self):
        for i in range(0, 10):
            key = "test" + str(i)
            self.test_bucket.put_object(key)
        options = MockOptions(
            source_path="qs://" + test_bucket1,
            dest_path="tmp/",
            exclude=None,
            include=None,
            delete=True
        )
        self.Sync.download_files(options)

    def tearDown(self):
        shutil.rmtree("tmp/")
        options = MockOptions(exclude=None, include=None, source_path="tmp/")
        self.Sync.remove_multiple_keys(test_bucket1, options=options)


if __name__ == "__main__":
    unittest.main(verbosity=2)
