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
bucket = qingstor.Bucket(test_data['bucket_name'], test_data['zone'])
bucket.put()
qsctl = sh.Command("qsctl")


@given(u'an empty QingStor bucket')
def step_impl(context):
    context.empty_bucket = qingstor.Bucket(
        test_data['bucket_name'] + "1", test_data['zone']
    )
    context.empty_bucket.put()


@when(u'delete bucket')
def step_impl(context):
    qsctl("rb", "qs://{bucket}".format(bucket=test_data['bucket_name'] + "1"))


@then(u'the bucket should be deleted')
def step_impl(context):
    assert_that(context.empty_bucket.head().status_code).is_equal_to(404)


@given(u'a QingSto bucket with files')
def step_impl(context):
    context.nonempty_bucket = qingstor.Bucket(
        test_data['bucket_name'] + "2", test_data['zone']
    )
    context.nonempty_bucket.put()
    context.nonempty_bucket.put("test_file")


@when(u'delete bucket forcibly')
def step_impl(context):
    qsctl(
        "rb",
        "qs://{bucket}".format(bucket=test_data['bucket_name'] + "2"),
        "-f"
    )


@then(u'the bucket should be deleted forcibly')
def step_impl(context):
    assert_that(context.nonempty_bucket.head().status_code).is_equal_to(404)
