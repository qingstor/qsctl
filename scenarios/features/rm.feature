Feature: Delete a QingStor key or keys under a prefix.

  Scenario: Delete QingStor keys
    Given several QingStor keys
      | name                                    |
      | 中文测试无目录                                 |
      | 中文目录测试/中文测试有目录                          |
      | test_file_without_directory             |
      | test_directory/test_file_with_directory |
    When delete keys
    Then QingStor keys should be deleted

  Scenario: Delete QingStor keys with prefix
    Given serveral QingStor keys with prefix
      | name                                    |
      | 中文测试无目录                                 |
      | 中文目录测试/中文测试有目录                          |
      | test_file_without_directory             |
      | test_directory/test_file_with_directory |
    When delete keys with prefix "中文目录测试"
    Then QingStor keys with prefix should be deleted, other files should keep
      | name                                    | deleted |
      | 中文测试无目录                                 | 0       |
      | 中文目录测试/中文测试有目录                          | 1       |
      | test_file_without_directory             | 0       |
      | test_directory/test_file_with_directory | 0       |
