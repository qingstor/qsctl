import os
import sys
import unittest

from mock import MockOptions

from qingstor.qsctl.commands.ls import LsCommand
from qingstor.qsctl.utils import load_conf

class TestLsCommand(unittest.TestCase):
    Ls = LsCommand

    def setUp(self):

        # Set the http connection
        conf = load_conf("~/.qingcloud/config.yaml")
        options = MockOptions(qs_path="qs://")
        self.Ls.conn = self.Ls.get_connection(conf, options)

        # We need some buckets for testing.
        valid_bucket1 = "validbucket1"
        valid_bucket2 = "validbucket2"
        resp = self.Ls.conn.make_request("PUT", valid_bucket1)
        resp = self.Ls.conn.make_request("HEAD", valid_bucket1)
        if resp.status != 200:
            self.fail("setUp failed: please use another bucket name")
        resp = self.Ls.conn.make_request("PUT", valid_bucket2)
        resp = self.Ls.conn.make_request("HEAD", valid_bucket2)
        if resp.status != 200:
            self.fail("setUp failed: please use another bucket name")
        resp.close()

        self.valid_bucket1 = valid_bucket1
        self.valid_bucket2 = valid_bucket2

    def test_list_buckets_1(self):
        options = MockOptions(zone="pek3a")
        self.Ls.list_buckets(options)

    def test_list_buckets_2(self):
        options = MockOptions(zone=None)
        self.Ls.list_buckets(options)

    def test_list_keys_1(self):
        # Set the http connection for listing keys
        options = MockOptions(
            qs_path="qs://validbucket1",
            recursive=False,
            page_size=None
        )
        conf = load_conf("~/.qingcloud/config.yaml")
        self.Ls.conn = self.Ls.get_connection(conf, options)

        # We need some keys for testing.
        for i in range(0, 10):
            resp = self.Ls.conn.make_request("PUT", self.valid_bucket1, "test"+str(i))
        resp.close()

        # testing
        self.Ls.list_keys(options)

        # clean keys after testing
        options = MockOptions(exclude=None, include=None)
        self.Ls.remove_multiple_keys(self.valid_bucket1, options=options)

    def test_list_keys_2(self):
        # Set the http connection for listing objects
        options = MockOptions(
            qs_path="qs://validbucket1/prefix/",
            recursive=True,
            page_size=None
        )
        conf = load_conf("~/.qingcloud/config.yaml")
        self.Ls.conn = self.Ls.get_connection(conf, options)

        # We need some keys for testing.
        resp = self.Ls.conn.make_request("PUT", self.valid_bucket1, "test.txt")
        for i in range(0, 10):
            resp = self.Ls.conn.make_request("PUT", self.valid_bucket1, "prefix/"+str(i))
        resp.close()

        # testing
        self.Ls.list_keys(options)

        # clean keys after testing
        options = MockOptions(exclude=None, include=None)
        self.Ls.remove_multiple_keys(self.valid_bucket1, options=options)

    def test_list_keys_3(self):
        # Set the http connection for listing objects
        options = MockOptions(
            qs_path="qs://validbucket1",
            recursive=False,
            page_size=None
        )
        conf = load_conf("~/.qingcloud/config.yaml")
        self.Ls.conn = self.Ls.get_connection(conf, options)

        # We need some keys for testing.
        resp = self.Ls.conn.make_request("PUT", self.valid_bucket1, "test.txt")
        for i in range(0, 10):
            resp = self.Ls.conn.make_request("PUT", self.valid_bucket1, "prefix/"+str(i))
        resp.close()

        # testing
        self.Ls.list_keys(options)

        # clean keys after testing
        options = MockOptions(exclude=None, include=None)
        self.Ls.remove_multiple_keys(self.valid_bucket1, options=options)

    def tearDown(self):
        resp = self.Ls.conn.make_request("DELETE", self.valid_bucket1)
        resp = self.Ls.conn.make_request("DELETE", self.valid_bucket2)
        resp.close()

if __name__ == "__main__":
    unittest.main(verbosity=2)
