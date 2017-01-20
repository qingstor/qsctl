# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import os

import sh
import yaml
from behave import *
from assertpy import assert_that

from qingstor.sdk.config import Config
from qingstor.sdk.service.qingstor import QingStor

config = Config().load_user_config()
qingstor = QingStor(config)
test_config_file_path = os.path.abspath(
    os.path.join(os.path.dirname(__file__), os.path.pardir)
)
with open(test_config_file_path + '/test_config.yaml') as f:
    test_data = yaml.load(f)
    f.close()
bucket = qingstor.Bucket(test_data["bucket_name"] + "1", test_data["zone"])
qsctl = sh.Command("qsctl")


@when(u'create a bucket')
def step_impl(context):
    qsctl("mb", "qs://{bucket}".format(bucket=test_data["bucket_name"] + "1"))


@then(u'should get a new bucket')
def step_impl(context):
    assert_that(bucket.head().status_code).is_equal_to(201)
