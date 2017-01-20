Feature: Sync between local directory and QingStor prefix.

  Scenario: sync local directory to QingStor prefix
    Given a local directory with files
      | name                                    |
      | 中文测试无目录                                 |
      | 中文目录测试/中文测试有目录                          |
      | test_file_without_directory             |
      | test_directory/test_file_with_directory |
    When sync local directory to QingStor prefix
    Then QingStor should have keys with prefix
      | name                                        |
      | tmp/中文测试无目录                                 |
      | tmp/中文目录测试/中文测试有目录                          |
      | tmp/test_file_without_directory             |
      | tmp/test_directory/test_file_with_directory |

  Scenario: sync QingStor prefix to local directory
    When sync QingStor prefix to local directory
    Then local should have files with prefix
      | name                                        |
      | tmp/中文测试无目录                                 |
      | tmp/中文目录测试/中文测试有目录                          |
      | tmp/test_file_without_directory             |
      | tmp/test_directory/test_file_with_directory |
