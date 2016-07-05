import os
import sys
import unittest

from qingstor.qsctl.commands.rb import RbCommand
from qingstor.qsctl.utils import load_conf

from mock import MockOptions

class TestMbCommand(unittest.TestCase):
    Rb = RbCommand

    def setUp(self):
        # Set the http connection
        conf = load_conf("~/.qingcloud/config.yaml")
        options = MockOptions()
        self.Rb.conn = self.Rb.get_connection(conf, options)

        # Prepare a bucket and some keys for testing
        valid_bucket = "validbucket"
        resp = self.Rb.conn.make_request("PUT", valid_bucket)
        resp = self.Rb.conn.make_request("HEAD", valid_bucket)
        if resp.status != 200:
            self.fail("setUp failed: please use another bucket name")

        resp = self.Rb.conn.make_request("PUT", valid_bucket, "testkey1")
        resp = self.Rb.conn.make_request("PUT", valid_bucket, "testkey2")
        resp.close()

        self.valid_bucket = valid_bucket

    def test_remove_bucket_1(self):
        options = MockOptions(bucket=self.valid_bucket, force=False)
        self.Rb.send_request(options)

    def test_remove_bucket_2(self):
        options = MockOptions(bucket=self.valid_bucket, force=True)
        self.Rb.send_request(options)

    def tearDown(self):
        resp = self.Rb.conn.make_request("DELETE", self.valid_bucket)
        resp.close()

if __name__ == "__main__":
    unittest.main(verbosity=2)
