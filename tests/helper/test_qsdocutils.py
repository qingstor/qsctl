import os
import unittest

from qingstor.qsctl.helper.qsdocutils import (
    to_rst_style_title, gen_rst_doc, gen_sphinx_doc, RstDocument
)

RST_STYLE_TITLE = "=========\nTESTTITLE\n========="


class TestQsdocutils(unittest.TestCase):

    def setUp(self):
        # Create a tmp file
        with open("tmp.file", 'w') as f:
            f.write("just for testing")
        self.tmp_file = "tmp.file"

    def test_to_rst_style_title(self):
        self.assertEqual(to_rst_style_title("TESTTITLE"), RST_STYLE_TITLE)

    def test_RstDocument(self):
        rst_doc = RstDocument()

        rst_doc.from_file(self.tmp_file)
        rst_doc.add_reporting_bug()
        rst_doc.add_see_also("qsctl-ls")
        rst_doc.add_copyright()
        print(rst_doc.getvalue())

    def test_gen_rst_doc(self):
        gen_rst_doc("qsctl-ls")

    def test_gen_sphinx_doc(self):
        gen_sphinx_doc("qsctl-ls")

    def tearDown(self):
        if os.path.exists(self.tmp_file):
            os.remove(self.tmp_file)


if __name__ == "__main__":
    unittest.main(verbosity=2)
