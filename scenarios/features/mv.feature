Feature: Move local file(s) to QingStor or QingStor key(s) to local.

  Scenario: Move local files to QingStor.
    Given several local files
      | name                                    |
      | 中文测试无目录                                 |
      | 中文目录测试/中文测试有目录                          |
      | test_file_without_directory             |
      | test_directory/test_file_with_directory |
    When move to QingStor
    Then QingStor should have same file and local files should be deleted

  Scenario: Move QingStor keys to local.
    Given several QingStor keys
      | name                                    |
      | 中文测试无目录                                 |
      | 中文目录测试/中文测试有目录                          |
      | test_file_without_directory             |
      | test_directory/test_file_with_directory |
    When move to local
    Then local should have same file and QingStor keys should be deleted
