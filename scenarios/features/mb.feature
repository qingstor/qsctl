Feature: Create a QingStor bucket.

  Scenario: create a QingStor bucket.
    When create a bucket
    Then should get a new bucket
