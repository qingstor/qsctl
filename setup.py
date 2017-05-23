# coding:utf-8

from sys import version_info
from setuptools import setup, find_packages
from qingstor.qsctl import __version__

install_requires = [
    'argparse >= 1.1',
    'PyYAML >= 3.1',
    'qingstor-sdk >= 2.1.0',
    'docutils >= 0.10',
    'tqdm >= 4.0.0'
]

if version_info[:3] < (2, 7, 9):
    install_requires.append("requests[security]")

setup(
    name='qsctl',
    version=__version__,
    description='Advanced command line tool for QingStor.',
    long_description=open('README.rst', 'rb').read().decode('utf-8'),
    keywords='yunify qingcloud qingstor qsctl object_storage',
    author='QingStor Dev Team',
    author_email='qs-devel@yunify.com',
    url='https://www.qingstor.com',
    scripts=['bin/qsctl', 'bin/qsctl.cmd'],
    packages=find_packages('.'),
    package_dir={'qsctl': 'qingstor'},
    namespace_packages=['qingstor'],
    include_package_data=True,
    install_requires=install_requires
)
