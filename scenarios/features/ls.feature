Feature: List QingStor keys under a prefix or all QingStor buckets.

  Scenario: list all QingStor buckets.
    When list all buckets
    Then should have this bucket

  Scenario: list QingStor keys.
    Given a bucket with files
      | name                                    |
      | 中文测试无目录                                 |
      | 中文目录测试/中文测试有目录                          |
      | test_file_without_directory             |
      | test_directory/test_file_with_directory |
      | !@#$%^&*()-_=+[]{}\?<>,.;:'             |
    When list keys
    Then should list all keys
      | name                        |
      | 中文测试无目录                     |
      | 中文目录测试/                     |
      | test_file_without_directory |
      | test_directory/             |
      | !@#$%^&*()-_=+[]{}\?<>,.;:' |
    When list keys with prefix
      | prefix         |
      | test_directory |
      | 中文目录测试         |
    Then should list keys with prefix
      | prefix         | should_show_up  | not_show_up     |
      | test_directory | test_directory/ | 中文目录测试/         |
      | 中文目录测试         | 中文目录测试/         | test_directory/ |
    When list keys recursively
    Then should list keys recursively
      | name                                    |
      | 中文测试无目录                                 |
      | 中文目录测试/中文测试有目录                          |
      | test_file_without_directory             |
      | test_directory/test_file_with_directory |
      | !@#$%^&*()-_=+[]{}\?<>,.;:'             |
