import os
import unittest

from mock_options import MockOptions

from tests.test_data import zone, test_bucket1, test_bucket2

from qingstor.qsctl.commands.base import BaseCommand
from qingstor.qsctl.utils import load_conf


class TestBaseCommand(unittest.TestCase):
    Base = BaseCommand

    def setUp(self):
        # Set the http connection
        conf = load_conf("~/.qingstor/config.yaml")
        options = MockOptions()
        self.Base.command = "base"
        self.Base.client = self.Base.get_client(conf)

        self.test_bucket = self.Base.client.Bucket(test_bucket1, zone)
        self.test_bucket.put()
        resp = self.test_bucket.head()
        if resp.status_code != 200:
            self.fail("setUp failed: please use another bucket name")
        self.test_bucket.put_object("existskey")

        self.local_path = "./tmp/tmp.file"

    def test_validate_bucket(self):
        self.Base.validate_bucket(test_bucket1)
        with self.assertRaises(SystemExit):
            self.Base.validate_bucket(test_bucket2)

    def test_validate_local_path(self):
        self.Base.validate_local_path(self.local_path)
        dirname = os.path.dirname(self.local_path)
        self.assertTrue(os.path.exists(dirname))

    def test_validate_qs_path(self):
        bucket, prefix = self.Base.validate_qs_path("qs://" + test_bucket1 +
                                                    "/prefix")
        self.assertEqual(bucket, test_bucket1)
        self.assertEqual(prefix, "prefix")

    def test_key_exists(self):
        self.assertTrue(self.Base.key_exists(test_bucket1, "existskey"))
        self.assertFalse(self.Base.key_exists(test_bucket1, "noneexistskey"))

    def test_remove_key(self):
        self.Base.remove_key(test_bucket1, "existskey")

    def test_confirm_key_remove(self):
        options = MockOptions(exclude="*", include="*.jpg")
        self.assertTrue(self.Base.confirm_key_remove("test.jpg", options))
        self.assertFalse(self.Base.confirm_key_remove("test.txt", options))

    def test_list_multiple_keys_1(self):
        for i in range(0, 9):
            key = "prefix" + "/" + str(i)
            self.test_bucket.put_object(key)
        keys, next_marker, dirs = self.Base.list_multiple_keys(
            test_bucket1, prefix="prefix", delimiter="/")
        self.assertEqual(len(keys), 0)
        self.assertEqual(next_marker, "")
        self.assertEqual(dirs, ['prefix/'])

    def test_list_multiple_keys_2(self):
        marker = "prefix/5"
        keys, next_marker, dirs = self.Base.list_multiple_keys(
            test_bucket1, marker=marker, prefix="prefix", limit=1)
        self.assertEqual(len(keys), 1)
        self.assertEqual(next_marker, "prefix/6")
        self.assertEqual(dirs, [])

    def test_remove_multiple_keys(self):
        options = MockOptions(exclude=None, include=None)
        self.Base.remove_multiple_keys(
            test_bucket1, prefix="prefix", options=options)

    def tearDown(self):
        self.test_bucket.delete_object("existskey")
        dirname = os.path.dirname(self.local_path)
        if os.path.exists(dirname):
            os.rmdir(dirname)


if __name__ == "__main__":
    unittest.main(verbosity=2)
