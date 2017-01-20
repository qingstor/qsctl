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


@given(u'several local files')
def step_impl(context):
    context.input = context.table
    sh.mkdir("tmp").wait()
    for row in context.input:
        dirs = row["name"].split("/")
        if len(dirs) > 1:
            sh.mkdir("tmp/" + dirs[0]).wait()
        sh.dd(
            "if=/dev/zero",
            "of=tmp/{filename}".format(filename=row["name"]),
            "bs=1048576",
            "count=1"
        ).wait()


@when(u'move to QingStor')
def step_impl(context):
    for row in context.input:
        qsctl(
            "mv",
            "tmp/{filename}".format(filename=row["name"]),
            "qs://{bucket}/{filename}".format(
                bucket=test_data["bucket_name"],
                filename=row["name"],
            )
        ).wait()


@then(u'QingStor should have same file and local files should be deleted')
def step_impl(context):
    for row in context.input:
        assert_that(bucket.head_object(row["name"]).status_code
                    ).is_equal_to(200)
        assert_that(os.path.isfile("tmp/" + row["name"])).is_equal_to(False)


@given(u'several QingStor keys')
def step_impl(context):
    context.input = context.table


@when(u'move to local')
def step_impl(context):
    for row in context.input:
        qsctl(
            "mv",
            "qs://{bucket}/{filename}".format(
                bucket=test_data["bucket_name"],
                filename=row["name"],
            ),
            "tmp/{filename}".format(
                filename=row["name"],
            )
        ).wait()


@then(u'local should have same file and QingStor keys should be deleted')
def step_impl(context):
    for row in context.input:
        assert_that(bucket.head_object(row["name"]).status_code
                    ).is_equal_to(404)
        assert_that(os.path.isfile("tmp/" + row["name"])).is_equal_to(True)

    sh.rm("-rf", "tmp").wait()
