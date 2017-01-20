Feature: Delete an empty QingStor bucket or forcibly delete nonempty QingStor bucket.

  Scenario: Delete an empty QingStor bucket.
    Given an empty QingStor bucket
    When delete bucket
    Then the bucket should be deleted

  Scenario: Forcibly delete nonempty QingStor bucket.
    Given a QingSto bucket with files
    When delete bucket forcibly
    Then the bucket should be deleted forcibly
