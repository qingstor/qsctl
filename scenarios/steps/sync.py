# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import os
import re

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


@given(u'a local directory with files')
def step_impl(context):
    sh.mkdir("tmp").wait()
    for row in context.table:
        dirs = row["name"].split("/")
        if len(dirs) > 1:
            sh.mkdir("tmp/" + dirs[0]).wait()
        sh.dd(
            "if=/dev/zero",
            "of=tmp/{filename}".format(filename=row["name"]),
            "bs=1048576",
            "count=1"
        ).wait()


@when(u'sync local directory to QingStor prefix')
def step_impl(context):
    qsctl(
        "sync",
        "tmp",
        "qs://{bucket}/{prefix}".format(
            bucket=test_data['bucket_name'],
            prefix="tmp",
        ),
    ).wait()


@then(u'QingStor should have keys with prefix')
def step_impl(context):
    for row in context.table:
        assert_that(bucket.head_object(row["name"]).status_code
                    ).is_equal_to(200)

    sh.rm("-rf", "tmp").wait()


@when(u'sync QingStor prefix to local directory')
def step_impl(context):
    qsctl(
        "sync",
        "qs://{bucket}/{prefix}".format(
            bucket=test_data['bucket_name'],
            prefix="tmp",
        ),
        "tmp",
    ).wait()


@then(u'local should have files with prefix')
def step_impl(context):
    for row in context.table:
        assert_that(os.path.isfile(row["name"])).is_equal_to(True)

    for row in context.table:
        bucket.delete_object(row["name"])

    sh.rm("-rf", "tmp").wait()


@given(u'several similar local directories with files')
def step_impl(context):
    sh.mkdir("tmp_similar").wait()
    for row in context.table:
        dirs = row["name"].split("/")
        if len(dirs) > 2:
            sh.mkdir("tmp_similar/" + dirs[1]).wait()
        sh.dd(
            "if=/dev/zero",
            "of={filename}".format(filename=row["name"]),
            "bs=1048576",
            "count=1"
        ).wait()


@when(u'sync local directories to QingStor prefix using wildcard')
def step_impl(context):
    pattern = context.table[0][1]
    pattern = pattern.replace("*", ".*")
    pattern = pattern.replace("?", "[\S]{1}")
    pattern = re.compile(pattern)
    for row in context.table:
        result = re.findall(pattern, row["name"])
        for x in result:
            if x == row["name"]:
                qsctl(
                    "sync",
                    "tmp_similar/{filename}".format(filename=row["name"]),
                    "qs://{bucket}/{prefix}".format(
                        bucket=test_data["bucket_name"],
                        prefix="tmp_similar")
                ).wait()


@then(u'QingStor should have keys with matched prefix')
def step_impl(context):
    for row in context.table:
        assert_that(os.path.isfile(row["name"])).is_equal_to(True)

    for row in context.table:
        bucket.delete_object(row["name"])

    sh.rm("-rf", "tmp_similar").wait()
