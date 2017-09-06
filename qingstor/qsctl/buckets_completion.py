def complete(output, cmdline, point=None):
    if point is None:
        point = len(cmdline)
    args = cmdline[0:point].split()
    current_arg = args[-1]
    current_arg_point = len(current_arg)
    cmd_args = [x for x in args if not x.startswith("-")]
    qs_path = [x for x in cmd_args if x.startswith("qs://")]
    qs_opts = ["ls", "cp", "mb", "mv", "rb", "rm", "sync", "presign"]
    if len(cmd_args) <= 2:
        if current_arg in qs_opts:
            # List buckets when press <tab> after qs_opts
            bucket_list = []
            for x in output["buckets"]:
                bucket_list.append((x["name"] + "/").encode("utf-8"))
            print(" \n".join(bucket_list))
        elif current_arg == "qsctl":
            print(" \n".join(qs_opts))
        else:
            # Return matching opts
            qs_opts_reply = []
            for x in qs_opts:
                if x.startswith(current_arg):
                    qs_opts_reply.append(x + " ")
            print(" \n".join(qs_opts_reply))
    elif 2 < len(cmd_args) < 5:
        bucket_list = []
        if len(cmd_args) == 4:
            # Avoid more completions
            for x in output["buckets"]:
                if x["name"].encode("utf-8") in cmd_args[-1]:
                    return 0
        if cmd_args[1] in ["ls", "mb", "rb", "rm", "presign"]:
            # If the cmd needs only one bucket
            for x in output["buckets"]:
                for w in cmd_args:
                    if x["name"].encode("utf-8") in w:
                        return 0
        if current_arg not in qs_path:
            # Input bucket without "qs://"
            for x in output["buckets"]:
                if x["name"].encode("utf-8").startswith(current_arg):
                    bucket_list.append((x["name"] + "/").encode("utf-8"))
            print(" \n".join(bucket_list))
        else:
            # Input bucket with "qs://"
            for x in output["buckets"]:
                if x["name"].encode("utf-8").startswith(
                        current_arg[5:current_arg_point + 1]):
                    bucket_list.append(("//" + x["name"] + "/").encode("utf-8"))
            print(" \n".join(bucket_list))
    else:
        return 0
