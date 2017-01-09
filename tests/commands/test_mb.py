import unittest

from mock import MockOptions

from tests.test_data import zone, test_bucket1, test_bucket2

from qingstor.qsctl.commands.mb import MbCommand
from qingstor.qsctl.utils import load_conf


class TestMbCommand(unittest.TestCase):
    Mb = MbCommand

    def setUp(self):
        # Set the http connection
        conf = load_conf("~/.qingstor/config.yaml")
        options = MockOptions()
        self.Mb.client = self.Mb.get_client(conf)

    def test_make_bucket_1(self):
        options = MockOptions(bucket=test_bucket2, zone=zone)
        self.Mb.send_request(options)

    def test_make_bucket_2(self):
        options = MockOptions(bucket=test_bucket2, zone=None)
        self.Mb.send_request(options)

    def tearDown(self):
        test_bucket = self.Mb.client.Bucket(test_bucket2, zone)
        test_bucket.delete()


if __name__ == "__main__":
    unittest.main(verbosity=2)
