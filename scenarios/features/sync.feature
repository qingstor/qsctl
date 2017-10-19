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

  Scenario: sync local directory to QingStor prefix using wildcard
    Given several similar local directories with files
      | name                                                |
      | tmp_similar/test_file_without_directory             |
      | tmp_similar/ab/test_file_with_directory             |
      | tmp_similar/aba/test_file_with_directory            |
      | tmp_similar/abb/test_file_with_directory            |
      | tmp_similar/aab/test_file_with_directory            |
      | tmp_similar/aaba/test_file_with_directory           |
      | tmp_similar/aabab/test_file_with_directory          |
      | tmp_similar/aabba/test_file_with_directory          |
    When sync local directories to QingStor prefix using wildcard
      | name                                                | pattern                |
      | tmp_similar/test_file_without_directory             | tmp_similar/a*b?       |
      | tmp_similar/ab/test_file_with_directory             |                        |
      | tmp_similar/aba/test_file_with_directory            |                        |
      | tmp_similar/abb/test_file_with_directory            |                        |
      | tmp_similar/aab/test_file_with_directory            |                        |
      | tmp_similar/aaba/test_file_with_directory           |                        |
      | tmp_similar/aabab/test_file_with_directory          |                        |
      | tmp_similar/aabba/test_file_with_directory          |                        |
    Then QingStor should have keys with matched prefix
      | name                                                |
      | tmp_similar/aba/test_file_with_directory            |
      | tmp_similar/abb/test_file_with_directory            |
      | tmp_similar/aaba/test_file_with_directory           |
      | tmp_similar/aabba/test_file_with_directory          |
