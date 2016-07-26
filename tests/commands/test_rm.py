import unittest

from mock import MockOptions

from qingstor.qsctl.commands.rm import RmCommand
from qingstor.qsctl.utils import load_conf

class TestRmCommand(unittest.TestCase):
    Rm = RmCommand

    def setUp(self):

        # Set the http connection
        conf = load_conf("~/.qingcloud/config.yaml")
        options = MockOptions()
        self.Rm.conn = self.Rm.get_connection(conf, options)

        # We need a bucket for testing.
        valid_bucket = "validbucket"
        resp = self.Rm.conn.make_request("PUT", valid_bucket)
        resp = self.Rm.conn.make_request("HEAD", valid_bucket)
        if resp.status != 200:
            self.fail("setUp failed: please use another bucket name")
        resp.close()

        self.valid_bucket = valid_bucket

    def test_remove_one_key(self):
        resp = self.Rm.conn.make_request("PUT", self.valid_bucket, "testkey")
        resp.close()
        options = MockOptions(qs_path="qs://validbucket/testkey", recursive=False)
        self.Rm.send_request(options)

    def test_remove_mutiple_keys_1(self):
        for i in range(0, 10):
            key = "prefix/" + str(i)
            resp = self.Rm.conn.make_request("PUT", self.valid_bucket, key)
        resp.close()

        options = MockOptions(
            qs_path="qs://validbucket/prefix/",
            recursive=True,
            exclude=None,
            include=None
        )
        self.Rm.send_request(options)

    def test_remove_mutiple_keys_2(self):
        for i in range(0, 10):
            key = "prefix/" + str(i) + ".txt"
            resp = self.Rm.conn.make_request("PUT", self.valid_bucket, key)
        resp = self.Rm.conn.make_request("PUT", self.valid_bucket, "prefix/test.jpg")
        resp.close()

        options = MockOptions(
            qs_path="qs://validbucket/prefix/",
            recursive=True,
            exclude="*.txt",
            include=None
        )
        self.Rm.send_request(options)

    def test_remove_mutiple_keys_3(self):
        for i in range(0, 10):
            key = "prefix/" + str(i) + ".txt"
            resp = self.Rm.conn.make_request("PUT", self.valid_bucket, key)
        resp = self.Rm.conn.make_request("PUT", self.valid_bucket, "prefix/test.jpg")
        resp.close()

        options = MockOptions(
            qs_path="qs://validbucket/prefix/",
            recursive=True,
            exclude="*",
            include="*.txt"
        )
        self.Rm.send_request(options)

    def tearDown(self):
        options = MockOptions(exclude=None, include=None)
        self.Rm.remove_multiple_keys(self.valid_bucket, options=options)
        resp = self.Rm.conn.make_request("DELETE", self.valid_bucket)
        resp.close()

if __name__ == "__main__":
    unittest.main()
