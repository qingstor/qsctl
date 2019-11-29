package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// initEnUS will init en_US support.
func initEnUS(tag language.Tag) {
	_ = message.SetString(tag, "%s  %s  %s  %s\n", "%s  %s  %s  %s\n")
	_ = message.SetString(tag, "-r is required to delete a directory", "-r is required to delete a directory")
	_ = message.SetString(tag, "AccessKey and SecretKey not found. Please setup your config now, or exit and setup manually.", "AccessKey and SecretKey not found. Please setup your config now, or exit and setup manually.")
	_ = message.SetString(tag, "Bucket <%s> created.\n", "Bucket <%s> created.\n")
	_ = message.SetString(tag, "Bucket <%s> removed.\n", "Bucket <%s> removed.\n")
	_ = message.SetString(tag, "Cat object: qsctl cat qs://prefix/a", "Cat object: qsctl cat qs://prefix/a")
	_ = message.SetString(tag, "Config not loaded, use default and environment value instead.", "Config not loaded, use default and environment value instead.")
	_ = message.SetString(tag, "Copy file: qsctl cp /path/to/file qs://prefix/a", "Copy file: qsctl cp /path/to/file qs://prefix/a")
	_ = message.SetString(tag, "Copy folder: qsctl cp qs://prefix/a /path/to/folder -r", "Copy folder: qsctl cp qs://prefix/a /path/to/folder -r")
	_ = message.SetString(tag, "Delete an empty qingstor bucket or forcely delete nonempty qingstor bucket.", "Delete an empty qingstor bucket or forcely delete nonempty qingstor bucket.")
	_ = message.SetString(tag, "Dir <%s> and <%s> synced.\n", "Dir <%s> and <%s> synced.\n")
	_ = message.SetString(tag, "Dir <%s> removed.\n", "Dir <%s> removed.\n")
	_ = message.SetString(tag, "Key <%s> copied.\n", "Key <%s> copied.\n")
	_ = message.SetString(tag, "Key <%s> moved.\n", "Key <%s> moved.\n")
	_ = message.SetString(tag, "List bucket's all objects: qsctl ls qs://bucket-name -R", "List bucket's all objects: qsctl ls qs://bucket-name -R")
	_ = message.SetString(tag, "List buckets: qsctl ls", "List buckets: qsctl ls")
	_ = message.SetString(tag, "List objects by long format: qsctl ls qs://bucket-name -l", "List objects by long format: qsctl ls qs://bucket-name -l")
	_ = message.SetString(tag, "List objects with prefix: qsctl ls qs://bucket-name/prefix", "List objects with prefix: qsctl ls qs://bucket-name/prefix")
	_ = message.SetString(tag, "Load config failed [%v]", "Load config failed [%v]")
	_ = message.SetString(tag, "Make bucket: qsctl mb bucket-name", "Make bucket: qsctl mb bucket-name")
	_ = message.SetString(tag, "Move file: qsctl mv /path/to/file qs://prefix/a", "Move file: qsctl mv /path/to/file qs://prefix/a")
	_ = message.SetString(tag, "Move folder: qsctl mv qs://prefix/a /path/to/folder -r", "Move folder: qsctl mv qs://prefix/a /path/to/folder -r")
	_ = message.SetString(tag, "Object <%s> removed.\n", "Object <%s> removed.\n")
	_ = message.SetString(tag, "Presign object: qsctl qs://bucket-name/object-name", "Presign object: qsctl qs://bucket-name/object-name")
	_ = message.SetString(tag, "Read from stdin: cat /path/to/file | qsctl cp - qs://prefix/stdin", "Read from stdin: cat /path/to/file | qsctl cp - qs://prefix/stdin")
	_ = message.SetString(tag, "Remove a single object: qsctl rm qs://bucket-name/object-key", "Remove a single object: qsctl rm qs://bucket-name/object-key")
	_ = message.SetString(tag, "Stat object: qsctl stat qs://prefix/a", "Stat object: qsctl stat qs://prefix/a")
	_ = message.SetString(tag, "Sync QS-Directory to local directory: qsctl sync qs://bucket-name/test/ test_local/", "Sync QS-Directory to local directory: qsctl sync qs://bucket-name/test/ test_local/")
	_ = message.SetString(tag, "Sync local directory to QS-Directory: qsctl sync . qs://bucket-name", "Sync local directory to QS-Directory: qsctl sync . qs://bucket-name")
	_ = message.SetString(tag, "Sync skip updating files that already exist on receiver: qsctl sync . qs://bucket-name --ignore-existing", "Sync skip updating files that already exist on receiver: qsctl sync . qs://bucket-name --ignore-existing")
	_ = message.SetString(tag, "Tee object: qsctl tee qs://prefix/a", "Tee object: qsctl tee qs://prefix/a")
	_ = message.SetString(tag, "Write to stdout: qsctl cp qs://prefix/b - > /path/to/file", "Write to stdout: qsctl cp qs://prefix/b - > /path/to/file")
	_ = message.SetString(tag, "Your config has been set to <%v>. You can still modify it manually.", "Your config has been set to <%v>. You can still modify it manually.")
	_ = message.SetString(tag, "assign config path manually", "assign config path manually")
	_ = message.SetString(tag, "cat a remote object to stdout", "cat a remote object to stdout")
	_ = message.SetString(tag, "cat qs://<bucket_name>/<object_key>", "cat qs://<bucket_name>/<object_key>")
	_ = message.SetString(tag, "copy directory recursively", "copy directory recursively")
	_ = message.SetString(tag, "copy from/to qingstor", "copy from/to qingstor")
	_ = message.SetString(tag, "cp <source-path> <dest-path>", "cp <source-path> <dest-path>")
	_ = message.SetString(tag, "delete a bucket", "delete a bucket")
	_ = message.SetString(tag, "delete an empty bucket: qsctl rb qs://bucket-name", "delete an empty bucket: qsctl rb qs://bucket-name")
	_ = message.SetString(tag, "enable benchmark or not", "enable benchmark or not")
	_ = message.SetString(tag, "forcely delete a nonempty bucket: qsctl rb qs://bucket-name -f", "forcely delete a nonempty bucket: qsctl rb qs://bucket-name -f")
	_ = message.SetString(tag, "get the pre-signed URL by the object key", "get the pre-signed URL by the object key")
	_ = message.SetString(tag, "help for this command", "help for this command")
	_ = message.SetString(tag, "in which zone to do the operation", "in which zone to do the operation")
	_ = message.SetString(tag, "in which zone to make the bucket (required)", "in which zone to make the bucket (required)")
	_ = message.SetString(tag, "list objects or buckets", "list objects or buckets")
	_ = message.SetString(tag, "ls [qs://<bucket-name/prefix>]", "ls [qs://<bucket-name/prefix>]")
	_ = message.SetString(tag, "make a new bucket", "make a new bucket")
	_ = message.SetString(tag, "move directory recursively", "move directory recursively")
	_ = message.SetString(tag, "move from/to qingstor", "move from/to qingstor")
	_ = message.SetString(tag, "print logs for debug", "print logs for debug")
	_ = message.SetString(tag, "qsctl cat can cat a remote object to stdout", "qsctl cat can cat a remote object to stdout")
	_ = message.SetString(tag, "qsctl cp can copy file/folder/stdin to qingstor or copy qingstor objects to local/stdout", "qsctl cp can copy file/folder/stdin to qingstor or copy qingstor objects to local/stdout")
	_ = message.SetString(tag, "qsctl mv can move file/folder to qingstor or move qingstor objects to local", "qsctl mv can move file/folder to qingstor or move qingstor objects to local")
	_ = message.SetString(tag, "qsctl rb delete a qingstor bucket", "qsctl rb delete a qingstor bucket")
	_ = message.SetString(tag, "qsctl rm remove the object with given object key", "qsctl rm remove the object with given object key")
	_ = message.SetString(tag, "qsctl stat show the detailed info of this object", "qsctl stat show the detailed info of this object")
	_ = message.SetString(tag, "recursively delete keys under a specific prefix", "recursively delete keys under a specific prefix")
	_ = message.SetString(tag, "recursively list subdirectories encountered", "recursively list subdirectories encountered")
	_ = message.SetString(tag, "remove a remote object", "remove a remote object")
	_ = message.SetString(tag, "stat a remote object", "stat a remote object")
	_ = message.SetString(tag, "sync between local directory and QS-Directory", "sync between local directory and QS-Directory")
	_ = message.SetString(tag, "tee a remote object from stdin", "tee a remote object from stdin")
	_ = message.SetString(tag, "the number of seconds until the pre-signed URL expires. Default is 300 seconds", "the number of seconds until the pre-signed URL expires. Default is 300 seconds")
	_ = message.SetString(tag, `qsctl ls can list all qingstor buckets or qingstor keys under a prefix.`, `qsctl ls can list all qingstor buckets or qingstor keys under a prefix.`)
	_ = message.SetString(tag, `qsctl mb can make a new bucket with the specific name,

bucket name should follow DNS name rule with:
* length between 6 and 63;
* can only contains lowercase letters, numbers and hyphen -
* must start and end with lowercase letter or number
* must not be an available IP address
	`, `qsctl mb can make a new bucket with the specific name,

bucket name should follow DNS name rule with:
* length between 6 and 63;
* can only contains lowercase letters, numbers and hyphen -
* must start and end with lowercase letter or number
* must not be an available IP address
	`)
	_ = message.SetString(tag, `qsctl presign can generate a pre-signed URL for the object.
Within the given expire time, anyone who receives this URL can retrieve
the object with an HTTP GET request. If an object belongs to a public bucket,
generate a URL spliced by bucket name, zone and its name, anyone who receives
this URL can always retrieve the object with an HTTP GET request.`, `qsctl presign can generate a pre-signed URL for the object.
Within the given expire time, anyone who receives this URL can retrieve
the object with an HTTP GET request. If an object belongs to a public bucket,
generate a URL spliced by bucket name, zone and its name, anyone who receives
this URL can always retrieve the object with an HTTP GET request.`)
	_ = message.SetString(tag, `qsctl sync between local directory and QS-Directory. The first path argument
is the source directory and second the destination directory.

When a key(file) already exists in the destination directory, program will compare
the modified time of source file(key) and destination key(file). The destination
key(file) will be overwritten only if the source one newer than destination one.`, `qsctl sync between local directory and QS-Directory. The first path argument
is the source directory and second the destination directory.

When a key(file) already exists in the destination directory, program will compare
the modified time of source file(key) and destination key(file). The destination
key(file) will be overwritten only if the source one newer than destination one.`)
	_ = message.SetString(tag, `qsctl tee can tee a remote object from stdin.

NOTICE: qsctl will not tee the content to stdout like linux tee command does.
`, `qsctl tee can tee a remote object from stdin.

NOTICE: qsctl will not tee the content to stdout like linux tee command does.
`)
	_ = message.SetString(tag, `skip creating new files in dest dirs, only copy newer by time`, `skip creating new files in dest dirs, only copy newer by time`)
	_ = message.SetString(tag, `use the specified FORMAT instead of the default;
output a newline after each use of FORMAT

The valid format sequences for files:

  %F   file type
  %h   content md5 of the file
  %n   file name
  %s   total size, in bytes
  %y   time of last data modification, human-readable, e.g: 2006-01-02 15:04:05 +0000 UTC
  %Y   time of last data modification, seconds since Epoch
	`, `use the specified FORMAT instead of the default;
output a newline after each use of FORMAT

The valid format sequences for files:

  %F   file type
  %h   content md5 of the file
  %n   file name
  %s   total size, in bytes
  %y   time of last data modification, human-readable, e.g: 2006-01-02 15:04:05 +0000 UTC
  %Y   time of last data modification, seconds since Epoch
	`)
}

// initZhCN will init zh_CN support.
func initZhCN(tag language.Tag) {
	_ = message.SetString(tag, "%s  %s  %s  %s\n", "%s  %s  %s  %s\n")
	_ = message.SetString(tag, "-r is required to delete a directory", "-r is required to delete a directory")
	_ = message.SetString(tag, "AccessKey and SecretKey not found. Please setup your config now, or exit and setup manually.", "AccessKey and SecretKey not found. Please setup your config now, or exit and setup manually.")
	_ = message.SetString(tag, "Bucket <%s> created.\n", "Bucket <%s> created.\n")
	_ = message.SetString(tag, "Bucket <%s> removed.\n", "Bucket <%s> removed.\n")
	_ = message.SetString(tag, "Cat object: qsctl cat qs://prefix/a", "Cat object: qsctl cat qs://prefix/a")
	_ = message.SetString(tag, "Config not loaded, use default and environment value instead.", "Config not loaded, use default and environment value instead.")
	_ = message.SetString(tag, "Copy file: qsctl cp /path/to/file qs://prefix/a", "Copy file: qsctl cp /path/to/file qs://prefix/a")
	_ = message.SetString(tag, "Copy folder: qsctl cp qs://prefix/a /path/to/folder -r", "Copy folder: qsctl cp qs://prefix/a /path/to/folder -r")
	_ = message.SetString(tag, "Delete an empty qingstor bucket or forcely delete nonempty qingstor bucket.", "Delete an empty qingstor bucket or forcely delete nonempty qingstor bucket.")
	_ = message.SetString(tag, "Dir <%s> and <%s> synced.\n", "Dir <%s> and <%s> synced.\n")
	_ = message.SetString(tag, "Dir <%s> removed.\n", "Dir <%s> removed.\n")
	_ = message.SetString(tag, "Key <%s> copied.\n", "文件 <%s> 已复制\n")
	_ = message.SetString(tag, "Key <%s> moved.\n", "Key <%s> moved.\n")
	_ = message.SetString(tag, "List bucket's all objects: qsctl ls qs://bucket-name -R", "List bucket's all objects: qsctl ls qs://bucket-name -R")
	_ = message.SetString(tag, "List buckets: qsctl ls", "List buckets: qsctl ls")
	_ = message.SetString(tag, "List objects by long format: qsctl ls qs://bucket-name -l", "List objects by long format: qsctl ls qs://bucket-name -l")
	_ = message.SetString(tag, "List objects with prefix: qsctl ls qs://bucket-name/prefix", "List objects with prefix: qsctl ls qs://bucket-name/prefix")
	_ = message.SetString(tag, "Load config failed [%v]", "Load config failed [%v]")
	_ = message.SetString(tag, "Make bucket: qsctl mb bucket-name", "Make bucket: qsctl mb bucket-name")
	_ = message.SetString(tag, "Move file: qsctl mv /path/to/file qs://prefix/a", "Move file: qsctl mv /path/to/file qs://prefix/a")
	_ = message.SetString(tag, "Move folder: qsctl mv qs://prefix/a /path/to/folder -r", "Move folder: qsctl mv qs://prefix/a /path/to/folder -r")
	_ = message.SetString(tag, "Object <%s> removed.\n", "Object <%s> removed.\n")
	_ = message.SetString(tag, "Presign object: qsctl qs://bucket-name/object-name", "Presign object: qsctl qs://bucket-name/object-name")
	_ = message.SetString(tag, "Read from stdin: cat /path/to/file | qsctl cp - qs://prefix/stdin", "Read from stdin: cat /path/to/file | qsctl cp - qs://prefix/stdin")
	_ = message.SetString(tag, "Remove a single object: qsctl rm qs://bucket-name/object-key", "Remove a single object: qsctl rm qs://bucket-name/object-key")
	_ = message.SetString(tag, "Stat object: qsctl stat qs://prefix/a", "Stat object: qsctl stat qs://prefix/a")
	_ = message.SetString(tag, "Sync QS-Directory to local directory: qsctl sync qs://bucket-name/test/ test_local/", "Sync QS-Directory to local directory: qsctl sync qs://bucket-name/test/ test_local/")
	_ = message.SetString(tag, "Sync local directory to QS-Directory: qsctl sync . qs://bucket-name", "Sync local directory to QS-Directory: qsctl sync . qs://bucket-name")
	_ = message.SetString(tag, "Sync skip updating files that already exist on receiver: qsctl sync . qs://bucket-name --ignore-existing", "Sync skip updating files that already exist on receiver: qsctl sync . qs://bucket-name --ignore-existing")
	_ = message.SetString(tag, "Tee object: qsctl tee qs://prefix/a", "Tee object: qsctl tee qs://prefix/a")
	_ = message.SetString(tag, "Write to stdout: qsctl cp qs://prefix/b - > /path/to/file", "Write to stdout: qsctl cp qs://prefix/b - > /path/to/file")
	_ = message.SetString(tag, "Your config has been set to <%v>. You can still modify it manually.", "Your config has been set to <%v>. You can still modify it manually.")
	_ = message.SetString(tag, "assign config path manually", "assign config path manually")
	_ = message.SetString(tag, "cat a remote object to stdout", "cat a remote object to stdout")
	_ = message.SetString(tag, "cat qs://<bucket_name>/<object_key>", "cat qs://<bucket_name>/<object_key>")
	_ = message.SetString(tag, "copy directory recursively", "copy directory recursively")
	_ = message.SetString(tag, "copy from/to qingstor", "copy from/to qingstor")
	_ = message.SetString(tag, "cp <source-path> <dest-path>", "cp <source-path> <dest-path>")
	_ = message.SetString(tag, "delete a bucket", "delete a bucket")
	_ = message.SetString(tag, "delete an empty bucket: qsctl rb qs://bucket-name", "delete an empty bucket: qsctl rb qs://bucket-name")
	_ = message.SetString(tag, "enable benchmark or not", "enable benchmark or not")
	_ = message.SetString(tag, "forcely delete a nonempty bucket: qsctl rb qs://bucket-name -f", "forcely delete a nonempty bucket: qsctl rb qs://bucket-name -f")
	_ = message.SetString(tag, "get the pre-signed URL by the object key", "get the pre-signed URL by the object key")
	_ = message.SetString(tag, "help for this command", "help for this command")
	_ = message.SetString(tag, "in which zone to do the operation", "in which zone to do the operation")
	_ = message.SetString(tag, "in which zone to make the bucket (required)", "in which zone to make the bucket (required)")
	_ = message.SetString(tag, "list objects or buckets", "list objects or buckets")
	_ = message.SetString(tag, "ls [qs://<bucket-name/prefix>]", "ls [qs://<bucket-name/prefix>]")
	_ = message.SetString(tag, "make a new bucket", "make a new bucket")
	_ = message.SetString(tag, "move directory recursively", "move directory recursively")
	_ = message.SetString(tag, "move from/to qingstor", "move from/to qingstor")
	_ = message.SetString(tag, "print logs for debug", "print logs for debug")
	_ = message.SetString(tag, "qsctl cat can cat a remote object to stdout", "qsctl cat can cat a remote object to stdout")
	_ = message.SetString(tag, "qsctl cp can copy file/folder/stdin to qingstor or copy qingstor objects to local/stdout", "qsctl cp can copy file/folder/stdin to qingstor or copy qingstor objects to local/stdout")
	_ = message.SetString(tag, "qsctl mv can move file/folder to qingstor or move qingstor objects to local", "qsctl mv can move file/folder to qingstor or move qingstor objects to local")
	_ = message.SetString(tag, "qsctl rb delete a qingstor bucket", "qsctl rb delete a qingstor bucket")
	_ = message.SetString(tag, "qsctl rm remove the object with given object key", "qsctl rm remove the object with given object key")
	_ = message.SetString(tag, "qsctl stat show the detailed info of this object", "qsctl stat show the detailed info of this object")
	_ = message.SetString(tag, "recursively delete keys under a specific prefix", "recursively delete keys under a specific prefix")
	_ = message.SetString(tag, "recursively list subdirectories encountered", "recursively list subdirectories encountered")
	_ = message.SetString(tag, "remove a remote object", "remove a remote object")
	_ = message.SetString(tag, "stat a remote object", "stat a remote object")
	_ = message.SetString(tag, "sync between local directory and QS-Directory", "sync between local directory and QS-Directory")
	_ = message.SetString(tag, "tee a remote object from stdin", "tee a remote object from stdin")
	_ = message.SetString(tag, "the number of seconds until the pre-signed URL expires. Default is 300 seconds", "the number of seconds until the pre-signed URL expires. Default is 300 seconds")
	_ = message.SetString(tag, `qsctl ls can list all qingstor buckets or qingstor keys under a prefix.`, `qsctl ls can list all qingstor buckets or qingstor keys under a prefix.`)
	_ = message.SetString(tag, `qsctl mb can make a new bucket with the specific name,

bucket name should follow DNS name rule with:
* length between 6 and 63;
* can only contains lowercase letters, numbers and hyphen -
* must start and end with lowercase letter or number
* must not be an available IP address
	`, `qsctl mb can make a new bucket with the specific name,

bucket name should follow DNS name rule with:
* length between 6 and 63;
* can only contains lowercase letters, numbers and hyphen -
* must start and end with lowercase letter or number
* must not be an available IP address
	`)
	_ = message.SetString(tag, `qsctl presign can generate a pre-signed URL for the object.
Within the given expire time, anyone who receives this URL can retrieve
the object with an HTTP GET request. If an object belongs to a public bucket,
generate a URL spliced by bucket name, zone and its name, anyone who receives
this URL can always retrieve the object with an HTTP GET request.`, `qsctl presign can generate a pre-signed URL for the object.
Within the given expire time, anyone who receives this URL can retrieve
the object with an HTTP GET request. If an object belongs to a public bucket,
generate a URL spliced by bucket name, zone and its name, anyone who receives
this URL can always retrieve the object with an HTTP GET request.`)
	_ = message.SetString(tag, `qsctl sync between local directory and QS-Directory. The first path argument
is the source directory and second the destination directory.

When a key(file) already exists in the destination directory, program will compare
the modified time of source file(key) and destination key(file). The destination
key(file) will be overwritten only if the source one newer than destination one.`, `qsctl sync between local directory and QS-Directory. The first path argument
is the source directory and second the destination directory.

When a key(file) already exists in the destination directory, program will compare
the modified time of source file(key) and destination key(file). The destination
key(file) will be overwritten only if the source one newer than destination one.`)
	_ = message.SetString(tag, `qsctl tee can tee a remote object from stdin.

NOTICE: qsctl will not tee the content to stdout like linux tee command does.
`, `qsctl tee can tee a remote object from stdin.

NOTICE: qsctl will not tee the content to stdout like linux tee command does.
`)
	_ = message.SetString(tag, `skip creating new files in dest dirs, only copy newer by time`, `skip creating new files in dest dirs, only copy newer by time`)
	_ = message.SetString(tag, `use the specified FORMAT instead of the default;
output a newline after each use of FORMAT

The valid format sequences for files:

  %F   file type
  %h   content md5 of the file
  %n   file name
  %s   total size, in bytes
  %y   time of last data modification, human-readable, e.g: 2006-01-02 15:04:05 +0000 UTC
  %Y   time of last data modification, seconds since Epoch
	`, `use the specified FORMAT instead of the default;
output a newline after each use of FORMAT

The valid format sequences for files:

  %F   file type
  %h   content md5 of the file
  %n   file name
  %s   total size, in bytes
  %y   time of last data modification, human-readable, e.g: 2006-01-02 15:04:05 +0000 UTC
  %Y   time of last data modification, seconds since Epoch
	`)
}
