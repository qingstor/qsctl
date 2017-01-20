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
bucket = qingstor.Bucket(test_data["bucket_name"], test_data["zone"])
bucket.put()
qsctl = sh.Command("qsctl")


@when(u'list all buckets')
def step_impl(context):
    context.output = qsctl("ls").stdout.decode("utf-8")


@then(u'should have this bucket')
def step_impl(context):
    assert_that(test_data['bucket_name'] in context.output).is_equal_to(True)


@given(u'a bucket with files')
def step_impl(context):
    for row in context.table:
        bucket.put_object(row["name"])


@when(u'list keys')
def step_impl(context):
    context.output = qsctl(
        "ls", "qs://{bucket}".format(bucket=test_data["bucket_name"])
    ).stdout.decode("utf-8")


@then(u'should list all keys')
def step_impl(context):
    ok = True
    for row in context.table:
        if row["name"] not in context.output:
            ok = False
            break
    assert_that(ok).is_equal_to(True)


@when(u'list keys with prefix')
def step_impl(context):
    context.output = {}
    for row in context.table:
        context.output[row["prefix"]] = qsctl(
            "ls",
            "qs://{bucket}/{prefix}".format(
                bucket=test_data["bucket_name"], prefix=row["prefix"]
            )
        ).stdout.decode("utf-8")


@then(u'should list keys with prefix')
def step_impl(context):
    for row in context.table:
        assert_that(row["should_show_up"] in context.output[row["prefix"]]
                    ).is_equal_to(True)
        assert_that(row["not_show_up"] in context.output[row["prefix"]]
                    ).is_equal_to(False)


@when(u'list keys recursively')
def step_impl(context):
    context.output = qsctl(
        "ls", "qs://{bucket}".format(bucket=test_data["bucket_name"]), "-r"
    ).stdout.decode("utf-8")


@Then(u'should list keys recursively')
def step_impl(context):
    for row in context.table:
        assert_that(row["name"] in context.output).is_equal_to(True)

    for row in context.table:
        bucket.delete_object(row["name"])
