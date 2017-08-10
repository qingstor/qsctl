# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import os

import sh
import yaml

from behave import *
from assertpy import assert_that
from requests.utils import urlparse

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


@given(u'a non-public bucket with objects')
def step_impl(context):
    context.input = context.table
    bucket.put_object(context.input[0]["name"])
    current_acl = bucket.get_acl()
    global is_changed
    is_changed = False  # Mark the initial state
    is_private = True  # Store the initial permission
    for v in current_acl["acl"]:
        if v["grantee"]["type"] == "group":
            is_private = False
    if is_private is False:
        for x in current_acl["acl"]:
            if x["grantee"]["type"] == "group":
                current_acl["acl"].remove(x)  # Delete the group permission
                bucket.put_acl(current_acl["acl"])
                is_changed = True


@when(u'generate the link with signature')
def step_impl(context):
    context.output = qsctl(
        "presign",
        "qs://{bucket}/{filename}".format(
            bucket=test_data["bucket_name"],
            filename=context.input[0]["name"],
        ),
        "-e %d" % int(context.input[0]["expire_seconds"])
    ).stdout.decode("utf-8")


@then(u'the link should include parameters like signature, etc.')
def step_impl(context):
    if is_changed is True:
        bucket.put_acl([{
            "grantee": {
                "type": "group",
                "name": "QS_ALL_USERS"
            },
            "permission": "READ"
        }])
    result = urlparse(context.output)
    params = result.query
    object_key = result.path
    assert_that(params).contains("signature")
    assert_that(params).contains("access_key_id")
    assert_that(params).contains("expires")
    assert_that(object_key.split("/")[-1]).is_equal_to(context.input[0]["name"])


@given(u'a public bucket with objects')
def step_impl(context):
    context.input = context.table
    bucket.put_object(context.input[0]["name"])
    current_acl = bucket.get_acl()
    global is_changed
    is_changed = False  # Mark the initial state
    is_private = True  # Store the initial permission
    for v in current_acl["acl"]:
        if v["grantee"]["type"] == "group":
            is_private = False
    if is_private is True:
        bucket.put_acl([{
            "grantee": {
                "type": "group",
                "name": "QS_ALL_USERS"
            },
            "permission": "READ"
        }])
        is_changed = True


@when(u'generate the spliced link')
def step_impl(context):
    context.output = qsctl(
        "presign",
        "qs://{bucket}/{filename}".format(
            bucket=test_data["bucket_name"],
            filename=context.input[0]["name"],
        )
    ).stdout.decode("utf-8")


@then(u'the link should include object_key')
def step_impl(context):
    current_acl = bucket.get_acl()
    if is_changed is True:
        for x in current_acl["acl"]:
            if x["grantee"]["type"] == "group":
                current_acl["acl"].remove(x)
                bucket.put_acl(current_acl["acl"])
    result = urlparse(context.output)
    object_key = result.path
    assert_that(object_key.split("/")[-1].strip()
                ).is_equal_to(context.input[0]["name"])
