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


@given(u'serveral QingStor keys')
def step_impl(context):
    context.input = context.table
    for row in context.input:
        bucket.put_object(row["name"])


@when(u'delete keys')
def step_impl(context):
    for row in context.input:
        qsctl(
            "rm",
            "qs://{bucket}/{filename}".format(
                bucket=test_data["bucket_name"],
                filename=row["name"],
            )
        ).wait()


@then(u'QingStor keys should be deleted')
def step_impl(context):
    for row in context.input:
        assert_that(bucket.head_object(row["name"]).status_code
                    ).is_equal_to(404)


@given(u'serveral QingStor keys with prefix')
def step_impl(context):
    for row in context.table:
        bucket.put_object(row["name"])


@when(u'delete keys with prefix "中文目录测试"')
def step_impl(context):
    qsctl(
        "rm",
        "qs://{bucket}/{prefix}".format(
            bucket=test_data['bucket_name'],
            prefix=context.text,
        ),
        "-r",
    ).wait()


@then(u'QingStor keys with prefix should be deleted, other files should keep')
def step_impl(context):
    for row in context.table:
        assert_that(bucket.head_object(row["name"]).status_code
                    ).is_equal_to(404 if row["deleted"] is "1" else 200)

    for row in context.table:
        bucket.delete_object(row["name"])
