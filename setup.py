# coding:utf-8

from setuptools import setup, find_packages

setup(
    name='qsctl',
    version='1.4.0',
    description='Advanced command line tool for QingStor.',
    long_description=open('README.rst', 'rb').read().decode('utf-8'),
    keywords='yunify qingcloud qingstor qsctl object_storage',
    author='QingStor Dev Team',
    author_email='qs_dev_group@yunify.com ',
    url='https://www.qingstor.com',
    scripts=['bin/qsctl', 'bin/qsctl.cmd'],
    packages=find_packages('.'),
    package_dir={'qsctl': 'qingstor'},
    namespace_packages=['qingstor'],
    include_package_data=True,
    install_requires=[
        'argparse >= 1.1',
        'PyYAML >= 3.1',
        'qingstor-sdk >= 2.1.0',
        'docutils >= 0.10',
        'tqdm >= 4.0.0'
    ])
