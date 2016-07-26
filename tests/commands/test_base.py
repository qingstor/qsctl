import os
import unittest

from mock import MockOptions

from qingstor.qsctl.commands.base import BaseCommand
from qingstor.qsctl.utils import load_conf

class TestBaseCommand(unittest.TestCase):
    Base = BaseCommand

    def setUp(self):
        # Set the http connection
        conf = load_conf("~/.qingcloud/config.yaml")
        options = MockOptions()
        self.Base.command = "base"
        self.Base.conn = self.Base.get_connection(conf, options)

        # We need some buckets and keys for testing.
        valid_bucket, invalid_bucket = "validbucket", "invalidbucket"
        resp = self.Base.conn.make_request("PUT", valid_bucket)
        resp = self.Base.conn.make_request("HEAD", valid_bucket)
        if resp.status != 200:
            self.fail("setUp failed: please use another bucket name")
        resp = self.Base.conn.make_request("PUT", valid_bucket, "existskey")
        resp.close()

        self.local_path = "./tmp/tmp.file"
        self.valid_bucket = valid_bucket
        self.invalid_bucket = invalid_bucket

    def test_validate_bucket(self):
        self.Base.validate_bucket(self.valid_bucket)
        with self.assertRaises(SystemExit):
            self.Base.validate_bucket(self.invalid_bucket)

    def test_validate_local_path(self):
        self.Base.validate_local_path(self.local_path)
        dirname = os.path.dirname(self.local_path)
        self.assertTrue(os.path.exists(dirname))

    def test_validate_qs_path(self):
        bucket, prefix = self.Base.validate_qs_path("qs://validbucket/prefix")
        self.assertEqual(bucket, "validbucket")
        self.assertEqual(prefix, "prefix")

    def test_key_exists(self):
        self.assertTrue(self.Base.key_exists(self.valid_bucket, "existskey"))
        self.assertFalse(self.Base.key_exists(self.valid_bucket, "noneexistskey"))

    def test_remove_key(self):
        self.Base.remove_key(self.valid_bucket, "existskey")

    def test_confirm_key_remove(self):
        options = MockOptions(exclude="*", include="*.jpg")
        self.assertTrue(self.Base.confirm_key_remove("test.jpg", options))
        self.assertFalse(self.Base.confirm_key_remove("test.txt", options))

    def test_list_multiple_keys_1(self):
        for i in range(0, 9):
            key = "prefix" + "/" + str(i)
            resp = self.Base.conn.make_request("PUT", self.valid_bucket, key)
        resp.close()
        params = {'prefix': 'prefix', 'delimiter': '/'}
        keys, marker, dirs= self.Base.list_multiple_keys(
            self.valid_bucket,
            "",
            params=params
        )
        self.assertEqual(len(keys), 0)
        self.assertEqual(marker, "prefix/8")
        self.assertEqual(dirs, ['prefix/'])

    def test_list_multiple_keys_2(self):
        params = {'prefix': 'prefix'}
        marker = "prefix/5"
        keys, marker, dirs= self.Base.list_multiple_keys(
            self.valid_bucket,
            marker,
            params=params
        )
        self.assertEqual(len(keys), 3)
        self.assertEqual(marker, "prefix/8")
        self.assertEqual(dirs, [])

    def test_remove_multiple_keys(self):
        options = MockOptions(exclude=None, include=None)
        self.Base.remove_multiple_keys(
            self.valid_bucket,
            prefix="prefix",
            options=options
        )

    def tearDown(self):
        resp = self.Base.conn.make_request("DELETE", self.valid_bucket, "existskey")
        resp = self.Base.conn.make_request("DELETE", self.valid_bucket)
        resp.close()
        dirname = os.path.dirname(self.local_path)
        if os.path.exists(dirname):
            os.rmdir(dirname)

if __name__ == "__main__":
    unittest.main(verbosity=2)
