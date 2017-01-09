import unittest

from mock import MockOptions

from tests.test_data import zone, test_bucket1, test_bucket2

from qingstor.qsctl.commands.ls import LsCommand
from qingstor.qsctl.utils import load_conf


class TestLsCommand(unittest.TestCase):
    Ls = LsCommand

    def setUp(self):

        # Set the http connection
        conf = load_conf("~/.qingstor/config.yaml")
        options = MockOptions(qs_path="qs://")
        self.Ls.client = self.Ls.get_client(conf)

        self.test_bucket1 = self.Ls.client.Bucket(test_bucket1, zone)
        self.test_bucket2 = self.Ls.client.Bucket(test_bucket2, zone)
        self.test_bucket1.put()
        resp = self.test_bucket1.head()
        if resp.status_code != 200:
            self.fail("setUp failed: please use another bucket name")
        self.test_bucket2.put()
        resp = self.test_bucket2.head()
        if resp.status_code != 200:
            self.fail("setUp failed: please use another bucket name")

    def test_list_buckets_1(self):
        options = MockOptions(zone="pek3a")
        self.Ls.list_buckets(options)

    def test_list_buckets_2(self):
        options = MockOptions(zone=None)
        self.Ls.list_buckets(options)

    def test_list_keys_1(self):
        # Set the http connection for listing keys
        options = MockOptions(
            qs_path="qs://" + test_bucket1, recursive=False, page_size=None)
        conf = load_conf("~/.qingstor/config.yaml")
        self.Ls.client = self.Ls.get_client(conf)

        # We need some keys for testing.
        for i in range(0, 10):
            self.test_bucket1.put_object("test" + str(i))

        # testing
        self.Ls.list_keys(options)

        # clean keys after testing
        options = MockOptions(exclude=None, include=None)
        self.Ls.remove_multiple_keys(test_bucket1, options=options)

    def test_list_keys_2(self):
        # Set the http connection for listing objects
        options = MockOptions(
            qs_path="qs://" + test_bucket1 + "/prefix/",
            recursive=True,
            page_size=None)
        conf = load_conf("~/.qingstor/config.yaml")
        self.Ls.client = self.Ls.get_client(conf)

        # We need some keys for testing.
        self.test_bucket1.put_object("test.txt")
        for i in range(0, 10):
            self.test_bucket1.put_object("prefix/" + str(i))

        # testing
        self.Ls.list_keys(options)

        # clean keys after testing
        options = MockOptions(exclude=None, include=None)
        self.Ls.remove_multiple_keys(test_bucket1, options=options)

    def test_list_keys_3(self):
        # Set the http connection for listing objects
        options = MockOptions(
            qs_path="qs://" + test_bucket1, recursive=False, page_size=None)
        conf = load_conf("~/.qingstor/config.yaml")
        self.Ls.client = self.Ls.get_client(conf)

        # We need some keys for testing.
        self.test_bucket1.put_object("test.txt")
        for i in range(0, 10):
            self.test_bucket1.put_object("prefix/" + str(i))

        # testing
        self.Ls.list_keys(options)

        # clean keys after testing
        options = MockOptions(exclude=None, include=None)
        self.Ls.remove_multiple_keys(test_bucket1, options=options)


if __name__ == "__main__":
    unittest.main(verbosity=2)
