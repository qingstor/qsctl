import os
import shutil
import unittest

from tests.test_data import zone, test_bucket1, test_bucket2

from qingstor.qsctl.commands.cp import CpCommand
from qingstor.qsctl.utils import load_conf

from mock import MockOptions


class TestCpCommand(unittest.TestCase):
    Cp = CpCommand

    def setUp(self):
        # Set the http connection
        conf = load_conf("~/.qingstor/config.yaml")
        options = MockOptions()
        self.Cp.client = self.Cp.get_client(conf)

        self.test_bucket = self.Cp.client.Bucket(test_bucket1, zone)
        self.test_bucket.put()
        resp = self.test_bucket.head()
        if resp.status_code != 200:
            self.fail("setUp failed: please use another bucket name")

        # Temp directory for testing
        if not os.path.exists("tmp/"):
            os.mkdir("tmp/")

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
            force=True)
        print(options.exclude)
        self.Cp.upload_files(options)

    def test_upload_file(self):
        # Create a large file(~8MB)
        with open("tmp/large_file", 'w') as f:
            f.seek(8 * 1024 * 1024)
            f.write("just for testing")
        options = MockOptions(
            source_path="tmp/large_file",
            dest_path="qs://" + test_bucket1,
            exclude=None,
            include=None,
            force=True)
        self.Cp.upload_file(options)

    def test_download_files(self):
        for i in range(0, 10):
            key = "test" + str(i)
            resp = self.test_bucket.put_object(key)
        options = MockOptions(
            source_path="qs://" + test_bucket1,
            dest_path="tmp/",
            exclude=None,
            include=None,
            force=True)
        self.Cp.download_files(options)

    def test_download_file(self):
        resp = self.test_bucket.put_object("test_file")
        options = MockOptions(
            source_path="qs://" + test_bucket1 + "/test_file",
            dest_path="tmp/",
            exclude=None,
            include=None,
            force=True)
        self.Cp.download_file(options)
        self.assertTrue((os.path.exists("tmp/test_file")))

    def tearDown(self):
        options = MockOptions(exclude=None, include=None)
        self.Cp.remove_multiple_keys(test_bucket1, options=options)
        shutil.rmtree("tmp/")


if __name__ == "__main__":
    unittest.main(verbosity=2)
