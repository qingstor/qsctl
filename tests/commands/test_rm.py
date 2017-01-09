import unittest

from mock import MockOptions

from tests.test_data import zone, test_bucket1, test_bucket2

from qingstor.qsctl.commands.rm import RmCommand
from qingstor.qsctl.utils import load_conf


class TestRmCommand(unittest.TestCase):
    Rm = RmCommand

    def setUp(self):

        # Set the http connection
        conf = load_conf("~/.qingstor/config.yaml")
        options = MockOptions()
        self.Rm.client = self.Rm.get_client(conf)

        self.test_bucket = self.Rm.client.Bucket(test_bucket1, zone)
        self.test_bucket.put()
        resp = self.test_bucket.head()
        if resp.status_code != 200:
            self.fail("setUp failed: please use another bucket name")

    def test_remove_one_key(self):
        self.test_bucket.put_object("testkey")
        options = MockOptions(
            qs_path="qs://" + test_bucket1 + "/testkey", recursive=False)
        self.Rm.send_request(options)

    def test_remove_mutiple_keys_1(self):
        for i in range(0, 10):
            key = "prefix/" + str(i)
            self.test_bucket.put_object(key)

        options = MockOptions(
            qs_path="qs://" + test_bucket1 + "/prefix/",
            recursive=True,
            exclude=None,
            include=None)
        self.Rm.send_request(options)

    def test_remove_mutiple_keys_2(self):
        for i in range(0, 10):
            key = "prefix/" + str(i) + ".txt"
            self.test_bucket.put_object(key)
        self.test_bucket.put_object("prefix/test.jpg")

        options = MockOptions(
            qs_path="qs://" + test_bucket1 + "/prefix/",
            recursive=True,
            exclude="*.txt",
            include=None)
        self.Rm.send_request(options)

    def test_remove_mutiple_keys_3(self):
        for i in range(0, 10):
            key = "prefix/" + str(i) + ".txt"
            self.test_bucket.put_object(key)
        self.test_bucket.put_object("prefix/test.jpg")

        options = MockOptions(
            qs_path="qs://" + test_bucket1 + "/prefix/",
            recursive=True,
            exclude="*",
            include="*.txt")
        self.Rm.send_request(options)

    def tearDown(self):
        options = MockOptions(exclude=None, include=None)
        self.Rm.remove_multiple_keys(test_bucket1, options=options)


if __name__ == "__main__":
    unittest.main()
