import unittest

from tests.test_data import zone, test_bucket1, test_bucket2

from qingstor.qsctl.commands.rb import RbCommand
from qingstor.qsctl.utils import load_conf

from mock import MockOptions


class TestMbCommand(unittest.TestCase):
    Rb = RbCommand

    def setUp(self):
        # Set the http connection
        conf = load_conf("~/.qingstor/config.yaml")
        options = MockOptions()
        self.Rb.client = self.Rb.get_client(conf)

        self.test_bucket = self.Rb.client.Bucket(test_bucket2, zone)
        self.test_bucket.put()
        resp = self.test_bucket.head()
        if resp.status_code != 200:
            self.fail("setUp failed: please use another bucket name")

        self.test_bucket.put_object("testkey1")
        self.test_bucket.put_object("testkey2")

    def test_remove_bucket_1(self):
        options = MockOptions(bucket=test_bucket2, force=False)
        self.Rb.send_request(options)

    def test_remove_bucket_2(self):
        options = MockOptions(bucket=test_bucket2, force=True)
        self.Rb.send_request(options)


if __name__ == "__main__":
    unittest.main(verbosity=2)
