import unittest

from mock import MockOptions

from qingstor.qsctl.commands.mb import MbCommand
from qingstor.qsctl.utils import load_conf

class TestMbCommand(unittest.TestCase):
    Mb = MbCommand

    def setUp(self):
        # Set the http connection
        conf = load_conf("~/.qingcloud/config.yaml")
        options = MockOptions()
        self.Mb.conn = self.Mb.get_connection(conf, options)

    def test_make_bucket_1(self):
        options = MockOptions(bucket="validbucket", zone="pek3a")
        self.Mb.send_request(options)

    def test_make_bucket_2(self):
        options = MockOptions(bucket="validbucket", zone=None)
        self.Mb.send_request(options)

    def tearDown(self):
        resp = self.Mb.conn.make_request("DELETE", "validbucket")
        resp.close()

if __name__ == "__main__":
    unittest.main(verbosity=2)
