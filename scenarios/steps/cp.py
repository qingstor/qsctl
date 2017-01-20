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


@given(u'a set of local files')
def step_impl(context):
    sh.mkdir("tmp").wait()
    for row in context.table:
        sh.dd(
            "if=/dev/zero", "of=tmp/" + row["name"], "bs=1048576",
            "count=" + row["count"]
        )


@when(u'copy to QingStor key')
def step_impl(context):
    for row in context.table:
        command = qsctl(
            "cp",
            "tmp/{filename}".format(filename=row["name"]),
            "qs://{bucket}/{filename}".format(
                bucket=test_data['bucket_name'], filename=row["name"]
            )
        )
        command.wait()


@then(u'QingStor should have key')
def step_impl(context):
    resp = bucket.list_objects()
    assert_that(sorted([i["key"] for i in resp["keys"]])
                ).is_equal_to(sorted([row["name"] for row in context.table]))

    for row in context.table:
        bucket.delete_object(row["name"])


@when(u'copy to QingStor keys recursively')
def step_impl(context):
    command = qsctl(
        "cp",
        "tmp",
        "qs://{bucket}".format(bucket=test_data['bucket_name']),
        "-r"
    )
    command.wait()


@then(u'QingStor should have keys')
def step_impl(context):
    resp = bucket.list_objects()
    assert_that(sorted([i["key"] for i in resp["keys"]])
                ).is_equal_to(sorted([row["name"] for row in context.table]))

    sh.rm("-rf", "tmp")


@when(u'copy to local file')
def step_impl(context):
    for row in context.table:
        command = qsctl(
            "cp",
            "qs://{bucket}/{filename}".format(
                bucket=test_data['bucket_name'], filename=row["name"]
            ),
            "tmp/{filename}".format(filename=row["name"]),
        )
        command.wait()


@then(u'local should have file')
def step_impl(context):
    output = sh.ls("tmp").stdout.decode("utf-8")
    ok = True
    for row in context.table:
        if row["name"] not in output:
            ok = False
            break
    assert_that(ok).is_equal_to(True)

    sh.rm("-rf", "tmp")


@when(u'copy to local files recursively')
def step_impl(context):
    command = qsctl(
        "cp",
        "qs://{bucket}".format(
            bucket=test_data["bucket_name"],
        ),
        "tmp",
        "-r",
    )
    command.wait()


@then(u'local should have files')
def step_impl(context):
    output = sh.ls("tmp").stdout.decode("utf-8")
    ok = True
    for row in context.table:
        if row["name"] not in output:
            ok = False
            break
    assert_that(ok).is_equal_to(True)

    sh.rm("-rf", "tmp")

    for row in context.table:
        bucket.delete_object(row["name"])
