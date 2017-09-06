# coding:utf-8

from sys import version_info
from setuptools import setup, find_packages
from qingstor.qsctl import __version__

install_requires = [
    'PyYAML >= 3.1',
    'qingstor-sdk >= 2.1.3',
    'docutils >= 0.10',
    'tqdm >= 4.0.0'
]

# python version under 2.7.9 should use requests[security]
if version_info[:3] < (2, 7, 9):
    install_requires.append("requests[security]")

# python version under 2.7 do not have argparse
if version_info[:3] < (2, 7):
    install_requires.append("argparse >= 1.1")

# python version under 3 do not have concurrent.futures
if version_info[:3] < (3, 2):
    install_requires.append("futures")

setup(
    name='qsctl',
    version=__version__,
    description='Advanced command line tool for QingStor.',
    long_description=open('README.rst', 'rb').read().decode('utf-8'),
    keywords='yunify qingcloud qingstor qsctl object_storage',
    author='QingStor Dev Team',
    author_email='qs-devel@yunify.com',
    url='https://www.qingstor.com',
    scripts=['bin/qsctl', 'bin/qsctl.cmd',
             'bin/qsctl_completer', 'bin/qsctl_completion_bash'],
    packages=find_packages('.'),
    package_dir={'qsctl': 'qingstor'},
    namespace_packages=['qingstor'],
    include_package_data=True,
    install_requires=install_requires
)
