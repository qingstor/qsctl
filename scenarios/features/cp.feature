Feature: Copy local file(s) to QingStor or QingStor key(s) to local.

  Scenario: Copy local files to QingStor keys
    Given a set of local files
      | name           | count |
      | 中文测试           | 1     |
      | qsctl_test_big | 50    |
    When copy to QingStor key
      | name           |
      | 中文测试           |
      | qsctl_test_big |
    Then QingStor should have key
      | name           |
      | 中文测试           |
      | qsctl_test_big |
    When copy to QingStor keys recursively
    Then QingStor should have keys
      | name           |
      | 中文测试           |
      | qsctl_test_big |

  Scenario: Copy QingStor keys to local files
    When copy to local file
      | name           |
      | 中文测试           |
      | qsctl_test_big |
    Then local should have file
      | name           |
      | 中文测试           |
      | qsctl_test_big |
    When copy to local files recursively
    Then local should have files
      | name           |
      | 中文测试           |
      | qsctl_test_big |

  Scenario: Copy local files to Qingstor keys using wildcard
    Given several similar local files
      | name           |
      | ab             |
      | aba            |
      | abb            |
      | aab            |
      | aaba           |
      | aabab          |
      | aabba          |
    When copy to QingStor key using wildcard
      | name           | pattern        |
      | ab             | a*b?           |
      | aba            |                |
      | abb            |                |
      | aab            |                |
      | aaba           |                |
      | aabab          |                |
      | aabba          |                |
    Then QingStor should have matched keys
      | name           |
      | aba            |
      | abb            |
      | aaba           |
      | aabba          |
